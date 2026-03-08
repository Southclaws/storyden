from __future__ import annotations

from typing import Annotated, Any, Dict, List, Literal, Union
from enum import Enum
from pydantic import AnyUrl, BaseModel, ConfigDict, Field




class DatagraphRef(BaseModel):
    model_config = ConfigDict(extra="forbid")

    """Resource ID"""
    id: str
    """Resource kind (e.g., 'post', 'node', 'account')"""
    kind: str


class Event(str, Enum):
    EVENTTHREADPUBLISHED = "EventThreadPublished"
    EVENTTHREADUNPUBLISHED = "EventThreadUnpublished"
    EVENTTHREADUPDATED = "EventThreadUpdated"
    EVENTTHREADDELETED = "EventThreadDeleted"
    EVENTTHREADREPLYCREATED = "EventThreadReplyCreated"
    EVENTTHREADREPLYDELETED = "EventThreadReplyDeleted"
    EVENTTHREADREPLYUPDATED = "EventThreadReplyUpdated"
    EVENTTHREADREPLYPUBLISHED = "EventThreadReplyPublished"
    EVENTTHREADREPLYUNPUBLISHED = "EventThreadReplyUnpublished"
    EVENTPOSTLIKED = "EventPostLiked"
    EVENTPOSTUNLIKED = "EventPostUnliked"
    EVENTPOSTREACTED = "EventPostReacted"
    EVENTPOSTUNREACTED = "EventPostUnreacted"
    EVENTCATEGORYUPDATED = "EventCategoryUpdated"
    EVENTCATEGORYDELETED = "EventCategoryDeleted"
    EVENTMEMBERMENTIONED = "EventMemberMentioned"
    EVENTNODECREATED = "EventNodeCreated"
    EVENTNODEUPDATED = "EventNodeUpdated"
    EVENTNODEDELETED = "EventNodeDeleted"
    EVENTNODEPUBLISHED = "EventNodePublished"
    EVENTNODESUBMITTEDFORREVIEW = "EventNodeSubmittedForReview"
    EVENTNODEUNPUBLISHED = "EventNodeUnpublished"
    EVENTACCOUNTCREATED = "EventAccountCreated"
    EVENTACCOUNTUPDATED = "EventAccountUpdated"
    EVENTACCOUNTSUSPENDED = "EventAccountSuspended"
    EVENTACCOUNTUNSUSPENDED = "EventAccountUnsuspended"
    EVENTREPORTCREATED = "EventReportCreated"
    EVENTREPORTUPDATED = "EventReportUpdated"
    EVENTACTIVITYCREATED = "EventActivityCreated"
    EVENTACTIVITYUPDATED = "EventActivityUpdated"
    EVENTACTIVITYDELETED = "EventActivityDeleted"
    EVENTACTIVITYPUBLISHED = "EventActivityPublished"
    EVENTSETTINGSUPDATED = "EventSettingsUpdated"


"""Emitted when a thread is visible as published, either on create or after a visibility change."""

class EventThreadPublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadPublished"]
    """Thread post ID"""
    id: str


"""Emitted when a previously published thread transitions to a non-published visibility."""

class EventThreadUnpublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadUnpublished"]
    """Thread post ID"""
    id: str


"""Emitted after a thread update succeeds, regardless of which fields changed."""

class EventThreadUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadUpdated"]
    """Thread post ID"""
    id: str


"""Emitted after a thread is deleted."""

class EventThreadDeleted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadDeleted"]
    """Thread post ID"""
    id: str


"""Emitted when a new reply is created and remains published (not moved to review by moderation)."""

class EventThreadReplyCreated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadReplyCreated"]
    """Reply author account ID"""
    reply_author_id: str
    """Reply post ID"""
    reply_id: str
    """Optional ID of the author being replied to"""
    reply_to_author_id: str | None = None
    """Optional ID of the post being replied to"""
    reply_to_target_id: str | None = None
    """Thread author account ID"""
    thread_author_id: str
    """Thread post ID"""
    thread_id: str


"""Emitted after a reply is deleted from a thread."""

class EventThreadReplyDeleted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadReplyDeleted"]
    """Reply post ID"""
    reply_id: str
    """Thread post ID"""
    thread_id: str


"""Emitted after a reply update succeeds."""

class EventThreadReplyUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadReplyUpdated"]
    """Reply post ID"""
    reply_id: str
    """Thread post ID"""
    thread_id: str


"""Emitted when a reply visibility transitions to published."""

class EventThreadReplyPublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadReplyPublished"]
    """Reply post ID"""
    reply_id: str
    """Thread post ID"""
    thread_id: str


"""Emitted when a previously published reply transitions to a non-published visibility."""

class EventThreadReplyUnpublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventThreadReplyUnpublished"]
    """Reply post ID"""
    reply_id: str
    """Thread post ID"""
    thread_id: str


