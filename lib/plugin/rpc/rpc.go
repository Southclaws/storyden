package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/event/event_ref"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/report"
	"github.com/rs/xid"
	"net/url"
)

type DatagraphRef struct {
	// Resource ID
	ID xid.ID `json:"id"`
	// Resource kind (e.g., 'post', 'node', 'account')
	Kind string `json:"kind"`
}

type Event string

const (
	EventEventThreadPublished        Event = "EventThreadPublished"
	EventEventThreadUnpublished      Event = "EventThreadUnpublished"
	EventEventThreadUpdated          Event = "EventThreadUpdated"
	EventEventThreadDeleted          Event = "EventThreadDeleted"
	EventEventThreadReplyCreated     Event = "EventThreadReplyCreated"
	EventEventThreadReplyDeleted     Event = "EventThreadReplyDeleted"
	EventEventThreadReplyUpdated     Event = "EventThreadReplyUpdated"
	EventEventThreadReplyPublished   Event = "EventThreadReplyPublished"
	EventEventThreadReplyUnpublished Event = "EventThreadReplyUnpublished"
	EventEventPostLiked              Event = "EventPostLiked"
	EventEventPostUnliked            Event = "EventPostUnliked"
	EventEventPostReacted            Event = "EventPostReacted"
	EventEventPostUnreacted          Event = "EventPostUnreacted"
	EventEventCategoryUpdated        Event = "EventCategoryUpdated"
	EventEventCategoryDeleted        Event = "EventCategoryDeleted"
	EventEventMemberMentioned        Event = "EventMemberMentioned"
	EventEventNodeCreated            Event = "EventNodeCreated"
	EventEventNodeUpdated            Event = "EventNodeUpdated"
	EventEventNodeDeleted            Event = "EventNodeDeleted"
	EventEventNodePublished          Event = "EventNodePublished"
	EventEventNodeSubmittedForReview Event = "EventNodeSubmittedForReview"
	EventEventNodeUnpublished        Event = "EventNodeUnpublished"
	EventEventAccountCreated         Event = "EventAccountCreated"
	EventEventAccountUpdated         Event = "EventAccountUpdated"
	EventEventAccountSuspended       Event = "EventAccountSuspended"
	EventEventAccountUnsuspended     Event = "EventAccountUnsuspended"
	EventEventReportCreated          Event = "EventReportCreated"
	EventEventReportUpdated          Event = "EventReportUpdated"
	EventEventActivityCreated        Event = "EventActivityCreated"
	EventEventActivityUpdated        Event = "EventActivityUpdated"
	EventEventActivityDeleted        Event = "EventActivityDeleted"
	EventEventActivityPublished      Event = "EventActivityPublished"
	EventEventSettingsUpdated        Event = "EventSettingsUpdated"
)

var EventValues = []Event{
	EventEventThreadPublished,
	EventEventThreadUnpublished,
	EventEventThreadUpdated,
	EventEventThreadDeleted,
	EventEventThreadReplyCreated,
	EventEventThreadReplyDeleted,
	EventEventThreadReplyUpdated,
	EventEventThreadReplyPublished,
	EventEventThreadReplyUnpublished,
	EventEventPostLiked,
	EventEventPostUnliked,
	EventEventPostReacted,
	EventEventPostUnreacted,
	EventEventCategoryUpdated,
	EventEventCategoryDeleted,
	EventEventMemberMentioned,
	EventEventNodeCreated,
	EventEventNodeUpdated,
	EventEventNodeDeleted,
	EventEventNodePublished,
	EventEventNodeSubmittedForReview,
	EventEventNodeUnpublished,
	EventEventAccountCreated,
	EventEventAccountUpdated,
	EventEventAccountSuspended,
	EventEventAccountUnsuspended,
	EventEventReportCreated,
	EventEventReportUpdated,
	EventEventActivityCreated,
	EventEventActivityUpdated,
	EventEventActivityDeleted,
	EventEventActivityPublished,
	EventEventSettingsUpdated,
}

type EventPayloadUnion interface {
	EventPayloadType() string
	isEventPayload()
}

type EventPayload struct {
	EventPayloadUnion
}

func (w EventPayload) MarshalJSON() ([]byte, error) {
	if w.EventPayloadUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.EventPayloadUnion)
}

