// Package recurrence evaluates a normalized JSON representation of the
// recurrence rule subset used by Storyden. Its semantics follow RFC 5545,
// while its public model is an AST rather than the RRULE text format.
package recurrence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

//go:generate go run github.com/Southclaws/enumerator

const LocalDateTimeLayout = "2006-01-02T15:04:05"

type frequencyEnum string

const (
	frequencyHourly  frequencyEnum = "hourly"
	frequencyDaily   frequencyEnum = "daily"
	frequencyWeekly  frequencyEnum = "weekly"
	frequencyMonthly frequencyEnum = "monthly"
	frequencyYearly  frequencyEnum = "yearly"
)

type weekdayEnum string

const (
	weekdayMonday    weekdayEnum = "monday"
	weekdayTuesday   weekdayEnum = "tuesday"
	weekdayWednesday weekdayEnum = "wednesday"
	weekdayThursday  weekdayEnum = "thursday"
	weekdayFriday    weekdayEnum = "friday"
	weekdaySaturday  weekdayEnum = "saturday"
	weekdaySunday    weekdayEnum = "sunday"
)

type Rule struct {
	Frequency  Frequency `json:"frequency"`
	Interval   int       `json:"interval"`
	ByWeekday  []Weekday `json:"by_weekday,omitempty"`
	ByMonth    []int     `json:"by_month,omitempty"`
	ByMonthDay []int     `json:"by_month_day,omitempty"`
	Count      *int      `json:"count,omitempty"`
}

type Schedule struct {
	Start    string `json:"start"`
	Timezone string `json:"timezone"`
	Rule     Rule   `json:"rule"`

	location *time.Location
	anchor   time.Time
}

// Parse decodes, validates, and normalizes a persisted recurrence definition.
func Parse(raw json.RawMessage) (*Schedule, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	var schedule Schedule
	if err := decoder.Decode(&schedule); err != nil {
		return nil, fmt.Errorf("parse recurrence: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return nil, fmt.Errorf("parse recurrence: %w", err)
	}
	return Compile(schedule)
}

// Compile validates and normalizes a typed recurrence definition.
func Compile(input Schedule) (*Schedule, error) {
	schedule := input
	schedule.Start = strings.TrimSpace(schedule.Start)
	schedule.Timezone = strings.TrimSpace(schedule.Timezone)

	location, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return nil, errors.New("timezone must be a valid IANA timezone")
	}
	parsed, err := time.Parse(LocalDateTimeLayout, schedule.Start)
	if err != nil {
		return nil, fmt.Errorf("start must use local date and time format %s", LocalDateTimeLayout)
	}
	anchor, ok := validLocal(location, parsed.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), parsed.Second())
	if !ok {
		return nil, errors.New("start must be a valid local date and time in the selected timezone")
	}

	if err := validateRule(schedule.Rule); err != nil {
		return nil, err
	}
	normaliseRule(&schedule.Rule)
	schedule.location = location
	schedule.anchor = anchor

	return &schedule, nil
}

// NormalisedJSON returns the canonical storage representation of a compiled
// recurrence definition.
func (s *Schedule) NormalisedJSON() (json.RawMessage, error) {
	encoded, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("encode recurrence: %w", err)
	}
	return encoded, nil
}

// NextAfter returns the first occurrence strictly after the supplied instant.
// Count applies to occurrences generated from the schedule anchor.
func (s *Schedule) NextAfter(after time.Time) (time.Time, bool) {
	if s == nil || s.location == nil || s.anchor.IsZero() {
		return time.Time{}, false
	}
	if s.Rule.Count == nil {
		return s.nextUnbounded(after)
	}

	cursor := s.anchor.Add(-time.Nanosecond)
	for occurrence := 0; occurrence < *s.Rule.Count; occurrence++ {
		next, ok := s.nextUnbounded(cursor)
		if !ok {
			return time.Time{}, false
		}
		if next.After(after) {
			return next, true
		}
		cursor = next
	}
	return time.Time{}, false
}

