package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
)

var IDPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

var AuthorPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

var (
	ErrInvalidPluginName = fmt.Errorf("plugin name cannot be empty")
	ErrPluginNameTooLong = fmt.Errorf("plugin name cannot be longer than 100 characters")
)

func ValidateName(name string) error {
	if name == "" {
		return ErrInvalidPluginName
	}
	if len(name) > 100 {
		return ErrPluginNameTooLong
	}
	return nil
}

func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse plugin manifest: %w", err)
	}
	return &m, nil
}

var validEventSet = func() map[Event]struct{} {
	s := make(map[Event]struct{}, len(EventValues))
	for _, event := range EventValues {
		s[event] = struct{}{}
	}
	return s
}()

func ManifestFromMap(m map[string]any) (*Manifest, error) {
	var manifest Manifest
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &manifest); err != nil {
		return nil, err
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (final *Manifest) UnmarshalJSON(data []byte) error {
	type rawManifest Manifest
	var m rawManifest

	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to unmarshal plugin manifest: %w", err)
	}

	*final = Manifest(m)

	if err := final.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *Manifest) Validate() error {
	internal := make([]string, 0)
	external := make([]string, 0)

	if !IDPattern.MatchString(string(m.ID)) {
		internal = append(internal, fmt.Sprintf("invalid plugin ID: '%s'", m.ID))
		external = append(external, `Field "id" is invalid. Use letters, numbers, and hyphens only.`)
	}

	if !AuthorPattern.MatchString(string(m.Author)) {
		internal = append(internal, fmt.Sprintf("invalid plugin author: '%s'", m.Author))
		external = append(external, `Field "author" is invalid. Use letters, numbers, and hyphens only.`)
	}

	if err := ValidateName(m.Name); err != nil {
		internal = append(internal, fmt.Sprintf("invalid plugin name: %v", err))
		switch {
		case errors.Is(err, ErrInvalidPluginName):
			external = append(external, `Field "name" cannot be empty.`)
		case errors.Is(err, ErrPluginNameTooLong):
			external = append(external, `Field "name" cannot be longer than 100 characters.`)
		default:
			external = append(external, `Field "name" is invalid.`)
		}
	}

	for _, event := range m.EventsConsumed {
		if _, ok := validEventSet[event]; !ok {
			internal = append(internal, fmt.Sprintf("invalid events_consumed value: %q", event))
			external = append(external, fmt.Sprintf(`Field "events_consumed" contains unknown event %q.`, event))
		}
	}

	if len(internal) > 0 {
		details := "plugin manifest validation failed: " + strings.Join(internal, "\n")
		return fault.New(
			"plugin manifest validation failed",
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(details, strings.Join(external, "\n")),
		)
	}

	return nil
}