func (w *EventPayload) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.EventPayloadUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"event"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("EventPayload: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("EventPayload: missing discriminator field %q", "event")
	}

	var v EventPayloadUnion
	switch peek.Type {
	case "EventThreadPublished":
		v = &EventThreadPublished{}
	case "EventThreadUnpublished":
		v = &EventThreadUnpublished{}
	case "EventThreadUpdated":
		v = &EventThreadUpdated{}
	case "EventThreadDeleted":
		v = &EventThreadDeleted{}
	case "EventThreadReplyCreated":
		v = &EventThreadReplyCreated{}
	case "EventThreadReplyDeleted":
		v = &EventThreadReplyDeleted{}
	case "EventThreadReplyUpdated":
		v = &EventThreadReplyUpdated{}
	case "EventThreadReplyPublished":
		v = &EventThreadReplyPublished{}
	case "EventThreadReplyUnpublished":
		v = &EventThreadReplyUnpublished{}
	case "EventPostLiked":
		v = &EventPostLiked{}
	case "EventPostUnliked":
		v = &EventPostUnliked{}
	case "EventPostReacted":
		v = &EventPostReacted{}
	case "EventPostUnreacted":
		v = &EventPostUnreacted{}
	case "EventCategoryUpdated":
		v = &EventCategoryUpdated{}
	case "EventCategoryDeleted":
		v = &EventCategoryDeleted{}
	case "EventMemberMentioned":
		v = &EventMemberMentioned{}
	case "EventNodeCreated":
		v = &EventNodeCreated{}
	case "EventNodeUpdated":
		v = &EventNodeUpdated{}
	case "EventNodeDeleted":
		v = &EventNodeDeleted{}
	case "EventNodePublished":
		v = &EventNodePublished{}
	case "EventNodeSubmittedForReview":
		v = &EventNodeSubmittedForReview{}
	case "EventNodeUnpublished":
		v = &EventNodeUnpublished{}
	case "EventAccountCreated":
		v = &EventAccountCreated{}
	case "EventAccountUpdated":
		v = &EventAccountUpdated{}
	case "EventAccountSuspended":
		v = &EventAccountSuspended{}
	case "EventAccountUnsuspended":
		v = &EventAccountUnsuspended{}
	case "EventReportCreated":
		v = &EventReportCreated{}
	case "EventReportUpdated":
		v = &EventReportUpdated{}
	case "EventActivityCreated":
		v = &EventActivityCreated{}
	case "EventActivityUpdated":
		v = &EventActivityUpdated{}
	case "EventActivityDeleted":
		v = &EventActivityDeleted{}
	case "EventActivityPublished":
		v = &EventActivityPublished{}
	case "EventSettingsUpdated":
		v = &EventSettingsUpdated{}
	default:
		return fmt.Errorf("EventPayload: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("EventPayload: invalid %q payload: %w", peek.Type, err)
	}

	w.EventPayloadUnion = v
	return nil
}

type EventThreadPublished struct {
	Event string `json:"event"`
	// Thread post ID
	ID post.ID `json:"id"`
}

func (EventThreadPublished) isEventPayload() {}

func (EventThreadPublished) EventPayloadType() string { return "EventThreadPublished" }

type EventThreadUnpublished struct {
	Event string `json:"event"`
	// Thread post ID
	ID post.ID `json:"id"`
}

func (EventThreadUnpublished) isEventPayload() {}

func (EventThreadUnpublished) EventPayloadType() string { return "EventThreadUnpublished" }

type EventThreadUpdated struct {
	Event string `json:"event"`
	// Thread post ID
	ID post.ID `json:"id"`
}

func (EventThreadUpdated) isEventPayload() {}

func (EventThreadUpdated) EventPayloadType() string { return "EventThreadUpdated" }

type EventThreadDeleted struct {
	Event string `json:"event"`
	// Thread post ID
	ID post.ID `json:"id"`
}

func (EventThreadDeleted) isEventPayload() {}

func (EventThreadDeleted) EventPayloadType() string { return "EventThreadDeleted" }

