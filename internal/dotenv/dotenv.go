// Package dotenv loads single-line KEY=value assignments without replacing
// existing environment variables. It supports comments, an optional export
// prefix, and single- or double-quoted values. Double quotes support \n, \r, \t,
// \\, and \" escapes. Variables and shell expressions are never expanded.
package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var doubleQuoteEscapes = strings.NewReplacer(
	`\n`, "\n",
	`\r`, "\r",
	`\t`, "\t",
	`\\`, `\`,
	`\"`, `"`,
)

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	values, err := parse(file)
	if err != nil {
		return fmt.Errorf("dotenv %s: %w", path, err)
	}

	for key, value := range values {
		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("dotenv %s: set %s: %w", path, key, err)
		}
	}

	return nil
}

func parse(reader io.Reader) (map[string]string, error) {
	values := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if lineNumber == 1 {
			line = strings.TrimPrefix(line, "\ufeff")
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") || strings.HasPrefix(line, "export\t") {
			line = strings.TrimSpace(line[len("export"):])
		}

		key, raw, found := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !found || !validKey(key) {
			return nil, fmt.Errorf("line %d: expected KEY=value with a valid environment variable name", lineNumber)
		}

		value, err := parseValue(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("line %d: read: %w", lineNumber+1, err)
	}

	return values, nil
}

func validKey(key string) bool {
	if key == "" {
		return false
	}

	for index, character := range key {
		letter := character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || character == '_'
		digit := character >= '0' && character <= '9'
		if !letter && !(index > 0 && digit) {
			return false
		}
	}

	return true
}

func parseValue(value string) (string, error) {
	if strings.ContainsRune(value, '\x00') {
		return "", fmt.Errorf("value contains a NUL byte")
	}

	if value == "" {
		return "", nil
	}

	quote := value[0]
	if quote != '\'' && quote != '"' {
		for index, character := range value {
			if character == '#' && (index == 0 || value[index-1] == ' ' || value[index-1] == '\t') {
				return strings.TrimSpace(value[:index]), nil
			}
		}

		return value, nil
	}

	for index := 1; index < len(value); index++ {
		if quote == '"' && value[index] == '\\' {
			index++
			continue
		}

		if value[index] != quote {
			continue
		}

		remainder := strings.TrimSpace(value[index+1:])
		if remainder != "" && !strings.HasPrefix(remainder, "#") {
			return "", fmt.Errorf("unexpected text after quoted value")
		}

		content := value[1:index]
		if quote == '"' {
			content = doubleQuoteEscapes.Replace(content)
		}

		return content, nil
	}

	return "", fmt.Errorf("unterminated quoted value")
}
