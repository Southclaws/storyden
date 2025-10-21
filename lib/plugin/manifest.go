package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
)

type Manifest struct {
	ID           ID              `json:"id"`
	Author       Author          `json:"author"`
	Name         Name            `json:"name"`
	Description  string          `json:"description,omitempty"`
	Version      *semver.Version `json:"version"`
	Capabilities []*Capability   `json:"capabilities"`
	Events       []string        `json:"events,omitempty"`
}

func (final *Manifest) UnmarshalJSON(data []byte) error {
	var errs []error

	type rawManifest Manifest
	var m rawManifest

	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to unmarshal plugin manifest: %w", err)
	}

	if !IDPattern.MatchString(string(m.ID)) {
		errs = append(errs, fmt.Errorf("invalid plugin ID: '%s'", m.ID))
	}

	if !AuthorPattern.MatchString(string(m.Author)) {
		errs = append(errs, fmt.Errorf("invalid plugin author: '%s'", m.Author))
	}

	if err := m.Name.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("invalid plugin name: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("plugin manifest validation failed: %v", errors.Join(errs...))
	}

	*final = Manifest(m)

	return nil
}

type ID string

func (id ID) String() string {
	return string(id)
}

var IDPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

type Author string

var AuthorPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

type Name string

func (name Name) String() string {
	return string(name)
}

var (
	ErrInvalidPluginName = fmt.Errorf("plugin name cannot be empty")
	ErrPluginNameTooLong = fmt.Errorf("plugin name cannot be longer than 100 characters")
)

var NameRegex = `^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,99}$`

func (name Name) Validate() error {
	if name == "" {
		return ErrInvalidPluginName
	}
	if len(name) > 100 {
		return ErrPluginNameTooLong
	}
	return nil
}