type EventThreadReplyCreated struct {
	Event string `json:"event"`
	// Reply author account ID
	ReplyAuthorID account.AccountID `json:"reply_author_id"`
	// Reply post ID
	ReplyID post.ID `json:"reply_id"`
	// Optional ID of the author being replied to
	ReplyToAuthorID opt.Optional[account.AccountID] `json:"reply_to_author_id,omitempty"`
	// Optional ID of the post being replied to
	ReplyToTargetID opt.Optional[post.ID] `json:"reply_to_target_id,omitempty"`
	// Thread author account ID
	ThreadAuthorID account.AccountID `json:"thread_author_id"`
	// Thread post ID
	ThreadID post.ID `json:"thread_id"`
}

func (EventThreadReplyCreated) isEventPayload() {}

func (EventThreadReplyCreated) EventPayloadType() string { return "EventThreadReplyCreated" }

type EventThreadReplyDeleted struct {
	Event string `json:"event"`
	// Reply post ID
	ReplyID post.ID `json:"reply_id"`
	// Thread post ID
	ThreadID post.ID `json:"thread_id"`
}

func (EventThreadReplyDeleted) isEventPayload() {}

func (EventThreadReplyDeleted) EventPayloadType() string { return "EventThreadReplyDeleted" }

type EventThreadReplyUpdated struct {
	Event string `json:"event"`
	// Reply post ID
	ReplyID post.ID `json:"reply_id"`
	// Thread post ID
	ThreadID post.ID `json:"thread_id"`
}

func (EventThreadReplyUpdated) isEventPayload() {}

func (EventThreadReplyUpdated) EventPayloadType() string { return "EventThreadReplyUpdated" }

type EventThreadReplyPublished struct {
	Event string `json:"event"`
	// Reply post ID
	ReplyID post.ID `json:"reply_id"`
	// Thread post ID
	ThreadID post.ID `json:"thread_id"`
}

func (EventThreadReplyPublished) isEventPayload() {}

func (EventThreadReplyPublished) EventPayloadType() string { return "EventThreadReplyPublished" }

type EventThreadReplyUnpublished struct {
	Event string `json:"event"`
	// Reply post ID
	ReplyID post.ID `json:"reply_id"`
	// Thread post ID
	ThreadID post.ID `json:"thread_id"`
}

func (EventThreadReplyUnpublished) isEventPayload() {}

func (EventThreadReplyUnpublished) EventPayloadType() string { return "EventThreadReplyUnpublished" }

type EventPostLiked struct {
	Event string `json:"event"`
	// Post ID that was liked
	PostID post.ID `json:"post_id"`
	// Root thread post ID
	RootPostID post.ID `json:"root_post_id"`
}

func (EventPostLiked) isEventPayload() {}

func (EventPostLiked) EventPayloadType() string { return "EventPostLiked" }

type EventPostUnliked struct {
	Event string `json:"event"`
	// Post ID that was unliked
	PostID post.ID `json:"post_id"`
	// Root thread post ID
	RootPostID post.ID `json:"root_post_id"`
}

func (EventPostUnliked) isEventPayload() {}

func (EventPostUnliked) EventPayloadType() string { return "EventPostUnliked" }

type EventPostReacted struct {
	Event string `json:"event"`
	// Post ID that was reacted to
	PostID post.ID `json:"post_id"`
	// Root thread post ID
	RootPostID post.ID `json:"root_post_id"`
}

func (EventPostReacted) isEventPayload() {}

func (EventPostReacted) EventPayloadType() string { return "EventPostReacted" }

type EventPostUnreacted struct {
	Event string `json:"event"`
	// Post ID that was unreacted
	PostID post.ID `json:"post_id"`
	// Root thread post ID
	RootPostID post.ID `json:"root_post_id"`
}

func (EventPostUnreacted) isEventPayload() {}

func (EventPostUnreacted) EventPayloadType() string { return "EventPostUnreacted" }

type EventCategoryUpdated struct {
	Event string `json:"event"`
	// Category ID
	ID xid.ID `json:"id"`
	// Category slug
	Slug string `json:"slug"`
}

func (EventCategoryUpdated) isEventPayload() {}

func (EventCategoryUpdated) EventPayloadType() string { return "EventCategoryUpdated" }

type EventCategoryDeleted struct {
	Event string `json:"event"`
	// Category ID
	ID xid.ID `json:"id"`
	// Category slug
	Slug string `json:"slug"`
}

func (EventCategoryDeleted) isEventPayload() {}

func (EventCategoryDeleted) EventPayloadType() string { return "EventCategoryDeleted" }

