package trail

import (
	"errors"
	"slices"

	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var (
	ErrInvalidScheduleTrigger = errors.New("invalid Trail schedule trigger")
	ErrInvalidEventTrigger    = errors.New("invalid Trail event trigger")
)

// Trigger is a validated Trail trigger definition. Its variants can only be
// constructed through the typed constructors below or decoded by Repository.
type Trigger struct {
	kind     TriggerType
	schedule *recurrence.Schedule
	event    *EventTrigger
}

type EventTrigger struct {
	Events []string
}

func NewScheduleTrigger(input recurrence.Schedule) (Trigger, error) {
	schedule, err := recurrence.Compile(cloneSchedule(&input))
	if err != nil {
		return Trigger{}, errors.Join(ErrInvalidScheduleTrigger, err)
	}

	return scheduleTrigger(schedule), nil
}

func NewEventTrigger(events []string) (Trigger, error) {
	if !validEvents(events) {
		return Trigger{}, ErrInvalidEventTrigger
	}

	return Trigger{
		kind:  TriggerTypeEvent,
		event: &EventTrigger{Events: slices.Clone(events)},
	}, nil
}

func (t Trigger) Type() TriggerType {
	return t.kind
}

func (t Trigger) Schedule() *recurrence.Schedule {
	if t.schedule == nil {
		return nil
	}

	schedule := cloneSchedule(t.schedule)
	return &schedule
}

func (t Trigger) Event() *EventTrigger {
	if t.event == nil {
		return nil
	}

	return &EventTrigger{Events: slices.Clone(t.event.Events)}
}

func scheduleTrigger(schedule *recurrence.Schedule) Trigger {
	if schedule == nil {
		return Trigger{kind: TriggerTypeSchedule}
	}

	cloned := cloneSchedule(schedule)
	return Trigger{kind: TriggerTypeSchedule, schedule: &cloned}
}

func cloneSchedule(input *recurrence.Schedule) recurrence.Schedule {
	cloned := *input
	cloned.Rule.ByWeekday = slices.Clone(input.Rule.ByWeekday)
	cloned.Rule.ByMonth = slices.Clone(input.Rule.ByMonth)
	cloned.Rule.ByMonthDay = slices.Clone(input.Rule.ByMonthDay)
	if input.Rule.Count != nil {
		count := *input.Rule.Count
		cloned.Rule.Count = &count
	}
	return cloned
}

func validEvents(names []string) bool {
	if len(names) == 0 {
		return false
	}

	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if _, exists := seen[name]; exists || !validEvent(name) {
			return false
		}
		seen[name] = struct{}{}
	}

	return true
}

func validEvent(name string) bool {
	for _, event := range rpc.EventValues {
		if string(event) == name {
			return true
		}
	}

	return false
}
