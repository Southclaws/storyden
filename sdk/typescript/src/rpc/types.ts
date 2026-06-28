export interface CapabilityConfigBase {
  description?: string;
  // Stable capability identifier. For Robot LLM providers, this is the provider identifier used in `provider/model` refs. For Robot tool providers, this is the provider namespace used to build fully qualified Robot tool names.
  id: string;
  // Human-readable capability name.
  name?: string;
  // Capability type discriminator.
  type: string;
  // Version of the host capability protocol implemented by this plugin.
  version: string;
}

export interface RobotToolAnnotations {
  destructiveHint?: boolean;
  idempotentHint?: boolean;
  openWorldHint?: boolean;
  readOnlyHint?: boolean;
}

export type RobotToolJSONSchema = unknown;

export interface RobotToolProviderToolConfig {
  annotations?: RobotToolAnnotations;
  description: string;
  // Stable provider-local tool identifier.
  id: string;
  input_schema: RobotToolJSONSchema;
  // Human-readable tool name.
  name: string;
  output_schema?: RobotToolJSONSchema;
  requires_confirmation?: boolean;
  // Optional display title for MCP clients and Robot builders.
  title?: string;
}


export interface RobotLLMProviderCapabilityConfig {
  type: "robot.llm_provider";
  description?: string;
  // Whether this provider implements the embedding RPC for semantic indexing and retrieval.
  embeddings?: boolean;
  // Stable capability identifier. For Robot LLM providers, this is the provider identifier used in `provider/model` refs. For Robot tool providers, this is the provider namespace used to build fully qualified Robot tool names.
  id: string;
  // Human-readable capability name.
  name?: string;
  // Whether this provider implements the structured prompt RPC for JSON Schema constrained single-shot outputs.
  structured_output?: boolean;
  // Version of the host capability protocol implemented by this plugin.
  version: string;
}

export interface RobotToolProviderCapabilityConfig {
  type: "robot.tool_provider";
  description?: string;
  // Stable capability identifier. For Robot LLM providers, this is the provider identifier used in `provider/model` refs. For Robot tool providers, this is the provider namespace used to build fully qualified Robot tool names.
  id: string;
  // Human-readable capability name.
  name?: string;
  // Robot tools statically provided by this plugin.
  tools: RobotToolProviderToolConfig[];
  // Version of the host capability protocol implemented by this plugin.
  version: string;
}

export type CapabilityConfig =
  | RobotLLMProviderCapabilityConfig
  | RobotToolProviderCapabilityConfig;

export function isRobotLLMProviderCapabilityConfig(value: CapabilityConfig): value is RobotLLMProviderCapabilityConfig {
  return value.type === "robot.llm_provider";
}

export function isRobotToolProviderCapabilityConfig(value: CapabilityConfig): value is RobotToolProviderCapabilityConfig {
  return value.type === "robot.tool_provider";
}

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
  | "EventNodeVersionDraftCreated"
  | "EventNodeVersionDraftUpdated"
  | "EventNodeVersionDraftDeleted"
  | "EventNodeVersionDraftApplied"
  | "EventAccountCreated"
  | "EventAccountUpdated"
  | "EventAccountSuspended"
  | "EventAccountUnsuspended"
  | "EventModerationNoteCreated"
  | "EventModerationNoteDeleted"
  | "EventAccountWarned"
  | "EventAccountWarningUpdated"
  | "EventAccountWarningDeleted"
  | "EventReportCreated"
  | "EventReportUpdated"
  | "EventActivityCreated"
  | "EventActivityUpdated"
  | "EventActivityDeleted"
  | "EventActivityPublished"
  | "EventSettingsUpdated";


// Emitted when a thread is visible as published, either on create or after a visibility change.
export interface EventThreadPublished {
  event: "EventThreadPublished";
  // Thread post ID
  id: string;
}

// Emitted when a previously published thread transitions to a non-published visibility.
export interface EventThreadUnpublished {
  event: "EventThreadUnpublished";
  // Thread post ID
  id: string;
}

// Emitted after a thread update succeeds, regardless of which fields changed.
export interface EventThreadUpdated {
  event: "EventThreadUpdated";
  // Thread post ID
  id: string;
}