type EventMemberMentioned struct {
	// Account ID of the member who mentioned
	By     account.AccountID `json:"by"`
	Event  string            `json:"event"`
	Item   DatagraphRef      `json:"item"`
	Source DatagraphRef      `json:"source"`
}

func (EventMemberMentioned) isEventPayload() {}

func (EventMemberMentioned) EventPayloadType() string { return "EventMemberMentioned" }

type EventNodeCreated struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodeCreated) isEventPayload() {}

func (EventNodeCreated) EventPayloadType() string { return "EventNodeCreated" }

type EventNodeUpdated struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodeUpdated) isEventPayload() {}

func (EventNodeUpdated) EventPayloadType() string { return "EventNodeUpdated" }

type EventNodeDeleted struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodeDeleted) isEventPayload() {}

func (EventNodeDeleted) EventPayloadType() string { return "EventNodeDeleted" }

type EventNodePublished struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodePublished) isEventPayload() {}

func (EventNodePublished) EventPayloadType() string { return "EventNodePublished" }

type EventNodeSubmittedForReview struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodeSubmittedForReview) isEventPayload() {}

func (EventNodeSubmittedForReview) EventPayloadType() string { return "EventNodeSubmittedForReview" }

type EventNodeUnpublished struct {
	Event string `json:"event"`
	// Library node ID
	ID library.NodeID `json:"id"`
	// Node slug
	Slug string `json:"slug"`
}

func (EventNodeUnpublished) isEventPayload() {}

func (EventNodeUnpublished) EventPayloadType() string { return "EventNodeUnpublished" }

type EventAccountCreated struct {
	Event string `json:"event"`
	// Account ID
	ID account.AccountID `json:"id"`
}

func (EventAccountCreated) isEventPayload() {}

func (EventAccountCreated) EventPayloadType() string { return "EventAccountCreated" }

type EventAccountUpdated struct {
	Event string `json:"event"`
	// Account ID
	ID account.AccountID `json:"id"`
}

func (EventAccountUpdated) isEventPayload() {}

func (EventAccountUpdated) EventPayloadType() string { return "EventAccountUpdated" }

type EventAccountSuspended struct {
	Event string `json:"event"`
	// Account ID
	ID account.AccountID `json:"id"`
}

func (EventAccountSuspended) isEventPayload() {}

func (EventAccountSuspended) EventPayloadType() string { return "EventAccountSuspended" }

type EventAccountUnsuspended struct {
	Event string `json:"event"`
	// Account ID
	ID account.AccountID `json:"id"`
}

func (EventAccountUnsuspended) isEventPayload() {}

func (EventAccountUnsuspended) EventPayloadType() string { return "EventAccountUnsuspended" }

type EventReportCreated struct {
	Event string `json:"event"`
	// Report ID
	ID report.ID `json:"id"`
	// Optional account ID of reporter
	ReportedBy opt.Optional[account.AccountID] `json:"reported_by,omitempty"`
	Target     opt.Optional[DatagraphRef]      `json:"target,omitempty"`
}

func (EventReportCreated) isEventPayload() {}

func (EventReportCreated) EventPayloadType() string { return "EventReportCreated" }

type EventReportUpdated struct {
	Event string `json:"event"`
	// Optional account ID of handler
	HandledBy opt.Optional[account.AccountID] `json:"handled_by,omitempty"`
	// Report ID
	ID report.ID `json:"id"`
	// Optional account ID of reporter
	ReportedBy opt.Optional[account.AccountID] `json:"reported_by,omitempty"`
	// Report status
	Status report.Status              `json:"status"`
	Target opt.Optional[DatagraphRef] `json:"target,omitempty"`
}

func (EventReportUpdated) isEventPayload() {}

func (EventReportUpdated) EventPayloadType() string { return "EventReportUpdated" }

type EventActivityCreated struct {
	Event string `json:"event"`
	// Activity/Event ID
	ID event_ref.EventID `json:"id"`
}

func (EventActivityCreated) isEventPayload() {}

func (EventActivityCreated) EventPayloadType() string { return "EventActivityCreated" }

type EventActivityUpdated struct {
	Event string `json:"event"`
	// Activity/Event ID
	ID event_ref.EventID `json:"id"`
}

func (EventActivityUpdated) isEventPayload() {}

