package event_ref

import (
	"net/url"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/event/location"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

type EventID xid.ID

func (e EventID) String() string { return xid.ID(e).String() }

func (e EventID) MarshalJSON() ([]byte, error) {
	return xid.ID(e).MarshalJSON()
}

func (e *EventID) UnmarshalJSON(data []byte) error {
	var id xid.ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	*e = EventID(id)
	return nil
}

type Event struct {
	ID           EventID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    opt.Optional[time.Time]
	IndexedAt    opt.Optional[time.Time]
	Name         string
	Slug         string
	Description  opt.Optional[string]
	TimeRange    TimeRange
	Asset        opt.Optional[asset.Asset]
	Policy       participation.Policy
	Visibility   visibility.Visibility
	Location     location.Location
	Capacity     opt.Optional[int]
	Participants participation.EventParticipants
	Meta         map[string]any
}

type TimeRange struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

func Map(in *ent.Event) (*Event, error) {
	assetEdge := opt.NewPtr(in.Edges.PrimaryImage)

	image := opt.Map(assetEdge, func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	policy, err := participation.NewPolicy(in.ParticipationPolicy)
	if err != nil {
		return nil, err
	}

	vis, err := visibility.NewVisibility(in.Visibility.String())
	if err != nil {
		return nil, err
	}

	loc, err := MapLocation(in)
	if err != nil {
		return nil, err
	}

	participants, err := dt.MapErr(in.Edges.Participants, participation.Map)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:           EventID(in.ID),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		DeletedAt:    opt.NewPtr(in.DeletedAt),
		IndexedAt:    opt.NewPtr(in.IndexedAt),
		Name:         in.Name,
		Slug:         in.Slug,
		Description:  opt.NewPtr(in.Description),
		TimeRange:    TimeRange{Start: in.StartTime, End: in.EndTime, Duration: in.EndTime.Sub(in.StartTime)},
		Asset:        image,
		Policy:       policy,
		Visibility:   vis,
		Location:     loc,
		Capacity:     opt.NewPtr(in.Capacity),
		Participants: participants,
		Meta:         in.Metadata,
	}, nil
}

func MapLocation(in *ent.Event) (location.Location, error) {
	if in.LocationType == nil {
		return nil, nil
	}

	lt, err := location.NewLocationType(*in.LocationType)
	if err != nil {
		return nil, err
	}

	locURL, err := opt.MapErr(opt.NewPtr(in.LocationURL), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, err
	}

	switch lt {
	case location.LocationTypePhysical:
		return &location.Physical{
			Name:      opt.NewPtr(in.LocationName).OrZero(),
			Address:   opt.NewPtr(in.LocationAddress),
			Latitude:  opt.NewPtr(in.LocationLatitude),
			Longitude: opt.NewPtr(in.LocationLongitude),
			URL:       locURL,
		}, nil

	case location.LocationTypeVirtual:
		return &location.Virtual{
			Name: opt.NewPtr(in.LocationName).OrZero(),
			URL:  locURL,
		}, nil

	default:
		return nil, fault.Newf("unhandled location type '%s'", lt)
	}
}
