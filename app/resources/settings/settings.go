package settings

import (
	"encoding/json"

	"dario.cat/mergo"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

// Settings is the global Storyden settings data that can be changed at runtime.
type Settings struct {
	// Title is the primary name of the instance, it's commonly used for things
	// like <title> tags in HTML for SEO as well as PWA and other app metadata.
	Title opt.Optional[string]

	// Description is similarly used for SEO tags, opengraph and other metadata.
	Description opt.Optional[string]

	// Content is a rich-text "about" field used for describing the instance.
	Content opt.Optional[datagraph.Content]

	// AccentColour is used for controlling frontend brand colour usage. Despite
	// being frontend-specific, it may be used for backend email/SMS templates.
	AccentColour opt.Optional[string]

	// Public is intended to be used to configure public access to the API. If
	// set to false any request to the API will require a verified user account.
	Public opt.Optional[bool]

	// The authentication mode is used to control which authentication methods
	// are exposed to members during the frontend registration and login flows.
	AuthenticationMode opt.Optional[authentication.Mode]

	Services opt.Optional[ServiceSettings]

	// Metadata is an arbitrary object which can be used by frontends/clients to
	// store vendor-specific configuration to control the client implementation.
	Metadata opt.Optional[map[string]any]
}

type ServiceSettings struct {
	Moderation opt.Optional[ModerationServiceSettings]
}

type ModerationServiceSettings struct {
	MaxThreadBodyLength opt.Optional[int]
	MaxReplyBodyLength  opt.Optional[int]
	WordBlocklist       opt.Optional[[]string]
}

// Merge will combine "updated" into "s" while overwriting any new values.
func (s *Settings) Merge(updated Settings) error {
	err := mergo.Merge(s, &updated, mergo.WithOverride)
	if err != nil {
		return err
	}

	return nil
}

func mapSettings(in *ent.Setting) (*Settings, error) {
	if in.ID != StorydenPrimarySettingsKey {
		return nil, fault.New("mapSettings was passed a non-system settings row")
	}

	var s Settings

	err := json.Unmarshal([]byte(in.Value), &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
