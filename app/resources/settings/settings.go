package settings

import (
	"encoding/json"
	"time"

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

	// The registration mode controls whether unauthenticated visitors may create
	// accounts publicly, only via invitation, or not at all.
	RegistrationMode opt.Optional[RegistrationMode]

	Services opt.Optional[ServiceSettings]

	// Metadata is an arbitrary object which can be used by frontends/clients to
	// store vendor-specific configuration to control the client implementation.
	Metadata opt.Optional[map[string]any]

	// Motd is an optional announcement banner shown to all site visitors.
	Motd opt.Optional[MessageOfTheDay]
}

// MessageOfTheDay is a date-bound rich text announcement.
type MessageOfTheDay struct {
	Content  datagraph.Content
	StartAt  opt.Optional[time.Time]
	EndAt    opt.Optional[time.Time]
	Metadata opt.Optional[map[string]any]
}

type ServiceSettings struct {
	ClientIP   opt.Optional[ClientIPServiceSettings]
	RateLimit  opt.Optional[RateLimitServiceSettings]
	Moderation opt.Optional[ModerationServiceSettings]
}

type ClientIPServiceSettings struct {
	ClientIPMode      opt.Optional[ClientIPMode]
	ClientIPHeader    opt.Optional[string]
	TrustedProxyCIDRs opt.Optional[[]string]
}

type RateLimitServiceSettings struct {
	RateLimit          opt.Optional[int]
	RateLimitPeriod    opt.Optional[time.Duration]
	RateLimitBucket    opt.Optional[time.Duration]
	RateLimitGuestCost opt.Optional[int]
	CostOverrides      opt.Optional[map[string]int]
}

type ModerationServiceSettings struct {
	ThreadBodyLengthMax opt.Optional[int]
	ReplyBodyLengthMax  opt.Optional[int]
	SignatureLengthMax  opt.Optional[int]
	WordBlockList       opt.Optional[[]string]
	WordReportList      opt.Optional[[]string]
}

// Merge will combine "updated" into "s" while overwriting any new values.
func (s *Settings) Merge(updated Settings) error {
	if updated.Motd.Ok() {
		next := updated.Motd.OrZero()

		// A fully empty MOTD patch is used to clear the existing announcement.
		if next.Content.IsEmpty() &&
			!next.StartAt.Ok() &&
			!next.EndAt.Ok() &&
			!next.Metadata.Ok() {
			s.Motd = opt.NewEmpty[MessageOfTheDay]()
			updated.Motd = opt.NewEmpty[MessageOfTheDay]()
		}
	}

	err := mergo.Merge(s, &updated, mergo.WithOverride)
	if err != nil {
		return err
	}

	return nil
}

func (s *Settings) Clone() (*Settings, error) {
	if s == nil {
		return nil, nil
	}

	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var clone Settings
	if err := json.Unmarshal(b, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
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