func Preview(schedule *Schedule, after time.Time, count int) []time.Time {
	if schedule == nil || count <= 0 {
		return nil
	}
	result := make([]time.Time, 0, count)
	cursor := after
	for len(result) < count {
		next, ok := schedule.NextAfter(cursor)
		if !ok {
			break
		}
		result = append(result, next.UTC())
		cursor = next
	}
	return result
}

// AdvancePast returns the latest occurrence at or before now and the next
// future occurrence.
func AdvancePast(schedule *Schedule, first, now time.Time) (time.Time, time.Time, bool) {
	latest := first
	cursor := first
	for steps := 0; steps < 10000; steps++ {
		next, ok := schedule.NextAfter(cursor)
		if !ok {
			return latest, time.Time{}, true
		}
		if next.After(now) {
			return latest, next.UTC(), true
		}
		latest = next
		cursor = next
	}
	return time.Time{}, time.Time{}, false
}

func validateRule(rule Rule) error {
	if rule.Interval < 1 {
		return errors.New("recurrence interval must be positive")
	}
	if rule.Count != nil && (*rule.Count < 1 || *rule.Count > 100000) {
		return errors.New("recurrence count must be between 1 and 100000")
	}
	if err := validateWeekdays(rule.ByWeekday); err != nil {
		return err
	}
	if err := validateIntegers("by_month", rule.ByMonth, func(value int) bool { return value >= 1 && value <= 12 }); err != nil {
		return err
	}
	if err := validateIntegers("by_month_day", rule.ByMonthDay, func(value int) bool {
		return value != 0 && value >= -31 && value <= 31
	}); err != nil {
		return err
	}

	switch rule.Frequency {
	case FrequencyHourly, FrequencyDaily:
		if len(rule.ByWeekday) > 0 || len(rule.ByMonth) > 0 || len(rule.ByMonthDay) > 0 {
			return fmt.Errorf("%s recurrence does not support selectors", rule.Frequency)
		}
	case FrequencyWeekly:
		if len(rule.ByMonth) > 0 || len(rule.ByMonthDay) > 0 {
			return errors.New("weekly recurrence supports only by_weekday")
		}
	case FrequencyMonthly:
		if len(rule.ByWeekday) > 0 {
			return errors.New("monthly recurrence does not support by_weekday")
		}
	case FrequencyYearly:
		if len(rule.ByWeekday) > 0 {
			return errors.New("yearly recurrence does not support by_weekday")
		}
		if len(rule.ByMonthDay) > 0 && len(rule.ByMonth) == 0 {
			return errors.New("yearly by_month_day requires by_month")
		}
	default:
		return fmt.Errorf("unsupported recurrence frequency %q", rule.Frequency)
	}
	return nil
}

func validateWeekdays(values []Weekday) error {
	seen := map[Weekday]bool{}
	for _, value := range values {
		if _, ok := weekdayValue(value); !ok || seen[value] {
			return errors.New("by_weekday contains an invalid or duplicate weekday")
		}
		seen[value] = true
	}
	return nil
}

func validateIntegers(name string, values []int, valid func(int) bool) error {
	seen := map[int]bool{}
	for _, value := range values {
		if !valid(value) || seen[value] {
			return fmt.Errorf("%s contains an invalid or duplicate value", name)
		}
		seen[value] = true
	}
	return nil
}

func normaliseRule(rule *Rule) {
	sort.Slice(rule.ByWeekday, func(i, j int) bool {
		left, _ := weekdayValue(rule.ByWeekday[i])
		right, _ := weekdayValue(rule.ByWeekday[j])
		return weekdayIndex(left) < weekdayIndex(right)
	})
	sort.Ints(rule.ByMonth)
	sort.Ints(rule.ByMonthDay)
}

