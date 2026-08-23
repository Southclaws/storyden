package trail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Southclaws/storyden/app/resources/recurrence"
)

type eventTriggerRecord struct {
	Events []string `json:"events"`
}

type triggerSnapshotRecord struct {
	Type     TriggerType     `json:"type"`
	Schedule json.RawMessage `json:"schedule,omitempty"`
	Events   []string        `json:"events,omitempty"`
}

type triggerEventRecord struct {
	TrailID      string          `json:"trail_id"`
	TrailRunID   string          `json:"trail_run_id"`
	Kind         RunKind         `json:"kind"`
	EventName    string          `json:"event_name,omitempty"`
	Trigger      json.RawMessage `json:"trigger"`
	Payload      json.RawMessage `json:"payload,omitempty"`
	ScheduledFor *time.Time      `json:"scheduled_for,omitempty"`
	ObservedAt   time.Time       `json:"observed_at"`
	InitiatedBy  string          `json:"initiated_by,omitempty"`
}

func encodeTrigger(trigger Trigger) (TriggerType, json.RawMessage, error) {
	switch trigger.Type() {
	case TriggerTypeSchedule:
		schedule := trigger.Schedule()
		if schedule == nil {
			return TriggerType{}, nil, ErrInvalidScheduleTrigger
		}

		config, err := schedule.NormalisedJSON()
		return TriggerTypeSchedule, config, err

	case TriggerTypeEvent:
		event := trigger.Event()
		if event == nil {
			return TriggerType{}, nil, ErrInvalidEventTrigger
		}

		config, err := json.Marshal(eventTriggerRecord{Events: event.Events})
		return TriggerTypeEvent, config, err

	default:
		return TriggerType{}, nil, ErrUnsupportedTrigger
	}
}

func decodeTrigger(triggerType string, config json.RawMessage) (Trigger, error) {
	kind, err := NewTriggerType(triggerType)
	if err != nil {
		return Trigger{}, err
	}

	switch kind {
	case TriggerTypeSchedule:
		schedule, err := recurrence.Parse(config)
		if err != nil {
			return Trigger{}, err
		}
		return scheduleTrigger(schedule), nil

	case TriggerTypeEvent:
		var record eventTriggerRecord
		if err := decodeStrictJSON(config, &record); err != nil {
			return Trigger{}, fmt.Errorf("decode Trail event trigger: %w", err)
		}
		return NewEventTrigger(record.Events)

	default:
		return Trigger{}, ErrUnsupportedTrigger
	}
}

func encodeTriggerSnapshot(trigger Trigger) (json.RawMessage, error) {
	kind, config, err := encodeTrigger(trigger)
	if err != nil {
		return nil, err
	}

	switch kind {
	case TriggerTypeSchedule:
		return json.Marshal(triggerSnapshotRecord{Type: kind, Schedule: config})
	case TriggerTypeEvent:
		event := trigger.Event()
		return json.Marshal(triggerSnapshotRecord{Type: kind, Events: event.Events})
	default:
		return nil, ErrUnsupportedTrigger
	}
}

func decodeTriggerSnapshot(raw json.RawMessage) (Trigger, error) {
	var record triggerSnapshotRecord
	if err := decodeStrictJSON(raw, &record); err != nil {
		return Trigger{}, fmt.Errorf("decode Trail trigger snapshot: %w", err)
	}

	switch record.Type {
	case TriggerTypeSchedule:
		return decodeTrigger(record.Type.String(), record.Schedule)
	case TriggerTypeEvent:
		return NewEventTrigger(record.Events)
	default:
		return Trigger{}, ErrUnsupportedTrigger
	}
}

func materialiseEventPayload(eventName string, payload json.RawMessage) (json.RawMessage, error) {
	if !validEvent(eventName) {
		return nil, ErrInvalidEventTrigger
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		return nil, fmt.Errorf("decode Trail source event payload: %w", err)
	}
	if fields == nil {
		return nil, errors.New("Trail source event payload must be an object")
	}

	event, err := json.Marshal(eventName)
	if err != nil {
		return nil, err
	}
	fields["event"] = event

	encoded, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("encode Trail source event payload: %w", err)
	}
	return encoded, nil
}

func encodeTriggerEvent(event TriggerEvent) (json.RawMessage, error) {
	if event.Kind == RunKindEvent && !validEvent(event.EventName) {
		return nil, ErrInvalidEventTrigger
	}

	trigger, err := encodeTriggerSnapshot(event.Trigger)
	if err != nil {
		return nil, err
	}

	return json.Marshal(triggerEventRecord{
		TrailID:      event.TrailID,
		TrailRunID:   event.TrailRunID,
		Kind:         event.Kind,
		EventName:    event.EventName,
		Trigger:      trigger,
		Payload:      event.Payload,
		ScheduledFor: event.ScheduledFor,
		ObservedAt:   event.ObservedAt,
		InitiatedBy:  event.InitiatedBy,
	})
}

func decodeTriggerEvent(raw json.RawMessage) (TriggerEvent, error) {
	var record triggerEventRecord
	if err := decodeStrictJSON(raw, &record); err != nil {
		return TriggerEvent{}, fmt.Errorf("decode Trail trigger event: %w", err)
	}

	trigger, err := decodeTriggerSnapshot(record.Trigger)
	if err != nil {
		return TriggerEvent{}, err
	}
	if record.Kind == RunKindEvent && !validEvent(record.EventName) {
		return TriggerEvent{}, ErrInvalidEventTrigger
	}

	return TriggerEvent{
		TrailID:      record.TrailID,
		TrailRunID:   record.TrailRunID,
		Kind:         record.Kind,
		EventName:    record.EventName,
		Trigger:      trigger,
		Payload:      record.Payload,
		ScheduledFor: record.ScheduledFor,
		ObservedAt:   record.ObservedAt,
		InitiatedBy:  record.InitiatedBy,
	}, nil
}

func decodeStrictJSON(raw json.RawMessage, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}