func (EventActivityUpdated) EventPayloadType() string { return "EventActivityUpdated" }

type EventActivityDeleted struct {
	Event string `json:"event"`
	// Activity/Event ID
	ID event_ref.EventID `json:"id"`
}

func (EventActivityDeleted) isEventPayload() {}

func (EventActivityDeleted) EventPayloadType() string { return "EventActivityDeleted" }

type EventActivityPublished struct {
	Event string `json:"event"`
	// Activity/Event ID
	ID event_ref.EventID `json:"id"`
}

func (EventActivityPublished) isEventPayload() {}

func (EventActivityPublished) EventPayloadType() string { return "EventActivityPublished" }

type EventSettingsUpdated struct {
	Event string `json:"event"`
	// Settings object
	Settings map[string]interface{} `json:"settings"`
}

func (EventSettingsUpdated) isEventPayload() {}

func (EventSettingsUpdated) EventPayloadType() string { return "EventSettingsUpdated" }

type JsonRpcRequest struct {
	ID      xid.ID `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

type RPCRequestPingParams struct {
}

type HostToPluginRequestUnion interface {
	HostToPluginRequestType() string
	isHostToPluginRequest()
}

type HostToPluginRequest struct {
	HostToPluginRequestUnion
}

func (w HostToPluginRequest) MarshalJSON() ([]byte, error) {
	if w.HostToPluginRequestUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.HostToPluginRequestUnion)
}

func (w *HostToPluginRequest) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.HostToPluginRequestUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"method"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("HostToPluginRequest: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("HostToPluginRequest: missing discriminator field %q", "method")
	}

	var v HostToPluginRequestUnion
	switch peek.Type {
	case "configure":
		v = &RPCRequestConfigure{}
	case "event":
		v = &RPCRequestEvent{}
	case "ping":
		v = &RPCRequestPing{}
	default:
		return fmt.Errorf("HostToPluginRequest: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("HostToPluginRequest: invalid %q payload: %w", peek.Type, err)
	}

	w.HostToPluginRequestUnion = v
	return nil
}

// Request sent by the host to the plugin to provide configuration settings. The params object can contain any key-value pairs defined by the plugin in its manifest "configuration_schema" field and the plugin should validate and apply these settings to its internal state.
// If configuration changes require a plugin to restart, the plugin should cleanly shut down with a zero exit code so that the Host can restart if it is a Supervised plugin. If it's an external plugin, the plugin itself is responsible for this behavior based on the plugin's lifecycle design.
type RPCRequestConfigure struct {
	ID      xid.ID                 `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

func (RPCRequestConfigure) isHostToPluginRequest() {}

func (RPCRequestConfigure) HostToPluginRequestType() string { return "configure" }

type RPCRequestEvent struct {
	ID      xid.ID       `json:"id"`
	Jsonrpc string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  EventPayload `json:"params"`
}

func (RPCRequestEvent) isHostToPluginRequest() {}

func (RPCRequestEvent) HostToPluginRequestType() string { return "event" }

type RPCRequestPing struct {
	ID      xid.ID                             `json:"id"`
	Jsonrpc string                             `json:"jsonrpc"`
	Method  string                             `json:"method"`
	Params  opt.Optional[RPCRequestPingParams] `json:"params,omitempty"`
}

func (RPCRequestPing) isHostToPluginRequest() {}

func (RPCRequestPing) HostToPluginRequestType() string { return "ping" }

type HostToPluginResponseError struct {
	Code    opt.Optional[int]    `json:"code,omitempty"`
	Message opt.Optional[string] `json:"message,omitempty"`
}

type HostToPluginResponseUnionUnion interface {
	HostToPluginResponseUnionType() string
	isHostToPluginResponseUnion()
}

type HostToPluginResponseUnion struct {
	HostToPluginResponseUnionUnion
}

func (w HostToPluginResponseUnion) MarshalJSON() ([]byte, error) {
	if w.HostToPluginResponseUnionUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.HostToPluginResponseUnionUnion)
}

func (w *HostToPluginResponseUnion) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.HostToPluginResponseUnionUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"method"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("HostToPluginResponseUnion: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("HostToPluginResponseUnion: missing discriminator field %q", "method")
	}

	var v HostToPluginResponseUnionUnion
	switch peek.Type {
	case "configure":
		v = &RPCResponseConfigure{}
	case "event":
		v = &RPCResponseEvent{}
	case "ping":
		v = &RPCResponsePing{}
	default:
		return fmt.Errorf("HostToPluginResponseUnion: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("HostToPluginResponseUnion: invalid %q payload: %w", peek.Type, err)
	}

	w.HostToPluginResponseUnionUnion = v
	return nil
}