func (s *Schedule) nextUnbounded(after time.Time) (time.Time, bool) {
	switch s.Rule.Frequency {
	case FrequencyHourly:
		return s.nextHourly(after)
	case FrequencyDaily:
		return s.nextDaily(after)
	case FrequencyWeekly:
		return s.nextWeekly(after)
	case FrequencyMonthly:
		return s.nextMonthly(after)
	case FrequencyYearly:
		return s.nextYearly(after)
	default:
		return time.Time{}, false
	}
}

func (s *Schedule) nextHourly(after time.Time) (time.Time, bool) {
	anchorWall := wallTime(s.anchor)
	afterWall := wallTime(after.In(s.location))
	period := int64(0)
	if difference := int64(afterWall.Sub(anchorWall) / time.Hour); difference > 0 {
		period = difference / int64(s.Rule.Interval)
	}
	for attempts := 0; attempts < 10000; attempts++ {
		wall := anchorWall.Add(time.Duration(period*int64(s.Rule.Interval)) * time.Hour)
		candidate, ok := validLocal(s.location, wall.Year(), wall.Month(), wall.Day(), wall.Hour(), wall.Minute(), wall.Second())
		if ok && !candidate.Before(s.anchor) && candidate.After(after) {
			return candidate, true
		}
		period++
	}
	return time.Time{}, false
}

func (s *Schedule) nextDaily(after time.Time) (time.Time, bool) {
	anchorDate := dateOnly(s.anchor)
	afterDate := dateOnly(after.In(s.location))
	days := int(afterDate.Sub(anchorDate) / (24 * time.Hour))
	period := 0
	if days > 0 {
		period = days / s.Rule.Interval
	}
	for attempts := 0; attempts < 10000; attempts++ {
		date := anchorDate.AddDate(0, 0, period*s.Rule.Interval)
		candidate, ok := validLocal(s.location, date.Year(), date.Month(), date.Day(), s.anchor.Hour(), s.anchor.Minute(), s.anchor.Second())
		if ok && !candidate.Before(s.anchor) && candidate.After(after) {
			return candidate, true
		}
		period++
	}
	return time.Time{}, false
}

func (s *Schedule) nextWeekly(after time.Time) (time.Time, bool) {
	anchorWeek := weekStart(s.anchor)
	afterWeek := weekStart(after.In(s.location))
	weeks := int(afterWeek.Sub(anchorWeek) / (7 * 24 * time.Hour))
	period := 0
	if weeks > 0 {
		period = weeks / s.Rule.Interval
	}
	weekdays := s.Rule.ByWeekday
	if len(weekdays) == 0 {
		weekdays = []Weekday{weekdayName(s.anchor.Weekday())}
	}
	for attempts := 0; attempts < 10000; attempts++ {
		start := anchorWeek.AddDate(0, 0, period*s.Rule.Interval*7)
		for _, weekday := range weekdays {
			value, _ := weekdayValue(weekday)
			date := start.AddDate(0, 0, weekdayIndex(value))
			candidate, ok := validLocal(s.location, date.Year(), date.Month(), date.Day(), s.anchor.Hour(), s.anchor.Minute(), s.anchor.Second())
			if ok && !candidate.Before(s.anchor) && candidate.After(after) {
				return candidate, true
			}
		}
		period++
	}
	return time.Time{}, false
}

func (s *Schedule) nextMonthly(after time.Time) (time.Time, bool) {
	anchorMonth := monthIndex(s.anchor.Year(), s.anchor.Month())
	afterLocal := after.In(s.location)
	months := monthIndex(afterLocal.Year(), afterLocal.Month()) - anchorMonth
	period := 0
	if months > 0 {
		period = months / s.Rule.Interval
	}
	for attempts := 0; attempts < 10000; attempts++ {
		year, month := monthFromIndex(anchorMonth + period*s.Rule.Interval)
		if matchesMonth(s.Rule.ByMonth, month) {
			for _, day := range resolvedMonthDays(year, month, s.Rule.ByMonthDay, s.anchor.Day()) {
				candidate, ok := validLocal(s.location, year, month, day, s.anchor.Hour(), s.anchor.Minute(), s.anchor.Second())
				if ok && !candidate.Before(s.anchor) && candidate.After(after) {
					return candidate, true
				}
			}
		}
		period++
	}
	return time.Time{}, false
}

