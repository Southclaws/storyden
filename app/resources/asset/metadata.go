package asset

import (
	"encoding/json"
	"strings"

	"github.com/Southclaws/storyden/internal/mime"
)

type Metadata struct {
	json.RawMessage
}

func NewMetadata(mediaType mime.Type, value any) (Metadata, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return Metadata{}, err
	}

	metadata := Metadata{RawMessage: encoded}
	if target := metadataForMIME(mediaType); target != nil {
		if err := metadata.Decode(target); err != nil {
			return Metadata{}, err
		}
	}

	return metadata, nil
}

func (m Metadata) Decode(target any) error {
	if len(m.RawMessage) == 0 {
		return nil
	}

	return json.Unmarshal(m.RawMessage, target)
}

func (m Metadata) Merge(value any) (Metadata, error) {
	fields := make(map[string]json.RawMessage)
	if err := m.Decode(&fields); err != nil {
		return Metadata{}, err
	}
	if fields == nil {
		fields = make(map[string]json.RawMessage)
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return Metadata{}, err
	}

	var updates map[string]json.RawMessage
	if err = json.Unmarshal(encoded, &updates); err != nil {
		return Metadata{}, err
	}
	for key, value := range updates {
		fields[key] = value
	}

	encoded, err = json.Marshal(fields)
	return Metadata{RawMessage: encoded}, err
}

func metadataForMIME(mediaType mime.Type) any {
	value := mediaType.String()
	if strings.HasPrefix(value, "image/") {
		return &ImageMetadata{}
	}

	switch value {
	case "text/css", "text/javascript", "application/javascript":
		return &ThemeMetadata{}
	default:
		return nil
	}
}

func decodeMetadata[T any](a Asset) (*T, error) {
	metadata, ok := metadataForMIME(a.MIME).(*T)
	if !ok {
		return nil, nil
	}

	if err := a.Metadata.Decode(metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