// Confirm that the configuration was received and applied correctly.
type RPCResponseConfigure struct {
	Method opt.Optional[string] `json:"method,omitempty"`
	Ok     bool                 `json:"ok"`
}

func (RPCResponseConfigure) isHostToPluginResponseUnion() {}

func (RPCResponseConfigure) HostToPluginResponseUnionType() string { return "configure" }

type RPCResponseEvent struct {
	Method opt.Optional[string] `json:"method,omitempty"`
	Ok     bool                 `json:"ok"`
}

func (RPCResponseEvent) isHostToPluginResponseUnion() {}

func (RPCResponseEvent) HostToPluginResponseUnionType() string { return "event" }

type RPCResponsePing struct {
	Method opt.Optional[string] `json:"method,omitempty"`
	Pong   bool                 `json:"pong"`
	// Optional status message
	Status opt.Optional[string] `json:"status,omitempty"`
	// How long the plugin has been running
	UptimeSeconds opt.Optional[float64] `json:"uptime_seconds,omitempty"`
}

func (RPCResponsePing) isHostToPluginResponseUnion() {}

func (RPCResponsePing) HostToPluginResponseUnionType() string { return "ping" }

type JsonRpcResponseError struct {
	Code    opt.Optional[int]    `json:"code,omitempty"`
	Message opt.Optional[string] `json:"message,omitempty"`
}

type JsonRpcResponse struct {
	Error   opt.Optional[JsonRpcResponseError] `json:"error,omitempty"`
	ID      xid.ID                             `json:"id"`
	Jsonrpc string                             `json:"jsonrpc"`
}

type HostToPluginResponse struct {
	Error   opt.Optional[HostToPluginResponseError] `json:"error,omitempty"`
	ID      xid.ID                                  `json:"id"`
	Jsonrpc string                                  `json:"jsonrpc"`
	Result  HostToPluginResponseUnion               `json:"result"`
}

type ManifestAccessExternalLink struct {
	Text string  `json:"text"`
	URL  url.URL `json:"url"`
}

type ManifestAccess struct {
	// Optional profile bio for the provisioned account.
	Bio opt.Optional[string] `json:"bio,omitempty"`
	// The account handle to provision for this plugin's API identity.
	//
	Handle string `json:"handle"`
	// Optional profile links for the provisioned account.
	Links []ManifestAccessExternalLink `json:"links,omitempty"`
	// Optional profile metadata for the provisioned account.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// The account display name to provision for this plugin's API identity.
	//
	Name string `json:"name"`
	// The list of permission names requested for API access.
	//
	Permissions []string `json:"permissions"`
}

type PluginConfigurationFieldUnion interface {
	PluginConfigurationFieldType() string
	isPluginConfigurationField()
}

type PluginConfigurationField struct {
	PluginConfigurationFieldUnion
}

func (w PluginConfigurationField) MarshalJSON() ([]byte, error) {
	if w.PluginConfigurationFieldUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.PluginConfigurationFieldUnion)
}

func (w *PluginConfigurationField) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.PluginConfigurationFieldUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("PluginConfigurationField: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("PluginConfigurationField: missing discriminator field %q", "type")
	}

	var v PluginConfigurationFieldUnion
	switch peek.Type {
	case "string":
		v = &PluginConfigurationFieldString{}
	case "number":
		v = &PluginConfigurationFieldNumber{}
	case "boolean":
		v = &PluginConfigurationFieldBoolean{}
	default:
		return fmt.Errorf("PluginConfigurationField: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("PluginConfigurationField: invalid %q payload: %w", peek.Type, err)
	}

	w.PluginConfigurationFieldUnion = v
	return nil
}

type PluginConfigurationFieldString struct {
	Default opt.Optional[string] `json:"default,omitempty"`
	// A description of the configuration field.
	Description string `json:"description"`
	// A unique identifier for this configuration field, used for
	// referencing the field in the plugin configuration object.
	//
	ID string `json:"id"`
	// A human-readable label for the configuration field.
	Label string `json:"label"`
	Type  string `json:"type"`
}

