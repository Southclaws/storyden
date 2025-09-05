package plugin

import (
	"reflect"
	"regexp"
	"strings"
)

// EventNameMapping maps Go event types to snake_case event names for plugins
var EventNameMapping = map[string]string{
	// Thread events
	"message.EventThreadPublished":    "thread_published",
	"message.EventThreadUnpublished":  "thread_unpublished",
	"message.EventThreadUpdated":      "thread_updated",
	"message.EventThreadDeleted":      "thread_deleted",
	"message.EventThreadReplyCreated": "thread_reply_created",
	"message.EventThreadReplyDeleted": "thread_reply_deleted",
	"message.EventThreadReplyUpdated": "thread_reply_updated",
	"message.EventPostLiked":          "post_liked",
	"message.EventPostUnliked":        "post_unliked",
	"message.EventPostReacted":        "post_reacted",
	"message.EventMemberMentioned":    "member_mentioned",

	// Library node events
	"message.EventNodeCreated":              "node_created",
	"message.EventNodeUpdated":              "node_updated",
	"message.EventNodeDeleted":              "node_deleted",
	"message.EventNodePublished":            "node_published",
	"message.EventNodeSubmittedForReview":   "node_submitted_for_review",
	"message.EventNodeUnpublished":          "node_unpublished",

	// Account events
	"message.EventAccountCreated": "account_created",
	"message.EventAccountUpdated": "account_updated",

	// Activity events
	"message.EventActivityCreated":   "activity_created",
	"message.EventActivityUpdated":   "activity_updated",
	"message.EventActivityDeleted":   "activity_deleted",
	"message.EventActivityPublished": "activity_published",
}

var camelCaseRegex = regexp.MustCompile("([a-z0-9])([A-Z])")

// GetEventName converts a Go event type name to snake_case format
func GetEventName(eventType interface{}) string {
	typeName := getTypeName(eventType)
	
	// Check if we have a predefined mapping
	if eventName, exists := EventNameMapping[typeName]; exists {
		return eventName
	}
	
	// Fallback: convert CamelCase to snake_case
	return toSnakeCase(typeName)
}

// GetAllEventNames returns all available event names in snake_case format
func GetAllEventNames() []string {
	eventNames := make([]string, 0, len(EventNameMapping))
	for _, eventName := range EventNameMapping {
		eventNames = append(eventNames, eventName)
	}
	return eventNames
}

// ValidateEventName checks if an event name is valid
func ValidateEventName(eventName string) bool {
	for _, validName := range EventNameMapping {
		if validName == eventName {
			return true
		}
	}
	return false
}

func getTypeName(eventType interface{}) string {
	t := reflect.TypeOf(eventType)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.String()
}

func toSnakeCase(str string) string {
	// Remove package prefix if present
	if idx := strings.LastIndex(str, "."); idx != -1 {
		str = str[idx+1:]
	}
	
	// Remove "Event" prefix if present
	str = strings.TrimPrefix(str, "Event")
	
	// Convert to snake_case
	snake := camelCaseRegex.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}