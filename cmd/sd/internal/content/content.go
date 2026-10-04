package content

import (
	"fmt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"io"
	"os"
)

func ToHTML(content string, markdown bool) (string, error) {
	if !markdown || content == "" {
		return content, nil
	}

	rt, err := datagraph.NewRichTextFromMarkdown(content)
	if err != nil {
		return "", fmt.Errorf("failed to convert markdown content: %w", err)
	}

	return rt.HTML(), nil
}

func Read(content string, contentFile string, stdin io.Reader) (string, error) {
	if content != "" && contentFile != "" {
		return "", fmt.Errorf("cannot specify both --content and --content-file")
	}

	if content != "" {
		return content, nil
	}

	if contentFile == "" {
		return "", nil
	}

	if contentFile == "-" {
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}

		return string(bytes), nil
	}

	bytes, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("failed to read content file: %w", err)
	}

	return string(bytes), nil
}