func (PluginConfigurationFieldString) isPluginConfigurationField() {}

func (PluginConfigurationFieldString) PluginConfigurationFieldType() string { return "string" }

type PluginConfigurationFieldNumber struct {
	Default opt.Optional[float64] `json:"default,omitempty"`
	// A description of the configuration field.
	Description string `json:"description"`
	// A unique identifier for this configuration field, used for
	// referencing the field in the plugin configuration object.
	//
	ID string `json:"id"`
	// A human-readable label for the configuration field.
	Label string `json:"label"`
	Type  string `json:"type"`
}

func (PluginConfigurationFieldNumber) isPluginConfigurationField() {}

func (PluginConfigurationFieldNumber) PluginConfigurationFieldType() string { return "number" }

type PluginConfigurationFieldBoolean struct {
	Default opt.Optional[bool] `json:"default,omitempty"`
	// A description of the configuration field.
	Description string `json:"description"`
	// A unique identifier for this configuration field, used for
	// referencing the field in the plugin configuration object.
	//
	ID string `json:"id"`
	// A human-readable label for the configuration field.
	Label string `json:"label"`
	Type  string `json:"type"`
}

func (PluginConfigurationFieldBoolean) isPluginConfigurationField() {}

func (PluginConfigurationFieldBoolean) PluginConfigurationFieldType() string { return "boolean" }

type PluginConfigurationFieldSchema = PluginConfigurationField

type ManifestConfigurationSchema struct {
	Fields []PluginConfigurationFieldSchema `json:"fields,omitempty"`
}

type Manifest struct {
	// Optional API access configuration for this plugin. When provided, the host can provision a bot account and access key for API calls via RPC.
	//
	Access opt.Optional[ManifestAccess] `json:"access,omitempty"`
	// Arguments passed to the "command" invocation.
	Args []string `json:"args,omitempty"`
	// The author of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
	// (NOTE: May change in future.)
	//
	Author string `json:"author"`
	// The executable or script used to launch your plugin. If your plugin is a binary (Go, Rust, C, etc) then this should be a path to that binary, it's best to put it in the root of your plugin archive like `./myplugin.exe` or `./myplugin`. If your plugin is a script (Python, Node, etc) then this should be the interpreter's `$PATH` executable (e.g. `python` or `node`)  and you should include the script in the `args` field.
	// Note that Storyden cannot guarantee that the runtime environment defined by the person hosting Storyden will have any language's interpreter on the `$PATH`. If you are running your own instance and building a custom plugin, you should `FROM` the Storyden base image for your deployment so that you know what runtimes are available.
	// If you are distributing a plugin for others to use, we highly recommend that you use a statically compiled language such as Go, Rust or Zig for your plugin so that it's guaranteed to be compatible with any runtime.
	//
	Command             string                                    `json:"command"`
	ConfigurationSchema opt.Optional[ManifestConfigurationSchema] `json:"configuration_schema,omitempty"`
	// The description of the plugin. Displayed in Plugin Registries as well as in UI of Storyden installations when installed.
	//
	Description string `json:"description"`
	// The list of events the plugin subscribes to and will receive from the host via RPC. Events allow your plugins to react to things that that humans or robots do on Storyden.
	//
	EventsConsumed []Event `json:"events_consumed,omitempty"`
	// The unique identifier of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
	// (NOTE: May change in future.)
	//
	ID string `json:"id"`
	// The name of the plugin. This is not a unique identifier and is only used for display purposes within the Plugin Registry and Storyden installation.
	//
	Name string `json:"name"`
	// The version of the plugin. This is not used for any versioning or compatibility purposes by the runtime and is only used for display purposes currently.
	//
	Version string `json:"version"`
}

type RPCRequestGetConfigParams struct {
	// Specific config keys to retrieve. If empty, returns all config.
	Keys []string `json:"keys,omitempty"`
}

type PluginToHostRequestUnion interface {
	PluginToHostRequestType() string
	isPluginToHostRequest()
}

type PluginToHostRequest struct {
	PluginToHostRequestUnion
}

func (w PluginToHostRequest) MarshalJSON() ([]byte, error) {
	if w.PluginToHostRequestUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.PluginToHostRequestUnion)
}

