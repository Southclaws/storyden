//go:build release

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	packageJsonPath = `./web/package.json`
	openapiPath     = `./api/openapi.yaml`
	apiVersionPath  = `./internal/config/version.go`
)

type version struct {
	one     int  // always v1
	year    int  // two-digit year
	release int  // release identifier
	post    bool // true if not a release but a post build
}

func (v *version) String() string {
	str := fmt.Sprintf("v%d.%02d.%d", v.one, v.year, v.release)
	if v.post {
		str += "-post"
	}
	return str
}

func parse(s string) (*version, error) {
	var v version
	_, err := fmt.Sscanf(s, "v%d.%d.%d", &v.one, &v.year, &v.release)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *version) next() (*version, error) {
	year, err := strconv.Atoi(time.Now().Format("06"))
	if err != nil {
		return nil, err
	}

	// Reset the release number if the year has changed.
	if v.year != year {
		v.year = year % 100 // Keep it two-digit
		v.release = 1
	} else {
		v.release++
	}

	return v, nil
}

func getMostRecentTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	var out bytes.Buffer
	var outerr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outerr

	err := cmd.Run()
	if err != nil {
		if strings.Contains(outerr.String(), "No names found, cannot describe anything.") {
			return "v1.25.0", nil
		}
		return "", errors.Join(err, fmt.Errorf("failed to get current version: %s", outerr.String()))
	}

	return strings.TrimSpace(out.String()), nil
}

func getCurrent() (*version, error) {
	tag, err := getMostRecentTag()
	if err != nil {
		return nil, err
	}

	return parse(tag)
}

func patchFile(path, pattern, replacement string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(pattern)
	out := re.ReplaceAll(data, []byte(replacement))
	return os.WriteFile(path, out, 0o644)
}

func writeVersion(v version) error {
	if err := patchFile(packageJsonPath, `"version":\s*".+?"`, fmt.Sprintf(`"version": "%s"`, v.String())); err != nil {
		return fmt.Errorf("failed to update package.json version: %w", err)
	}

	if err := patchFile(openapiPath, `version:\s*".+?"`, fmt.Sprintf(`version: "%s"`, v.String())); err != nil {
		return fmt.Errorf("failed to update openapi.yaml version: %w", err)
	}

	if err := patchFile(apiVersionPath, `Version\s*=\s*".+?"`, fmt.Sprintf(`Version = "%s"`, v.String())); err != nil {
		return fmt.Errorf("failed to update version.go: %w", err)
	}

	return nil
}

func runNext(current version, write bool) error {
	next, err := current.next()
	if err != nil {
		return fmt.Errorf("failed to generate next version: %w", err)
	}

	if write {
		if err := writeVersion(*next); err != nil {
			return fmt.Errorf("failed to write new version: %w", err)
		}
	}

	fmt.Println(next.String())

	return nil
}

func runPostRelease(current version) error {
	current.post = true
	if err := writeVersion(current); err != nil {
		return fmt.Errorf("failed to write new version: %w", err)
	}
	fmt.Println(current.String())
	return nil
}

func run(write, post bool) error {
	current, err := getCurrent()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if post {
		if err := runPostRelease(*current); err != nil {
			return fmt.Errorf("failed to run post version: %w", err)
		}
	} else {
		if err := runNext(*current, write); err != nil {
			return fmt.Errorf("failed to run next version: %w", err)
		}
	}

	return nil
}

func main() {
	write := flag.Bool("w", false, "Write the new version to files")
	post := flag.Bool("c", false, "For post-release run, run this script again with -c to write a post version to files")
	flag.Parse()

	if err := run(*write, *post); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run release script: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