"""Emitted when a member adds a like to a post."""

class EventPostLiked(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventPostLiked"]
    """Post ID that was liked"""
    post_id: str
    """Root thread post ID"""
    root_post_id: str


"""Emitted when a member removes a like from a post."""

class EventPostUnliked(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventPostUnliked"]
    """Post ID that was unliked"""
    post_id: str
    """Root thread post ID"""
    root_post_id: str


"""Emitted when a member adds an emoji reaction to a post."""

class EventPostReacted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventPostReacted"]
    """Post ID that was reacted to"""
    post_id: str
    """Root thread post ID"""
    root_post_id: str


"""Emitted when a member removes an emoji reaction from a post."""

class EventPostUnreacted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventPostUnreacted"]
    """Post ID that was unreacted"""
    post_id: str
    """Root thread post ID"""
    root_post_id: str


"""Emitted when a category is created, updated, or moved."""

class EventCategoryUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventCategoryUpdated"]
    """Category ID"""
    id: str
    """Category slug"""
    slug: str


"""Emitted after a category is deleted."""

class EventCategoryDeleted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventCategoryDeleted"]
    """Category ID"""
    id: str
    """Category slug"""
    slug: str


"""Emitted once per mention target when a member mention is detected (self-mentions are skipped)."""

class EventMemberMentioned(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventMemberMentioned"]
    """Account ID of the member who mentioned"""
    by: str
    """Reference to the mentioned item."""
    item: DatagraphRef
    """Reference to the content item where the mention was found."""
    source: DatagraphRef


"""Emitted after a library node is created."""

class EventNodeCreated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodeCreated"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted after a library node is updated, moved, re-ordered, or affected by property schema changes."""

class EventNodeUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodeUpdated"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted after a library node is deleted."""

class EventNodeDeleted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodeDeleted"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted when a library node becomes published, either on create or after a visibility change."""

class EventNodePublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodePublished"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted when a library node transitions to review visibility."""

class EventNodeSubmittedForReview(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodeSubmittedForReview"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted when a previously published library node transitions to draft, unlisted, or review."""

class EventNodeUnpublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventNodeUnpublished"]
    """Library node ID"""
    id: str
    """Node slug"""
    slug: str


"""Emitted after a new account is created."""

class EventAccountCreated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventAccountCreated"]
    """Account ID"""
    id: str


"""Emitted after account profile, email, or role assignment changes."""

class EventAccountUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventAccountUpdated"]
    """Account ID"""
    id: str


"""Emitted when an account is suspended."""

class EventAccountSuspended(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventAccountSuspended"]
    """Account ID"""
    id: str


"""Emitted when a suspended account is reinstated."""

class EventAccountUnsuspended(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventAccountUnsuspended"]
    """Account ID"""
    id: str


"""Emitted when a new member or system report is created."""

class EventReportCreated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventReportCreated"]
    """Report ID"""
    id: str
    """Optional account ID of reporter, not set if it was an automated system report."""
    reported_by: str | None = None
    """Reference to the item that was reported."""
    target: DatagraphRef | None = None


"""Emitted when a report is updated, including status and handler changes."""

class EventReportUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventReportUpdated"]
    """Optional account ID of handler"""
    handled_by: str | None = None
    """Report ID"""
    id: str
    """Optional account ID of reporter"""
    reported_by: str | None = None
    """Report status"""
    status: str
    """Reference to the item that the report is about."""
    target: DatagraphRef | None = None


"""Emitted after an activity event is created."""

class EventActivityCreated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventActivityCreated"]
    """Activity/Event ID"""
    id: str


"""Emitted after an activity event is updated."""

class EventActivityUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventActivityUpdated"]
    """Activity/Event ID"""
    id: str


"""Emitted after an activity event is deleted."""

class EventActivityDeleted(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventActivityDeleted"]
    """Activity/Event ID"""
    id: str


"""Emitted when an activity event becomes published, either on create or after a visibility change."""

class EventActivityPublished(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventActivityPublished"]
    """Activity/Event ID"""
    id: str


"""Emitted when site settings are updated via the admin settings manager."""

class EventSettingsUpdated(BaseModel):
    model_config = ConfigDict(extra="forbid")
    event: Literal["EventSettingsUpdated"]
    """Settings object"""
    settings: Dict[str, Any]

EventPayload = Annotated[
    Union[EventThreadPublished, EventThreadUnpublished, EventThreadUpdated, EventThreadDeleted, EventThreadReplyCreated, EventThreadReplyDeleted, EventThreadReplyUpdated, EventThreadReplyPublished, EventThreadReplyUnpublished, EventPostLiked, EventPostUnliked, EventPostReacted, EventPostUnreacted, EventCategoryUpdated, EventCategoryDeleted, EventMemberMentioned, EventNodeCreated, EventNodeUpdated, EventNodeDeleted, EventNodePublished, EventNodeSubmittedForReview, EventNodeUnpublished, EventAccountCreated, EventAccountUpdated, EventAccountSuspended, EventAccountUnsuspended, EventReportCreated, EventReportUpdated, EventActivityCreated, EventActivityUpdated, EventActivityDeleted, EventActivityPublished, EventSettingsUpdated],
    Field(discriminator="event"),
]



class JsonRpcRequest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    id: str
    jsonrpc: str


class RPCRequestPingParams(BaseModel):
    model_config = ConfigDict(extra="forbid")

    pass


"""
Request sent by the host to the plugin to provide configuration settings. The params object can contain any key-value pairs defined by the plugin in its manifest `configuration_schema` field and the plugin should validate and apply these settings to its internal state.
If configuration changes require a plugin to restart, the plugin should cleanly shut down with a zero exit code so that the host can restart it if it is a supervised plugin. If it is an external plugin, the plugin itself is responsible for this behavior based on the plugin's lifecycle design.
"""

class RPCRequestConfigure(JsonRpcRequest):
    method: Literal["configure"]
    params: Dict[str, Any]


"""Delivers a subscribed Storyden event payload to the plugin."""

class RPCRequestEvent(JsonRpcRequest):
    method: Literal["event"]
    params: EventPayload


"""Health-check request sent by the host to verify plugin responsiveness."""

class RPCRequestPing(JsonRpcRequest):
    method: Literal["ping"]
    params: RPCRequestPingParams | None = None

HostToPluginRequest = Annotated[
    Union[RPCRequestConfigure, RPCRequestEvent, RPCRequestPing],
    Field(discriminator="method"),
]



class HostToPluginResponseError(BaseModel):
    model_config = ConfigDict(extra="forbid")

    code: int | None = None
    message: str | None = None


"""Confirms that the configuration was received and applied correctly."""

class RPCResponseConfigure(BaseModel):
    model_config = ConfigDict(extra="forbid")
    method: Literal["configure"]
    ok: bool


"""Acknowledges that the plugin received the event payload."""

class RPCResponseEvent(BaseModel):
    model_config = ConfigDict(extra="forbid")
    method: Literal["event"]
    ok: bool


"""Health-check response from the plugin."""

class RPCResponsePing(BaseModel):
    model_config = ConfigDict(extra="forbid")
    method: Literal["ping"]
    pong: bool
    """Optional status message"""
    status: str | None = None
    """How long the plugin has been running"""
    uptime_seconds: float | None = None

HostToPluginResponseUnion = Annotated[
    Union[RPCResponseConfigure, RPCResponseEvent, RPCResponsePing],
    Field(discriminator="method"),
]



class JsonRpcResponseError(BaseModel):
    model_config = ConfigDict(extra="forbid")

    code: int | None = None
    message: str | None = None


class JsonRpcResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    error: JsonRpcResponseError | None = None
    id: str
    jsonrpc: str


class HostToPluginResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    result: HostToPluginResponseUnion


class ManifestAccessExternalLink(BaseModel):
    model_config = ConfigDict(extra="forbid")

    text: str
    url: AnyUrl


class ManifestAccess(BaseModel):
    model_config = ConfigDict(extra="forbid")

    """Optional profile bio for the provisioned account."""
    bio: str | None = None
    """The account handle to provision for this plugin's API identity."""
    handle: str
    """Optional profile links for the provisioned account."""
    links: List[ManifestAccessExternalLink] | None = None
    """Optional profile metadata for the provisioned account."""
    metadata: Dict[str, Any] | None = None
    """The account display name to provision for this plugin's API identity."""
    name: str
    """The list of permission names requested for API access. See https://storyden.org/docs/introduction/members/permissions for available values and descriptions."""
    permissions: List[str]


class PluginConfigurationFieldString(BaseModel):
    model_config = ConfigDict(extra="forbid")
    type: Literal["string"]
    default: str | None = None
    """A description of the configuration field."""
    description: str
    """
    A unique identifier for this configuration field, used for
    referencing the field in the plugin configuration object.
    """
    id: str
    """A human-readable label for the configuration field."""
    label: str


class PluginConfigurationFieldNumber(BaseModel):
    model_config = ConfigDict(extra="forbid")
    type: Literal["number"]
    default: float | None = None
    """A description of the configuration field."""
    description: str
    """
    A unique identifier for this configuration field, used for
    referencing the field in the plugin configuration object.
    """
    id: str
    """A human-readable label for the configuration field."""
    label: str


class PluginConfigurationFieldBoolean(BaseModel):
    model_config = ConfigDict(extra="forbid")
    type: Literal["boolean"]
    default: bool | None = None
    """A description of the configuration field."""
    description: str
    """
    A unique identifier for this configuration field, used for
    referencing the field in the plugin configuration object.
    """
    id: str
    """A human-readable label for the configuration field."""
    label: str

PluginConfigurationField = Annotated[
    Union[PluginConfigurationFieldString, PluginConfigurationFieldNumber, PluginConfigurationFieldBoolean],
    Field(discriminator="type"),
]



PluginConfigurationFieldSchema = PluginConfigurationField


class ManifestConfigurationSchema(BaseModel):
    model_config = ConfigDict(extra="forbid")

    fields: List[PluginConfigurationFieldSchema] | None = None


class Manifest(BaseModel):
    model_config = ConfigDict(extra="forbid")

    """Optional API access configuration for this plugin. When provided, the host can provision a bot account and access key for API calls via RPC."""
    access: ManifestAccess | None = None
    """
    Arguments passed to the "command" invocation.
    This field is used only for Supervised plugins. External plugins can omit it or provide placeholder values.
    """
    args: List[str] | None = None
    """
    The author of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
    (NOTE: May change in future.)
    """
    author: str
    """
    The executable or script used to launch your plugin. If your plugin is a binary (Go, Rust, C, etc) then this should be a path to that binary, it's best to put it in the root of your plugin archive like `./myplugin.exe` or `./myplugin`. If your plugin is a script (Python, Node, etc) then this should be the interpreter's `$PATH` executable (e.g. `python` or `node`)  and you should include the script in the `args` field.
    This field is used only for Supervised plugins. External plugins can provide a placeholder value and it will be ignored by the runtime.
    Note that Storyden cannot guarantee that the runtime environment defined by the person hosting Storyden will have any language's interpreter on the `$PATH`. If you are running your own instance and building a custom plugin, you should `FROM` the Storyden base image for your deployment so that you know what runtimes are available.
    If you are distributing a plugin for others to use, we highly recommend that you use a statically compiled language such as Go, Rust or Zig for your plugin so that it's guaranteed to be compatible with any runtime.
    """
    command: str
    configuration_schema: ManifestConfigurationSchema | None = None
    """The description of the plugin. Displayed in Plugin Registries as well as in UI of Storyden installations when installed."""
    description: str
    """The list of events the plugin subscribes to and will receive from the host via RPC. Events allow your plugins to react to things that humans or robots do on Storyden."""
    events_consumed: List[Event] | None = None
    """
    The unique identifier of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
    (NOTE: May change in future.)
    """
    id: str
    """The name of the plugin. This is not a unique identifier and is only used for display purposes within the Plugin Registry and Storyden installation."""
    name: str
    """The version of the plugin. This is not used for any versioning or compatibility purposes by the runtime and is only used for display purposes currently."""
    version: str


class RPCRequestGetConfigParams(BaseModel):
    model_config = ConfigDict(extra="forbid")

    """Specific config keys to retrieve. If empty, returns all config."""
    keys: List[str] | None = None


"""Request API access credentials provisioned for this plugin."""

class RPCRequestAccessGet(JsonRpcRequest):
    method: Literal["access_get"]


"""Request the plugin's current stored configuration from the host."""

class RPCRequestGetConfig(JsonRpcRequest):
    method: Literal["get_config"]
    params: RPCRequestGetConfigParams | None = None

PluginToHostRequest = Annotated[
    Union[RPCRequestAccessGet, RPCRequestGetConfig],
    Field(discriminator="method"),
]



class PluginToHostResponseError(BaseModel):
    model_config = ConfigDict(extra="forbid")

    code: int | None = None
    message: str | None = None


class RPCResponseAccessGetError(BaseModel):
    model_config = ConfigDict(extra="forbid")

    code: int | None = None
    message: str | None = None


class RPCResponseAccessGetResult(BaseModel):
    model_config = ConfigDict(extra="forbid")

    """Bearer access key for API authentication."""
    access_key: str
    """Base URL for API requests."""
    apibase_url: AnyUrl


"""Returns API base URL and bearer access key for authenticated API calls."""

class RPCResponseAccessGet(JsonRpcResponse):
    method: Literal["access_get"]
    result: RPCResponseAccessGetResult


"""Returns current configuration values for this plugin."""

class RPCResponseGetConfig(BaseModel):
    model_config = ConfigDict(extra="forbid")
    method: Literal["get_config"]
    """Configuration key-value pairs"""
    config: Dict[str, Any]

PluginToHostResponseUnion = Annotated[
    Union[RPCResponseAccessGet, RPCResponseGetConfig],
    Field(discriminator="method"),
]



class PluginToHostResponse(BaseModel):
    model_config = ConfigDict(extra="forbid")

    result: PluginToHostResponseUnion