func (w *PluginToHostRequest) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.PluginToHostRequestUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"method"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("PluginToHostRequest: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("PluginToHostRequest: missing discriminator field %q", "method")
	}

	var v PluginToHostRequestUnion
	switch peek.Type {
	case "access_get":
		v = &RPCRequestAccessGet{}
	case "get_config":
		v = &RPCRequestGetConfig{}
	default:
		return fmt.Errorf("PluginToHostRequest: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("PluginToHostRequest: invalid %q payload: %w", peek.Type, err)
	}

	w.PluginToHostRequestUnion = v
	return nil
}

type RPCRequestAccessGet struct {
	ID      xid.ID `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
}

func (RPCRequestAccessGet) isPluginToHostRequest() {}

func (RPCRequestAccessGet) PluginToHostRequestType() string { return "access_get" }

type RPCRequestGetConfig struct {
	ID      xid.ID                                  `json:"id"`
	Jsonrpc string                                  `json:"jsonrpc"`
	Method  string                                  `json:"method"`
	Params  opt.Optional[RPCRequestGetConfigParams] `json:"params,omitempty"`
}

func (RPCRequestGetConfig) isPluginToHostRequest() {}

func (RPCRequestGetConfig) PluginToHostRequestType() string { return "get_config" }

type PluginToHostResponseError struct {
	Code    opt.Optional[int]    `json:"code,omitempty"`
	Message opt.Optional[string] `json:"message,omitempty"`
}

type RPCResponseAccessGetError struct {
	Code    opt.Optional[int]    `json:"code,omitempty"`
	Message opt.Optional[string] `json:"message,omitempty"`
}

type RPCResponseAccessGetResult struct {
	// Bearer access key for API authentication.
	AccessKey string `json:"access_key"`
	// Base URL for API requests.
	APIBaseURL url.URL `json:"api_base_url"`
}

type PluginToHostResponseUnionUnion interface {
	PluginToHostResponseUnionType() string
	isPluginToHostResponseUnion()
}

type PluginToHostResponseUnion struct {
	PluginToHostResponseUnionUnion
}

func (w PluginToHostResponseUnion) MarshalJSON() ([]byte, error) {
	if w.PluginToHostResponseUnionUnion == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.PluginToHostResponseUnionUnion)
}

func (w *PluginToHostResponseUnion) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		w.PluginToHostResponseUnionUnion = nil
		return nil
	}

	var peek struct {
		Type string `json:"method"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("PluginToHostResponseUnion: invalid JSON: %w", err)
	}
	if peek.Type == "" {
		return fmt.Errorf("PluginToHostResponseUnion: missing discriminator field %q", "method")
	}

	var v PluginToHostResponseUnionUnion
	switch peek.Type {
	case "access_get":
		v = &RPCResponseAccessGet{}
	case "get_config":
		v = &RPCResponseGetConfig{}
	default:
		return fmt.Errorf("PluginToHostResponseUnion: unknown type %q", peek.Type)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("PluginToHostResponseUnion: invalid %q payload: %w", peek.Type, err)
	}

	w.PluginToHostResponseUnionUnion = v
	return nil
}

type RPCResponseAccessGet struct {
	Error   opt.Optional[RPCResponseAccessGetError] `json:"error,omitempty"`
	ID      xid.ID                                  `json:"id"`
	Jsonrpc string                                  `json:"jsonrpc"`
	Method  opt.Optional[string]                    `json:"method,omitempty"`
	Result  RPCResponseAccessGetResult              `json:"result"`
}

func (RPCResponseAccessGet) isPluginToHostResponseUnion() {}

func (RPCResponseAccessGet) PluginToHostResponseUnionType() string { return "access_get" }

type RPCResponseGetConfig struct {
	// Configuration key-value pairs
	Config map[string]interface{} `json:"config"`
	Method string                 `json:"method"`
}

func (RPCResponseGetConfig) isPluginToHostResponseUnion() {}

func (RPCResponseGetConfig) PluginToHostResponseUnionType() string { return "get_config" }

type PluginToHostResponse struct {
	Error   opt.Optional[PluginToHostResponseError] `json:"error,omitempty"`
	ID      xid.ID                                  `json:"id"`
	Jsonrpc string                                  `json:"jsonrpc"`
	Result  PluginToHostResponseUnion               `json:"result"`
}