// Emitted after a thread is deleted.
export interface EventThreadDeleted {
  event: "EventThreadDeleted";
  // Thread post ID
  id: string;
}

// Emitted when a new reply is created and remains published (not moved to review by moderation).
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

// Emitted after a reply is deleted from a thread.
export interface EventThreadReplyDeleted {
  event: "EventThreadReplyDeleted";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

// Emitted after a reply update succeeds.
export interface EventThreadReplyUpdated {
  event: "EventThreadReplyUpdated";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

// Emitted when a reply visibility transitions to published.
export interface EventThreadReplyPublished {
  event: "EventThreadReplyPublished";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

// Emitted when a previously published reply transitions to a non-published visibility.
export interface EventThreadReplyUnpublished {
  event: "EventThreadReplyUnpublished";
  // Reply post ID
  reply_id: string;
  // Thread post ID
  thread_id: string;
}

// Emitted when a member adds a like to a post.
export interface EventPostLiked {
  event: "EventPostLiked";
  // Post ID that was liked
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

// Emitted when a member removes a like from a post.
export interface EventPostUnliked {
  event: "EventPostUnliked";
  // Post ID that was unliked
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

// Emitted when a member adds an emoji reaction to a post.
export interface EventPostReacted {
  event: "EventPostReacted";
  // Post ID that was reacted to
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

// Emitted when a member removes an emoji reaction from a post.
export interface EventPostUnreacted {
  event: "EventPostUnreacted";
  // Post ID that was unreacted
  post_id: string;
  // Root thread post ID
  root_post_id: string;
}

// Emitted when a category is created, updated, or moved.
export interface EventCategoryUpdated {
  event: "EventCategoryUpdated";
  // Category ID
  id: string;
  // Category slug
  slug: string;
}

// Emitted after a category is deleted.
export interface EventCategoryDeleted {
  event: "EventCategoryDeleted";
  // Category ID
  id: string;
  // Category slug
  slug: string;
}

// Emitted once per mention target when a member mention is detected (self-mentions are skipped).
export interface EventMemberMentioned {
  event: "EventMemberMentioned";
  // Account ID of the member who mentioned
  by: string;
  // Reference to the mentioned item.
  item: DatagraphRef;
  // Reference to the content item where the mention was found.
  source: DatagraphRef;
}

// Emitted after a library node is created.
export interface EventNodeCreated {
  event: "EventNodeCreated";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted after a library node is updated, moved, re-ordered, or affected by property schema changes.
export interface EventNodeUpdated {
  event: "EventNodeUpdated";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted after a library node is deleted.
export interface EventNodeDeleted {
  event: "EventNodeDeleted";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted when a library node becomes published, either on create or after a visibility change.
export interface EventNodePublished {
  event: "EventNodePublished";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted when a library node transitions to review visibility.
export interface EventNodeSubmittedForReview {
  event: "EventNodeSubmittedForReview";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted when a previously published library node transitions to draft, unlisted, or review.
export interface EventNodeUnpublished {
  event: "EventNodeUnpublished";
  // Library node ID
  id: string;
  // Node slug
  slug: string;
}

// Emitted when a draft version is created for a library node.
export interface EventNodeVersionDraftCreated {
  event: "EventNodeVersionDraftCreated";
  // Draft author account ID
  author_id: string;
  // Library node ID
  node_id: string;
  // Node slug
  node_slug: string;
  // Draft version ID
  version_id: string;
}

// Emitted when a draft version is updated.
export interface EventNodeVersionDraftUpdated {
  event: "EventNodeVersionDraftUpdated";
  // Draft author account ID
  author_id: string;
  // Library node ID
  node_id: string;
  // Node slug
  node_slug: string;
  // Draft version ID
  version_id: string;
}

// Emitted when a draft version is discarded.
export interface EventNodeVersionDraftDeleted {
  event: "EventNodeVersionDraftDeleted";
  // Draft author account ID
  author_id: string;
  // Library node ID
  node_id: string;
  // Node slug
  node_slug: string;
  // Draft version ID
  version_id: string;
}

// Emitted when a draft version is applied to a library node.
export interface EventNodeVersionDraftApplied {
  event: "EventNodeVersionDraftApplied";
  // Draft author account ID
  author_id: string;
  // Library node ID
  node_id: string;
  // Node slug
  node_slug: string;
  // Draft version ID
  version_id: string;
}

// Emitted after a new account is created.
export interface EventAccountCreated {
  event: "EventAccountCreated";
  // Account ID
  id: string;
}

// Emitted after account profile, email, or role assignment changes.
export interface EventAccountUpdated {
  event: "EventAccountUpdated";
  // Account ID
  id: string;
}

// Emitted when an account is suspended.
export interface EventAccountSuspended {
  event: "EventAccountSuspended";
  // Account ID
  id: string;
}

// Emitted when a suspended account is reinstated.
export interface EventAccountUnsuspended {
  event: "EventAccountUnsuspended";
  // Account ID
  id: string;
}

// Emitted when a moderator creates an internal account moderation note.
export interface EventModerationNoteCreated {
  event: "EventModerationNoteCreated";
  // Target account ID
  account_id: string;
  // Moderation note ID
  note_id: string;
}

// Emitted when a moderator deletes an internal account moderation note.
export interface EventModerationNoteDeleted {
  event: "EventModerationNoteDeleted";
  // Target account ID
  account_id: string;
  // Moderation note ID
  note_id: string;
}

// Emitted when a moderation warning is issued to an account.
export interface EventAccountWarned {
  event: "EventAccountWarned";
  // Account ID that received the warning.
  account_id: string;
  // Account ID that issued the warning.
  author_id: string;
  // Warning record ID.
  warning_id: string;
}

// Emitted when a moderation warning is edited.
export interface EventAccountWarningUpdated {
  event: "EventAccountWarningUpdated";
  // Account ID whose warning was edited.
  account_id: string;
  // Account ID that edited the warning.
  author_id: string;
  // Warning reason before the edit.
  previous_reason: string;
  // Updated warning reason.
  reason: string;
  // Warning record ID.
  warning_id: string;
}

// Emitted when a moderation warning is permanently deleted.
export interface EventAccountWarningDeleted {
  event: "EventAccountWarningDeleted";
  // Account ID whose warning was deleted.
  account_id: string;
  // Account ID that deleted the warning.
  author_id: string;
  // Warning record ID.
  warning_id: string;
}

// Emitted when a new member or system report is created.
export interface EventReportCreated {
  event: "EventReportCreated";
  // Report ID
  id: string;
  // Optional account ID of reporter, not set if it was an automated system report.
  reported_by?: string;
  // Reference to the item that was reported.
  target?: DatagraphRef;
}

// Emitted when a report is updated, including status and handler changes.
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
  // Reference to the item that the report is about.
  target?: DatagraphRef;
}

// Emitted after an activity event is created.
export interface EventActivityCreated {
  event: "EventActivityCreated";
  // Activity/Event ID
  id: string;
}

// Emitted after an activity event is updated.
export interface EventActivityUpdated {
  event: "EventActivityUpdated";
  // Activity/Event ID
  id: string;
}

// Emitted after an activity event is deleted.
export interface EventActivityDeleted {
  event: "EventActivityDeleted";
  // Activity/Event ID
  id: string;
}

// Emitted when an activity event becomes published, either on create or after a visibility change.
export interface EventActivityPublished {
  event: "EventActivityPublished";
  // Activity/Event ID
  id: string;
}

// Emitted when site settings are updated via the admin settings manager.
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
  | EventNodeVersionDraftCreated
  | EventNodeVersionDraftUpdated
  | EventNodeVersionDraftDeleted
  | EventNodeVersionDraftApplied
  | EventAccountCreated
  | EventAccountUpdated
  | EventAccountSuspended
  | EventAccountUnsuspended
  | EventModerationNoteCreated
  | EventModerationNoteDeleted
  | EventAccountWarned
  | EventAccountWarningUpdated
  | EventAccountWarningDeleted
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

export function isEventNodeVersionDraftCreated(value: EventPayload): value is EventNodeVersionDraftCreated {
  return value.event === "EventNodeVersionDraftCreated";
}

export function isEventNodeVersionDraftUpdated(value: EventPayload): value is EventNodeVersionDraftUpdated {
  return value.event === "EventNodeVersionDraftUpdated";
}

export function isEventNodeVersionDraftDeleted(value: EventPayload): value is EventNodeVersionDraftDeleted {
  return value.event === "EventNodeVersionDraftDeleted";
}

export function isEventNodeVersionDraftApplied(value: EventPayload): value is EventNodeVersionDraftApplied {
  return value.event === "EventNodeVersionDraftApplied";
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

export function isEventModerationNoteCreated(value: EventPayload): value is EventModerationNoteCreated {
  return value.event === "EventModerationNoteCreated";
}

export function isEventModerationNoteDeleted(value: EventPayload): value is EventModerationNoteDeleted {
  return value.event === "EventModerationNoteDeleted";
}

export function isEventAccountWarned(value: EventPayload): value is EventAccountWarned {
  return value.event === "EventAccountWarned";
}

export function isEventAccountWarningUpdated(value: EventPayload): value is EventAccountWarningUpdated {
  return value.event === "EventAccountWarningUpdated";
}

export function isEventAccountWarningDeleted(value: EventPayload): value is EventAccountWarningDeleted {
  return value.event === "EventAccountWarningDeleted";
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

export interface RPCRequestRobotModelProviderEmbedTextParams {
  // Provider identifier declared by the plugin manifest.
  provider: string;
  // Text to embed.
  text: string;
}

export interface RobotModelProviderMessage {
  content?: string;
  name?: string;
  role: string;
  tool_call_id?: string;
}

export interface RobotModelProviderTool {
  description?: string;
  name: string;
  parameters?: Record<string, unknown>;
}

export interface RPCRequestRobotModelProviderGenerateParams {
  messages: RobotModelProviderMessage[];
  // Model name within the provider namespace.
  model: string;
  // Provider identifier declared by the plugin manifest.
  provider: string;
  system?: string;
  tools?: RobotModelProviderTool[];
}

export interface RPCRequestRobotModelProviderListModelsParams {
  // Provider identifier declared by the plugin manifest.
  provider: string;
}

export interface RPCRequestRobotModelProviderStructuredPromptParams {
  // Human-readable description of the structured output task.
  description: string;
  // User input or prompt text for the structured output request.
  input: string;
  // Model name within the provider namespace.
  model: string;
  // Provider identifier declared by the plugin manifest.
  provider: string;
  // JSON Schema object the response must satisfy.
  schema: Record<string, unknown>;
}

export interface RPCRequestRobotToolCallParams {
  // Account identifier for the user who invoked the Robot run.
  account_id: string;
  arguments: Record<string, unknown>;
  // Host-generated tool call identifier.
  call_id: string;
  // Tool provider identifier declared by the plugin manifest.
  provider_id: string;
  // Active Robot identifier when the call came from a custom Robot.
  robot_id?: string;
  // Robot session identifier.
  session_id: string;
  // Provider-local tool identifier declared by the plugin manifest.
  tool_id: string;
}


// Request sent by the host to the plugin to provide configuration settings. The params object can contain any key-value pairs defined by the plugin in its manifest `configuration_schema` field and the plugin should validate and apply these settings to its internal state.
// If configuration changes require a plugin to restart, the plugin should cleanly shut down with a zero exit code so that the host can restart it if it is a supervised plugin. If it is an external plugin, the plugin itself is responsible for this behavior based on the plugin's lifecycle design.
export interface RPCRequestConfigure {
  method: "configure";
  id: string;
  jsonrpc: string;
  params: Record<string, unknown>;
}

// Delivers a subscribed Storyden event payload to the plugin.
export interface RPCRequestEvent {
  method: "event";
  id: string;
  jsonrpc: string;
  params: EventPayload;
}

// Health-check request sent by the host to verify plugin responsiveness.
export interface RPCRequestPing {
  method: "ping";
  id: string;
  jsonrpc: string;
  params?: RPCRequestPingParams;
}

// Requests the current list of models exposed by a plugin-backed Robot model provider.
export interface RPCRequestRobotModelProviderListModels {
  method: "robot_model_provider_list_models";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotModelProviderListModelsParams;
}

// Runs one non-streaming model generation through a plugin-backed Robot model provider.
export interface RPCRequestRobotModelProviderGenerate {
  method: "robot_model_provider_generate";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotModelProviderGenerateParams;
}

// Runs one single-shot structured output request through a plugin-backed Robot model provider.
export interface RPCRequestRobotModelProviderStructuredPrompt {
  method: "robot_model_provider_structured_prompt";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotModelProviderStructuredPromptParams;
}

// Requests one text embedding from a plugin-backed Robot model provider.
export interface RPCRequestRobotModelProviderEmbedText {
  method: "robot_model_provider_embed_text";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotModelProviderEmbedTextParams;
}

// Executes a manifest-declared Robot tool provided by a plugin.
export interface RPCRequestRobotToolCall {
  method: "robot_tool_call";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotToolCallParams;
}

export type HostToPluginRequest =
  | RPCRequestConfigure
  | RPCRequestEvent
  | RPCRequestPing
  | RPCRequestRobotModelProviderListModels
  | RPCRequestRobotModelProviderGenerate
  | RPCRequestRobotModelProviderStructuredPrompt
  | RPCRequestRobotModelProviderEmbedText
  | RPCRequestRobotToolCall;

export function isRPCRequestConfigure(value: HostToPluginRequest): value is RPCRequestConfigure {
  return value.method === "configure";
}

export function isRPCRequestEvent(value: HostToPluginRequest): value is RPCRequestEvent {
  return value.method === "event";
}

export function isRPCRequestPing(value: HostToPluginRequest): value is RPCRequestPing {
  return value.method === "ping";
}

export function isRPCRequestRobotModelProviderListModels(value: HostToPluginRequest): value is RPCRequestRobotModelProviderListModels {
  return value.method === "robot_model_provider_list_models";
}

export function isRPCRequestRobotModelProviderGenerate(value: HostToPluginRequest): value is RPCRequestRobotModelProviderGenerate {
  return value.method === "robot_model_provider_generate";
}

export function isRPCRequestRobotModelProviderStructuredPrompt(value: HostToPluginRequest): value is RPCRequestRobotModelProviderStructuredPrompt {
  return value.method === "robot_model_provider_structured_prompt";
}

export function isRPCRequestRobotModelProviderEmbedText(value: HostToPluginRequest): value is RPCRequestRobotModelProviderEmbedText {
  return value.method === "robot_model_provider_embed_text";
}

export function isRPCRequestRobotToolCall(value: HostToPluginRequest): value is RPCRequestRobotToolCall {
  return value.method === "robot_tool_call";
}

export interface HostToPluginResponseError {
  code?: number;
  message?: string;
}

export interface RobotModelProviderModel {
  // Context window in tokens, when known.
  context_window?: number;
  // Human-readable model name.
  display_name?: string;
  // Maximum output tokens, when known.
  max_output_tokens?: number;
  // Model name within the provider namespace.
  name: string;
  // Provider-specific model metadata.
  raw?: Record<string, unknown>;
  // Model release time, when known.
  released_at?: Date;
}

export interface RobotModelProviderToolCall {
  arguments?: Record<string, unknown>;
  id?: string;
  name: string;
}


// Confirms that the configuration was received and applied correctly.
export interface RPCResponseConfigure {
  method: "configure";
  ok: boolean;
}

// Acknowledges that the plugin received the event payload.
export interface RPCResponseEvent {
  method: "event";
  ok: boolean;
}

// Health-check response from the plugin.
export interface RPCResponsePing {
  method: "ping";
  pong: boolean;
  // Optional status message
  status?: string;
  // How long the plugin has been running
  uptime_seconds?: number;
}

// Returns models exposed by a plugin-backed Robot model provider.
export interface RPCResponseRobotModelProviderListModels {
  method: "robot_model_provider_list_models";
  models: RobotModelProviderModel[];
}

// Returns one non-streaming model generation from a plugin-backed Robot model provider.
export interface RPCResponseRobotModelProviderGenerate {
  method: "robot_model_provider_generate";
  content?: string;
  error?: string;
  finish_reason?: string;
  tool_calls?: RobotModelProviderToolCall[];
}

// Returns JSON text generated by a plugin-backed Robot model provider for a JSON Schema constrained structured output request.
export interface RPCResponseRobotModelProviderStructuredPrompt {
  method: "robot_model_provider_structured_prompt";
  // JSON text matching the requested schema.
  content?: string;
  error?: string;
}

// Returns one text embedding from a plugin-backed Robot model provider.
export interface RPCResponseRobotModelProviderEmbedText {
  method: "robot_model_provider_embed_text";
  embedding?: number[];
  error?: string;
}

// Returns a structured result from a plugin-provided Robot tool.
export interface RPCResponseRobotToolCall {
  method: "robot_tool_call";
  error?: string;
  output?: Record<string, unknown>;
}

export type HostToPluginResponseUnion =
  | RPCResponseConfigure
  | RPCResponseEvent
  | RPCResponsePing
  | RPCResponseRobotModelProviderListModels
  | RPCResponseRobotModelProviderGenerate
  | RPCResponseRobotModelProviderStructuredPrompt
  | RPCResponseRobotModelProviderEmbedText
  | RPCResponseRobotToolCall;

export function isRPCResponseConfigure(value: HostToPluginResponseUnion): value is RPCResponseConfigure {
  return value.method === "configure";
}

export function isRPCResponseEvent(value: HostToPluginResponseUnion): value is RPCResponseEvent {
  return value.method === "event";
}

export function isRPCResponsePing(value: HostToPluginResponseUnion): value is RPCResponsePing {
  return value.method === "ping";
}

export function isRPCResponseRobotModelProviderListModels(value: HostToPluginResponseUnion): value is RPCResponseRobotModelProviderListModels {
  return value.method === "robot_model_provider_list_models";
}

export function isRPCResponseRobotModelProviderGenerate(value: HostToPluginResponseUnion): value is RPCResponseRobotModelProviderGenerate {
  return value.method === "robot_model_provider_generate";
}

export function isRPCResponseRobotModelProviderStructuredPrompt(value: HostToPluginResponseUnion): value is RPCResponseRobotModelProviderStructuredPrompt {
  return value.method === "robot_model_provider_structured_prompt";
}

export function isRPCResponseRobotModelProviderEmbedText(value: HostToPluginResponseUnion): value is RPCResponseRobotModelProviderEmbedText {
  return value.method === "robot_model_provider_embed_text";
}

export function isRPCResponseRobotToolCall(value: HostToPluginResponseUnion): value is RPCResponseRobotToolCall {
  return value.method === "robot_tool_call";
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
  // The list of permission names requested for API access. See https://storyden.org/docs/introduction/members/permissions for available values and descriptions.
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
  // This field is used only for Supervised plugins. External plugins can omit it or provide placeholder values.
  args?: string[];
  // The author of the plugin. Must match the pattern `^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`.
  // (NOTE: May change in future.)
  author: string;
  // Platform capabilities this plugin provides. Capabilities extend Storyden with host-defined integration points such as Robot model providers, authentication providers, or Robot tools.
  capabilities?: CapabilityConfig[];
  // The executable or script used to launch your plugin. If your plugin is a binary (Go, Rust, C, etc) then this should be a path to that binary, it's best to put it in the root of your plugin archive like `./myplugin.exe` or `./myplugin`. If your plugin is a script (Python, Node, etc) then this should be the interpreter's `$PATH` executable (e.g. `python` or `node`)  and you should include the script in the `args` field.
  // This field is used only for Supervised plugins. External plugins can provide a placeholder value and it will be ignored by the runtime.
  // Note that Storyden cannot guarantee that the runtime environment defined by the person hosting Storyden will have any language's interpreter on the `$PATH`. If you are running your own instance and building a custom plugin, you should `FROM` the Storyden base image for your deployment so that you know what runtimes are available.
  // If you are distributing a plugin for others to use, we highly recommend that you use a statically compiled language such as Go, Rust or Zig for your plugin so that it's guaranteed to be compatible with any runtime.
  command: string;
  configuration_schema?: ManifestConfigurationSchema;
  // The description of the plugin. Displayed in Plugin Registries as well as in UI of Storyden installations when installed.
  description: string;
  // The list of events the plugin subscribes to and will receive from the host via RPC. Events allow your plugins to react to things that humans or robots do on Storyden.
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

// Optional workspace mount request. Provide either workspace_id or workspace_instance_id, not both.
export interface RPCRequestRobotRunParamsWorkspace {
  // Workspace template to instantiate and mount for this run.
  workspace_id?: string;
  // Existing workspace instance to mount for this run.
  workspace_instance_id?: string;
}

export interface RPCRequestRobotRunParams {
  // Input message for the robot.
  message: string;
  // The Robot to invoke.
  robot_id: string;
  // Optional workspace mount request. Provide either workspace_id or workspace_instance_id, not both.
  workspace?: RPCRequestRobotRunParamsWorkspace;
}


// Request API access credentials provisioned for this plugin.
export interface RPCRequestAccessGet {
  method: "access_get";
  id: string;
  jsonrpc: string;
}

// Request the plugin's current stored configuration from the host.
export interface RPCRequestGetConfig {
  method: "get_config";
  id: string;
  jsonrpc: string;
  params?: RPCRequestGetConfigParams;
}

// Run a one-shot robot invocation and return the assistant's final text response.
export interface RPCRequestRobotRun {
  method: "robot_run";
  id: string;
  jsonrpc: string;
  params: RPCRequestRobotRunParams;
}

export type PluginToHostRequest =
  | RPCRequestAccessGet
  | RPCRequestGetConfig
  | RPCRequestRobotRun;

export function isRPCRequestAccessGet(value: PluginToHostRequest): value is RPCRequestAccessGet {
  return value.method === "access_get";
}

export function isRPCRequestGetConfig(value: PluginToHostRequest): value is RPCRequestGetConfig {
  return value.method === "get_config";
}

export function isRPCRequestRobotRun(value: PluginToHostRequest): value is RPCRequestRobotRun {
  return value.method === "robot_run";
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


// Reason an unattended Robot invocation needs human attention.
export type RobotRunAttentionReason =
  | "missing_input"
  | "needs_approval"
  | "policy_blocked"
  | "tool_unavailable"
  | "error";

export interface RobotRunAttention {
  // Actionable explanation for the user.
  message: string;
  reason: RobotRunAttentionReason;
}


// Overall status of an unattended Robot invocation.
export type RobotRunStatus =
  | "completed"
  | "blocked"
  | "failed";

export interface RobotRunOutput {
  attention?: RobotRunAttention;
  status: RobotRunStatus;
  // Concise user-facing summary of what happened.
  summary: string;
}


// Returns API base URL and bearer access key for authenticated API calls.
export interface RPCResponseAccessGet {
  method: "access_get";
  error?: RPCResponseAccessGetError;
  id: string;
  jsonrpc: string;
  result: RPCResponseAccessGetResult;
}

// Returns current configuration values for this plugin.
export interface RPCResponseGetConfig {
  method: "get_config";
  // Configuration key-value pairs
  config: Record<string, unknown>;
}

// Final result of a one-shot robot invocation.
export interface RPCResponseRobotRun {
  method: "robot_run";
  // Error message if invocation failed.
  error?: string;
  output?: RobotRunOutput;
  // Robot session ID containing the persisted invocation log.
  session_id?: string;
}

export type PluginToHostResponseUnion =
  | RPCResponseAccessGet
  | RPCResponseGetConfig
  | RPCResponseRobotRun;

export function isRPCResponseAccessGet(value: PluginToHostResponseUnion): value is RPCResponseAccessGet {
  return value.method === "access_get";
}

export function isRPCResponseGetConfig(value: PluginToHostResponseUnion): value is RPCResponseGetConfig {
  return value.method === "get_config";
}

export function isRPCResponseRobotRun(value: PluginToHostResponseUnion): value is RPCResponseRobotRun {
  return value.method === "robot_run";
}

export interface PluginToHostResponse {
  error?: PluginToHostResponseError;
  id: string;
  jsonrpc: string;
  result: PluginToHostResponseUnion;
}
