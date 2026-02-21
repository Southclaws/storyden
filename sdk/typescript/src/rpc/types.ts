export interface DatagraphRef {
  // Resource ID
  id: string;
  // Resource kind (e.g., 'post', 'node', 'account')
  kind: string;
}

export type Event =
  | "EventThreadPublished"
  | "EventThreadUnpublished"
  | "EventThreadUpdated"
  | "EventThreadDeleted"
  | "EventThreadReplyCreated"
  | "EventThreadReplyDeleted"
  | "EventThreadReplyUpdated"
  | "EventThreadReplyPublished"
  | "EventThreadReplyUnpublished"
  | "EventPostLiked"
  | "EventPostUnliked"
  | "EventPostReacted"
  | "EventPostUnreacted"
  | "EventCategoryUpdated"
  | "EventCategoryDeleted"
  | "EventMemberMentioned"
  | "EventNodeCreated"
  | "EventNodeUpdated"
  | "EventNodeDeleted"
  | "EventNodePublished"
  | "EventNodeSubmittedForReview"
  | "EventNodeUnpublished"
  | "EventAccountCreated"
  | "EventAccountUpdated"
  | "EventAccountSuspended"
  | "EventAccountUnsuspended"
  | "EventReportCreated"
  | "EventReportUpdated"
  | "EventActivityCreated"
  | "EventActivityUpdated"
  | "EventActivityDeleted"
  | "EventActivityPublished"
  | "EventSettingsUpdated";


export interface EventThreadPublished {
  event: "EventThreadPublished";
  // Thread post ID
  id: string;
}

export interface EventThreadUnpublished {
  event: "EventThreadUnpublished";
  // Thread post ID
  id: string;
}

export interface EventThreadUpdated {
  event: "EventThreadUpdated";
  // Thread post ID
  id: string;
}

export interface EventThreadDeleted {
  event: "EventThreadDeleted";
  // Thread post ID
  id: string;
}

export interface EventThreadReplyCreated {
  event: "EventThreadReplyCreated";
  // Reply author account ID
  reply_author_id: string;
  // Reply post ID
  reply_id: string;
  // Optional ID of the author being replied to
  reply_to_author_id?: string;
  // Optional ID of the post being replied to
  reply_to_target_id?: string;
  // Thread author account ID
  thread_author_id: string;
  // Thread post ID
  thread_id: string;
}

