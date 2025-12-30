package timerange

import (
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/opt"
)

type TimeRange struct {
	Start opt.Optional[time.Time]
	End   opt.Optional[time.Time]
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err == nil {
		return t, nil
	}

	return time.Parse(time.RFC3339, s)
}

func Parse(rangeStr string) (TimeRange, error) {
	if rangeStr == "" {
		return TimeRange{}, nil
	}

	parts := strings.Split(rangeStr, "/")

	switch len(parts) {
	case 1:
		start, err := parseTime(parts[0])
		if err != nil {
			return TimeRange{}, err
		}
		return TimeRange{
			Start: opt.New(start),
			End:   opt.NewEmpty[time.Time](),
		}, nil

	case 2:
		if parts[0] == "" {
			end, err := parseTime(parts[1])
			if err != nil {
				return TimeRange{}, err
			}
			return TimeRange{
				Start: opt.NewEmpty[time.Time](),
				End:   opt.New(end),
			}, nil
		} else if parts[1] == "" {
			start, err := parseTime(parts[0])
			if err != nil {
				return TimeRange{}, err
			}
			return TimeRange{
				Start: opt.New(start),
				End:   opt.NewEmpty[time.Time](),
			}, nil
		} else {
			start, err := parseTime(parts[0])
			if err != nil {
				return TimeRange{}, err
			}
			end, err := parseTime(parts[1])
			if err != nil {
				return TimeRange{}, err
			}
			return TimeRange{
				Start: opt.New(start),
				End:   opt.New(end),
			}, nil
		}

	default:
		return TimeRange{}, fmt.Errorf("invalid time range format: expected 'start/end' or 'start/' or '/end' or 'start', got %q", rangeStr)
	}
}