func (s *Schedule) nextYearly(after time.Time) (time.Time, bool) {
	years := after.In(s.location).Year() - s.anchor.Year()
	period := 0
	if years > 0 {
		period = years / s.Rule.Interval
	}
	months := s.Rule.ByMonth
	if len(months) == 0 {
		months = []int{int(s.anchor.Month())}
	}
	for attempts := 0; attempts < 10000; attempts++ {
		year := s.anchor.Year() + period*s.Rule.Interval
		if year > 9999 {
			return time.Time{}, false
		}
		for _, monthValue := range months {
			month := time.Month(monthValue)
			for _, day := range resolvedMonthDays(year, month, s.Rule.ByMonthDay, s.anchor.Day()) {
				candidate, ok := validLocal(s.location, year, month, day, s.anchor.Hour(), s.anchor.Minute(), s.anchor.Second())
				if ok && !candidate.Before(s.anchor) && candidate.After(after) {
					return candidate, true
				}
			}
		}
		period++
	}
	return time.Time{}, false
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("recurrence contains more than one JSON value")
}

func validLocal(location *time.Location, year int, month time.Month, day, hour, minute, second int) (time.Time, bool) {
	if location == nil || day < 1 {
		return time.Time{}, false
	}
	candidate := time.Date(year, month, day, hour, minute, second, 0, location)
	if candidate.Year() != year || candidate.Month() != month || candidate.Day() != day || candidate.Hour() != hour || candidate.Minute() != minute || candidate.Second() != second {
		return time.Time{}, false
	}
	return candidate, true
}

func wallTime(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(), 0, time.UTC)
}

func dateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func weekStart(value time.Time) time.Time {
	date := dateOnly(value)
	return date.AddDate(0, 0, -weekdayIndex(value.Weekday()))
}

func weekdayIndex(value time.Weekday) int {
	return (int(value) + 6) % 7
}

func weekdayValue(day Weekday) (time.Weekday, bool) {
	values := map[Weekday]time.Weekday{
		WeekdayMonday: time.Monday, WeekdayTuesday: time.Tuesday, WeekdayWednesday: time.Wednesday,
		WeekdayThursday: time.Thursday, WeekdayFriday: time.Friday, WeekdaySaturday: time.Saturday, WeekdaySunday: time.Sunday,
	}
	value, ok := values[day]
	return value, ok
}

func weekdayName(day time.Weekday) Weekday {
	values := map[time.Weekday]Weekday{
		time.Monday: WeekdayMonday, time.Tuesday: WeekdayTuesday, time.Wednesday: WeekdayWednesday,
		time.Thursday: WeekdayThursday, time.Friday: WeekdayFriday, time.Saturday: WeekdaySaturday, time.Sunday: WeekdaySunday,
	}
	return values[day]
}

func matchesMonth(values []int, wanted time.Month) bool {
	if len(values) == 0 {
		return true
	}
	for _, value := range values {
		if time.Month(value) == wanted {
			return true
		}
	}
	return false
}

func resolvedMonthDays(year int, month time.Month, values []int, fallback int) []int {
	if len(values) == 0 {
		values = []int{fallback}
	}
	last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	seen := map[int]bool{}
	result := make([]int, 0, len(values))
	for _, value := range values {
		day := value
		if value < 0 {
			day = last + value + 1
		}
		if day < 1 || day > last || seen[day] {
			continue
		}
		seen[day] = true
		result = append(result, day)
	}
	sort.Ints(result)
	return result
}

func monthIndex(year int, month time.Month) int {
	return year*12 + int(month) - 1
}

func monthFromIndex(value int) (int, time.Month) {
	return value / 12, time.Month(value%12 + 1)
}