export interface EventThreadReplyDeleted {
  event: "EventThreadReplyDeleted";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

export interface EventThreadReplyUpdated {
  event: "EventThreadReplyUpdated";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

export interface EventThreadReplyPublished {
  event: "EventThreadReplyPublished";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

export interface EventThreadReplyUnpublished {
  event: "EventThreadReplyUnpublished";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

export interface EventPostLiked {
  event: "EventPostLiked";
  // Post ID that was liked
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

export interface EventPostUnliked {
  event: "EventPostUnliked";
  // Post ID that was unliked
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

export interface EventPostReacted {
  event: "EventPostReacted";
  // Post ID that was reacted to
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

export interface EventPostUnreacted {
  event: "EventPostUnreacted";
  // Post ID that was unreacted
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

export interface EventCategoryUpdated {
  event: "EventCategoryUpdated";
  // Category ID
  id: string;
  // Category slug
  slug: string;
}

export interface EventCategoryDeleted {
  event: "EventCategoryDeleted";
  // Category ID
  id: string;
  // Category slug
  slug: string;
}

export interface EventMemberMentioned {
  event: "EventMemberMentioned";
  // Account ID of the member who mentioned
  by: string;
  item: DatagraphRef;
  source: DatagraphRef;
}

export interface EventNodeCreated {
  event: "EventNodeCreated";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventNodeUpdated {
  event: "EventNodeUpdated";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventNodeDeleted {
  event: "EventNodeDeleted";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventNodePublished {
  event: "EventNodePublished";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventNodeSubmittedForReview {
  event: "EventNodeSubmittedForReview";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventNodeUnpublished {
  event: "EventNodeUnpublished";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

export interface EventAccountCreated {
  event: "EventAccountCreated";
  // Account ID
  id: string;
}

export interface EventAccountUpdated {
  event: "EventAccountUpdated";
  // Account ID
  id: string;
}

export interface EventAccountSuspended {
  event: "EventAccountSuspended";
  // Account ID
  id: string;
}

export interface EventAccountUnsuspended {
  event: "EventAccountUnsuspended";
  // Account ID
  id: string;
}

export interface EventReportCreated {
  event: "EventReportCreated";
  // Report ID
  id: string;
  // Optional account ID of reporter
  reported_by?: string;
  target?: DatagraphRef;
}

export interface EventReportUpdated {
  event: "EventReportUpdated";
  // Optional account ID of handler
  handled_by?: string;
  // Report ID
  id: string;
  // Optional account ID of reporter
  reported_by?: string;
  // Report status
  status: string;
  target?: DatagraphRef;
}

export interface EventActivityCreated {
  event: "EventActivityCreated";
  // Activity/Event ID
  id: string;
}

export interface EventActivityUpdated {
  event: "EventActivityUpdated";
  // Activity/Event ID
  id: string;
}

export interface EventActivityDeleted {
  event: "EventActivityDeleted";
  // Activity/Event ID
  id: string;
}

export interface EventActivityPublished {
  event: "EventActivityPublished";
  // Activity/Event ID
  id: string;
}

export interface EventSettingsUpdated {
  event: "EventSettingsUpdated";
  // Settings object
  settings: Record<string, unknown>;
}

export type EventPayload =
  | EventThreadPublished
  | EventThreadUnpublished
  | EventThreadUpdated
  | EventThreadDeleted
  | EventThreadReplyCreated
  | EventThreadReplyDeleted
  | EventThreadReplyUpdated
  | EventThreadReplyPublished
  | EventThreadReplyUnpublished
  | EventPostLiked
  | EventPostUnliked
  | EventPostReacted
  | EventPostUnreacted
  | EventCategoryUpdated
  | EventCategoryDeleted
  | EventMemberMentioned
  | EventNodeCreated
  | EventNodeUpdated
  | EventNodeDeleted
  | EventNodePublished
  | EventNodeSubmittedForReview
  | EventNodeUnpublished
  | EventAccountCreated
  | EventAccountUpdated
  | EventAccountSuspended
  | EventAccountUnsuspended
  | EventReportCreated
  | EventReportUpdated
  | EventActivityCreated
  | EventActivityUpdated
  | EventActivityDeleted
  | EventActivityPublished
  | EventSettingsUpdated;

export function isEventThreadPublished(value: EventPayload): value is EventThreadPublished {
  return value.event === "EventThreadPublished";
}

export function isEventThreadUnpublished(value: EventPayload): value is EventThreadUnpublished {
  return value.event === "EventThreadUnpublished";
}

export function isEventThreadUpdated(value: EventPayload): value is EventThreadUpdated {
  return value.event === "EventThreadUpdated";
}

export function isEventThreadDeleted(value: EventPayload): value is EventThreadDeleted {
  return value.event === "EventThreadDeleted";
}

export function isEventThreadReplyCreated(value: EventPayload): value is EventThreadReplyCreated {
  return value.event === "EventThreadReplyCreated";
}

export function isEventThreadReplyDeleted(value: EventPayload): value is EventThreadReplyDeleted {
  return value.event === "EventThreadReplyDeleted";
}

export function isEventThreadReplyUpdated(value: EventPayload): value is EventThreadReplyUpdated {
  return value.event === "EventThreadReplyUpdated";
}

export function isEventThreadReplyPublished(value: EventPayload): value is EventThreadReplyPublished {
  return value.event === "EventThreadReplyPublished";
}

export function isEventThreadReplyUnpublished(value: EventPayload): value is EventThreadReplyUnpublished {
  return value.event === "EventThreadReplyUnpublished";
}

export function isEventPostLiked(value: EventPayload): value is EventPostLiked {
  return value.event === "EventPostLiked";
}

export function isEventPostUnliked(value: EventPayload): value is EventPostUnliked {
  return value.event === "EventPostUnliked";
}

export function isEventPostReacted(value: EventPayload): value is EventPostReacted {
  return value.event === "EventPostReacted";
}

export function isEventPostUnreacted(value: EventPayload): value is EventPostUnreacted {
  return value.event === "EventPostUnreacted";
}

export function isEventCategoryUpdated(value: EventPayload): value is EventCategoryUpdated {
  return value.event === "EventCategoryUpdated";
}

export function isEventCategoryDeleted(value: EventPayload): value is EventCategoryDeleted {
  return value.event === "EventCategoryDeleted";
}

export function isEventMemberMentioned(value: EventPayload): value is EventMemberMentioned {
  return value.event === "EventMemberMentioned";
}

export function isEventNodeCreated(value: EventPayload): value is EventNodeCreated {
  return value.event === "EventNodeCreated";
}

export function isEventNodeUpdated(value: EventPayload): value is EventNodeUpdated {
  return value.event === "EventNodeUpdated";
}

export function isEventNodeDeleted(value: EventPayload): value is EventNodeDeleted {
  return value.event === "EventNodeDeleted";
}

export function isEventNodePublished(value: EventPayload): value is EventNodePublished {
  return value.event === "EventNodePublished";
}

export function isEventNodeSubmittedForReview(value: EventPayload): value is EventNodeSubmittedForReview {
  return value.event === "EventNodeSubmittedForReview";
}

export function isEventNodeUnpublished(value: EventPayload): value is EventNodeUnpublished {
  return value.event === "EventNodeUnpublished";
}

export function isEventAccountCreated(value: EventPayload): value is EventAccountCreated {
  return value.event === "EventAccountCreated";
}

export function isEventAccountUpdated(value: EventPayload): value is EventAccountUpdated {
  return value.event === "EventAccountUpdated";
}

export function isEventAccountSuspended(value: EventPayload): value is EventAccountSuspended {
  return value.event === "EventAccountSuspended";
}

export function isEventAccountUnsuspended(value: EventPayload): value is EventAccountUnsuspended {
  return value.event === "EventAccountUnsuspended";
}

export function isEventReportCreated(value: EventPayload): value is EventReportCreated {
  return value.event === "EventReportCreated";
}

export function isEventReportUpdated(value: EventPayload): value is EventReportUpdated {
  return value.event === "EventReportUpdated";
}

export function isEventActivityCreated(value: EventPayload): value is EventActivityCreated {
  return value.event === "EventActivityCreated";
}

export function isEventActivityUpdated(value: EventPayload): value is EventActivityUpdated {
  return value.event === "EventActivityUpdated";
}

export function isEventActivityDeleted(value: EventPayload): value is EventActivityDeleted {
  return value.event === "EventActivityDeleted";
}

export function isEventActivityPublished(value: EventPayload): value is EventActivityPublished {
  return value.event === "EventActivityPublished";
}

export function isEventSettingsUpdated(value: EventPayload): value is EventSettingsUpdated {
  return value.event === "EventSettingsUpdated";
}

export interface JsonRpcRequest {
  id: string;
  jsonrpc: string;
}

export interface RPCRequestPingParams {
}


// Request sent by the host to the plugin to provide configuration settings. The params object can contain any key-value pairs defined by the plugin in its manifest "configuration_schema" field and the plugin should validate and apply these settings to its internal state.
// If configuration changes require a plugin to restart, the plugin should cleanly shut down with a zero exit code so that the Host can restart if it is a Supervised plugin. If it's an external plugin, the plugin itself is responsible for this behavior based on the plugin's lifecycle design.
export interface RPCRequestConfigure {
  method: "configure";
  id: string;
  jsonrpc: string;
  params: Record<string, unknown>;
}

export interface RPCRequestEvent {
  method: "event";
  id: string;
  jsonrpc: string;
  params: EventPayload;
}

export interface RPCRequestPing {
  method: "ping";
  id: string;
  jsonrpc: string;
  params?: RPCRequestPingParams;
}

export type HostToPluginRequest =
  | RPCRequestConfigure
  | RPCRequestEvent
  | RPCRequestPing;

export function isRPCRequestConfigure(value: HostToPluginRequest): value is RPCRequestConfigure {
  return value.method === "configure";
}

export function isRPCRequestEvent(value: HostToPluginRequest): value is RPCRequestEvent {
  return value.method === "event";
}

export function isRPCRequestPing(value: HostToPluginRequest): value is RPCRequestPing {
  return value.method === "ping";
}

export interface HostToPluginResponseError {
  code?: number;
  message?: string;
}


// Confirm that the configuration was received and applied correctly.
export interface RPCResponseConfigure {
  method: "configure";
  ok: boolean;
}

export interface RPCResponseEvent {
  method: "event";
  ok: boolean;
}

export interface RPCResponsePing {
  method: "ping";
  pong: boolean;
  // Optional status message
  status?: string;
  // How long the plugin has been running
  uptime_seconds?: number;
}

export type HostToPluginResponseUnion =
  | RPCResponseConfigure
  | RPCResponseEvent
  | RPCResponsePing;

export function isRPCResponseConfigure(value: HostToPluginResponseUnion): value is RPCResponseConfigure {
  return value.method === "configure";
}

export function isRPCResponseEvent(value: HostToPluginResponseUnion): value is RPCResponseEvent {
  return value.method === "event";
}

export function isRPCResponsePing(value: HostToPluginResponseUnion): value is RPCResponsePing {
  return value.method === "ping";
}

export interface JsonRpcResponseError {
  code?: number;
  message?: string;
}

export interface JsonRpcResponse {
  error?: JsonRpcResponseError;
  id: string;
  jsonrpc: string;
}

export interface HostToPluginResponse {
  error?: HostToPluginResponseError;
  id: string;
  jsonrpc: string;
  result: HostToPluginResponseUnion;
}

export interface ManifestAccessExternalLink {
  text: string;
  url: string;
}

export interface ManifestAccess {
  // Optional profile bio for the provisioned account.
  bio?: string;
  // The account handle to provision for this plugin's API identity.
  handle: string;
  // Optional profile links for the provisioned account.
  links?: ManifestAccessExternalLink[];
  // Optional profile metadata for the provisioned account.
  metadata?: Record<string, unknown>;
  // The account display name to provision for this plugin's API identity.
  name: string;
  // The list of permission names requested for API access.
  permissions: string[];
}


export interface PluginConfigurationFieldString {
  type: "string";
  default?: string;
  // A description of the configuration field.
  description: string;
  // A unique identifier for this configuration field, used for
  // referencing the field in the plugin configuration object.
  id: string;
  // A human-readable label for the configuration field.
  label: string;
}

export interface PluginConfigurationFieldNumber {
  type: "number";
  default?: number;
  // A description of the configuration field.
  description: string;
  // A unique identifier for this configuration field, used for
  // referencing the field in the plugin configuration object.
  id: string;
  // A human-readable label for the configuration field.
  label: string;
}

export interface PluginConfigurationFieldBoolean {
  type: "boolean";
  default?: boolean;
  // A description of the configuration field.
  description: string;
  // A unique identifier for this configuration field, used for
  // referencing the field in the plugin configuration object.
  id: string;
  // A human-readable label for the configuration field.
  label: string;
}

export type PluginConfigurationField =
  | PluginConfigurationFieldString
  | PluginConfigurationFieldNumber
  | PluginConfigurationFieldBoolean;

export function isPluginConfigurationFieldString(value: PluginConfigurationField): value is PluginConfigurationFieldString {
  return value.type === "string";
}

export function isPluginConfigurationFieldNumber(value: PluginConfigurationField): value is PluginConfigurationFieldNumber {
  return value.type === "number";
}

export function isPluginConfigurationFieldBoolean(value: PluginConfigurationField): value is PluginConfigurationFieldBoolean {
  return value.type === "boolean";
}

export type PluginConfigurationFieldSchema = PluginConfigurationField;

export interface ManifestConfigurationSchema {
  fields?: PluginConfigurationFieldSchema[];
}

export interface Manifest {
  // Optional API access configuration for this plugin. When provided, the host can provision a bot account and access key for API calls via RPC.
  access?: ManifestAccess;
  // Arguments passed to the "command" invocation.
  args?: string[];
  // The author of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
  // (NOTE: May change in future.)
  author: string;
  // The executable or script used to launch your plugin. If your plugin is a binary (Go, Rust, C, etc) then this should be a path to that binary, it's best to put it in the root of your plugin archive like `./myplugin.exe` or `./myplugin`. If your plugin is a script (Python, Node, etc) then this should be the interpreter's `$PATH` executable (e.g. `python` or `node`)  and you should include the script in the `args` field.
  // Note that Storyden cannot guarantee that the runtime environment defined by the person hosting Storyden will have any language's interpreter on the `$PATH`. If you are running your own instance and building a custom plugin, you should `FROM` the Storyden base image for your deployment so that you know what runtimes are available.
  // If you are distributing a plugin for others to use, we highly recommend that you use a statically compiled language such as Go, Rust or Zig for your plugin so that it's guaranteed to be compatible with any runtime.
  command: string;
  configuration_schema?: ManifestConfigurationSchema;
  // The description of the plugin. Displayed in Plugin Registries as well as in UI of Storyden installations when installed.
  description: string;
  // The list of events the plugin subscribes to and will receive from the host via RPC. Events allow your plugins to react to things that that humans or robots do on Storyden.
  events_consumed?: Event[];
  // The unique identifier of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
  // (NOTE: May change in future.)
  id: string;
  // The name of the plugin. This is not a unique identifier and is only used for display purposes within the Plugin Registry and Storyden installation.
  name: string;
  // The version of the plugin. This is not used for any versioning or compatibility purposes by the runtime and is only used for display purposes currently.
  version: string;
}

export interface RPCRequestGetConfigParams {
  // Specific config keys to retrieve. If empty, returns all config.
  keys?: string[];
}


export interface RPCRequestAccessGet {
  method: "access_get";
  id: string;
  jsonrpc: string;
}

export interface RPCRequestGetConfig {
  method: "get_config";
  id: string;
  jsonrpc: string;
  params?: RPCRequestGetConfigParams;
}

export type PluginToHostRequest =
  | RPCRequestAccessGet
  | RPCRequestGetConfig;

export function isRPCRequestAccessGet(value: PluginToHostRequest): value is RPCRequestAccessGet {
  return value.method === "access_get";
}

export function isRPCRequestGetConfig(value: PluginToHostRequest): value is RPCRequestGetConfig {
  return value.method === "get_config";
}

export interface PluginToHostResponseError {
  code?: number;
  message?: string;
}

export interface RPCResponseAccessGetError {
  code?: number;
  message?: string;
}

export interface RPCResponseAccessGetResult {
  // Bearer access key for API authentication.
  access_key: string;
  // Base URL for API requests.
  api_base_url: string;
}


export interface RPCResponseAccessGet {
  method: "access_get";
  error?: RPCResponseAccessGetError;
  id: string;
  jsonrpc: string;
  result: RPCResponseAccessGetResult;
}

export interface RPCResponseGetConfig {
  method: "get_config";
  // Configuration key-value pairs
  config: Record<string, unknown>;
}

export type PluginToHostResponseUnion =
  | RPCResponseAccessGet
  | RPCResponseGetConfig;

export function isRPCResponseAccessGet(value: PluginToHostResponseUnion): value is RPCResponseAccessGet {
  return value.method === "access_get";
}

export function isRPCResponseGetConfig(value: PluginToHostResponseUnion): value is RPCResponseGetConfig {
  return value.method === "get_config";
}

export interface PluginToHostResponse {
  error?: PluginToHostResponseError;
  id: string;
  jsonrpc: string;
  result: PluginToHostResponseUnion;
}
