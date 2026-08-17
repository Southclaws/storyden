/* eslint-disable */
/**
 * This file was automatically generated from mcp/schema.yaml
 * DO NOT MODIFY IT BY HAND. Run: node mcp/generate-ts.js
 */

export interface MCPTools {
  ToolContentSearch?: ContentSearch;
  ToolDocumentGet?: DocumentGet;
  ToolDocumentSearch?: DocumentSearch;
  ToolDocumentList?: DocumentList;
  ToolDocumentClose?: DocumentClose;
  ToolToolSearch?: ToolSearch;
  ToolToolGet?: ToolGet;
  ToolToolLoad?: ToolLoad;
  ToolRobotCreate?: RobotCreate;
  ToolRobotList?: RobotList;
  ToolRobotGet?: RobotGet;
  ToolRobotUpdate?: RobotUpdate;
  ToolRobotDelete?: RobotDelete;
  ToolRobotSearch?: RobotSearch;
  ToolToolsetSearch?: ToolsetSearch;
  ToolToolsetLoad?: ToolsetLoad;
  ToolToolsetCreate?: ToolsetCreate;
  ToolToolsetList?: ToolsetList;
  ToolToolsetGet?: ToolsetGet;
  ToolToolsetUpdate?: ToolsetUpdate;
  ToolToolsetDelete?: ToolsetDelete;
  ToolLibraryPageTree?: LibraryPageList;
  ToolLibraryRequestPage?: LibraryRequestPage;
  ToolLibraryPageGet?: LibraryPageGet;
  ToolLibraryPageOpen?: LibraryPageOpen;
  ToolLibraryPageCreate?: CreateLibraryPage;
  ToolLibraryPageUpdate?: UpdateLibraryPage;
  ToolLibrarySearchPages?: LibrarySearchPages;
  ToolLibraryPagePropertySchemaGet?: LibraryPagePropertySchemaGet;
  ToolLibraryPagePropertySchemaUpdate?: LibraryPagePropertySchemaUpdate;
  ToolLibraryPagePropertiesUpdate?: LibraryPagePropertiesUpdate;
  ToolTagList?: TagList;
  ToolLinkCreate?: LinkCreate;
  ToolWebFetch?: WebFetch;
  ToolWebOpen?: WebOpen;
  ToolThreadCreate?: ThreadCreate;
  ToolThreadList?: ThreadList;
  ToolThreadGet?: ThreadGet;
  ToolThreadOpen?: ThreadOpen;
  ToolThreadUpdate?: ThreadUpdate;
  ToolThreadReply?: ThreadReply;
  ToolCategoryList?: CategoryList;
  ToolThreadSearch?: ThreadSearch;
  ToolReplySearch?: ReplySearch;
  ToolPostSearch?: PostSearch;
  ToolMemberSearch?: MemberSearch;
  [k: string]: unknown;
}
/**
 * Search all published content in the community, including library pages, forum threads, and replies.
 */
export interface ContentSearch {
  input: ToolContentSearchInput;
  output: ToolContentSearchOutput;
  [k: string]: unknown;
}
export interface ToolContentSearchInput {
  /**
   * Plain text keyword search string. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces, for example: robot agent automation.
   */
  query: string;
  /**
   * Filter by content types.
   */
  kind?: ("post" | "thread" | "reply" | "node" | "collection" | "profile" | "event")[];
  /**
   * Filter by author handles (usernames). Do not use '@' prefix.
   */
  authors?: string[];
  /**
   * Filter by category names (for forum threads). Category names are case-insensitive.
   */
  categories?: string[];
  /**
   * Filter by tag names. Tags are case-sensitive.
   */
  tags?: string[];
  /**
   * Maximum number of results to return (default 10, max 100)
   */
  max_results?: number;
}
export interface ToolContentSearchOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of search results
   */
  items: SearchedItem[];
}
export interface SearchedItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * Type of content
   */
  kind: string;
  /**
   * URL friendly identifier
   */
  slug: string;
  /**
   * Title or name
   */
  name: string;
  /**
   * Brief description or excerpt
   */
  description?: string;
}
/**
 * Inspect one structural page of the active document or a specified location and make it the current document cursor.
 */
export interface DocumentGet {
  input: ToolDocumentGetInput;
  output: RobotDocumentProjection;
  [k: string]: unknown;
}
export interface ToolDocumentGetInput {
  /**
   * Open document identifier; omit to inspect the active document.
   */
  document_id?: string;
  /**
   * Structural location returned by an outline or search result; omit to continue from that document's current node.
   */
  node_id?: string;
  /**
   * One-indexed structural page to inspect; omit to retain the current page when node_id is also omitted, or use the first page of a specified node.
   */
  page?: number;
}
/**
 * A bounded view into an open document with stable identifiers for further navigation.
 */
export interface RobotDocumentProjection {
  /**
   * Conversation-local document identifier accepted by document navigation tools.
   */
  document_id: string;
  /**
   * Origin kind used to distinguish documents with similar titles.
   */
  source_type: "library_page" | "thread" | "web";
  /**
   * Source identifier or URL used to open this snapshot.
   */
  source_id: string;
  /**
   * Document title for identifying the current result.
   */
  title: string;
  /**
   * Structural location represented by this projection.
   */
  node_id: string;
  /**
   * One-indexed structural page currently shown for this location.
   */
  page: number;
  /**
   * Number of structural pages available for this location.
   */
  total_pages: number;
  /**
   * Previous structural page when one exists.
   */
  previous_page?: number;
  /**
   * Next structural page when one exists.
   */
  next_page?: number;
  /**
   * One-indexed position of the first structural item shown on this page.
   */
  item_start?: number;
  /**
   * One-indexed position of the final structural item shown on this page.
   */
  item_end?: number;
  /**
   * Total number of structural items available at this location.
   */
  total_items?: number;
  /**
   * Bounded plain-text outline, preview, or complete leaf content for this location.
   */
  projection: string;
  /**
   * Whether additional content exists beyond this projection.
   */
  truncated: boolean;
  /**
   * Recommended document navigation action when more detail is needed.
   */
  next_action: string;
}
/**
 * Search structural locations within the active document or selected open document. Whole-term matches are returned before substring-only matches, with document order preserved within each relevance tier.
 */
export interface DocumentSearch {
  input: ToolDocumentSearchInput;
  output: ToolDocumentSearchOutput;
  [k: string]: unknown;
}
export interface ToolDocumentSearchInput {
  /**
   * Plain-text terms that must all occur in each matching structural location.
   */
  query: string;
  /**
   * Open document identifier; omit to search the active document.
   */
  document_id?: string;
  /**
   * Structural locations whose subtrees bound the search; omit to search the whole document.
   *
   * @minItems 1
   */
  node_ids?: [string, ...string[]];
  /**
   * Maximum number of matching locations to return; omit to use ten.
   */
  max_results?: number;
}
export interface ToolDocumentSearchOutput {
  /**
   * Open document searched by this call.
   */
  document_id: string;
  /**
   * Normalized query used to select matches.
   */
  query: string;
  /**
   * Matching locations ranked by whole-term relevance, then document order, with enough context to choose the smallest sufficient location before using document_get.
   */
  matches: DocumentSearchMatch[];
  /**
   * Whether more matching locations exist beyond this result.
   */
  truncated: boolean;
  /**
   * Recommended action for inspecting a selected match or refining the search.
   */
  next_action: string;
}
export interface DocumentSearchMatch {
  /**
   * Structural location accepted by document_get and scoped document_search.
   */
  node_id: string;
  /**
   * Structural role such as heading, paragraph, list item, code, table, or image.
   */
  kind: string;
  /**
   * Heading path that locates this match within the document.
   */
  path: string;
  /**
   * Bounded matching text used to judge relevance before expanding the location.
   */
  preview: string;
}
/**
 * List document snapshots currently available for navigation in this conversation.
 */
export interface DocumentList {
  input: ToolDocumentListInput;
  output: ToolDocumentListOutput;
  [k: string]: unknown;
}
export interface ToolDocumentListInput {}
export interface ToolDocumentListOutput {
  /**
   * Open document snapshots with their navigation identifiers and active status.
   */
  documents: OpenDocument[];
  /**
   * Recommended action for selecting or opening a document.
   */
  next_action: string;
}
export interface OpenDocument {
  /**
   * Conversation-local identifier accepted by document navigation tools.
   */
  document_id: string;
  /**
   * Origin kind used to distinguish documents with similar titles.
   */
  source_type: "library_page" | "thread" | "web";
  /**
   * Source identifier or URL used to open this snapshot.
   */
  source_id: string;
  /**
   * Document title for selecting the intended snapshot.
   */
  title: string;
  /**
   * Whether tools use this document when document_id is omitted.
   */
  active: boolean;
  /**
   * Most recently inspected structural location in this snapshot.
   */
  active_node_id: string;
  /**
   * Most recently inspected structural page at this location.
   */
  page: number;
  /**
   * Number of structural pages at the most recently inspected location.
   */
  total_pages: number;
  /**
   * First structural item shown by the current cursor page.
   */
  item_start?: number;
  /**
   * Final structural item shown by the current cursor page.
   */
  item_end?: number;
  /**
   * Total structural items at the current cursor location.
   */
  total_items?: number;
}
/**
 * Remove an open document snapshot from this conversation without changing its source.
 */
export interface DocumentClose {
  input: ToolDocumentCloseInput;
  output: ToolDocumentCloseOutput;
  [k: string]: unknown;
}
export interface ToolDocumentCloseInput {
  /**
   * Open document identifier; omit to close the active document.
   */
  document_id?: string;
}
export interface ToolDocumentCloseOutput {
  /**
   * Identifier of the removed document snapshot.
   */
  document_id: string;
  /**
   * Number of document snapshots still available for navigation.
   */
  remaining: number;
  /**
   * Newly active document identifier when another snapshot remains.
   */
  active_document_id?: string;
  /**
   * Recommended action after the snapshot was removed.
   */
  next_action: string;
}
/**
 * Search for individual tools that can perform one narrow task. Use tool_get to inspect a candidate's schema before loading it or assigning it to a Robot.
 */
export interface ToolSearch {
  input: ToolToolSearchInput;
  output: ToolToolSearchOutput;
  [k: string]: unknown;
}
export interface ToolToolSearchInput {
  /**
   * Capability, task, or tool name to search for.
   */
  query: string;
  /**
   * Maximum number of ranked candidates to return; omit to use the default.
   */
  max_results?: number;
}
export interface ToolToolSearchOutput {
  /**
   * Ranked tool candidates with enough context to choose which one to inspect.
   */
  tools: RobotToolCatalogueItem[];
  /**
   * Recommended next step for inspecting and activating a selected candidate.
   */
  next_action: string;
}
/**
 * An individually composable tool candidate returned by capability discovery; use its ID with tool_get before choosing or loading it.
 */
export interface RobotToolCatalogueItem {
  /**
   * Stable tool identifier accepted by tool_get, tool_load, and Robot configuration tools.
   */
  id: string;
  /**
   * Human-readable tool title that helps distinguish similar capabilities.
   */
  name: string;
  /**
   * The outcome this tool can achieve, used to decide whether its full schema is worth inspecting.
   */
  description: string;
}
/**
 * Get the full schema, Toolset membership, and runtime preconditions of one tool after finding its ID with tool_search. This is inspection only; it does not load or activate the tool.
 */
export interface ToolGet {
  input: ToolToolGetInput;
  output: ToolToolGetOutput;
  [k: string]: unknown;
}
export interface ToolToolGetInput {
  /**
   * Stable tool ID returned by tool_search.
   */
  id: string;
}
export interface ToolToolGetOutput {
  /**
   * Stable tool identifier accepted by tool_load and Robot configuration tools unless toolset_only is true.
   */
  id: string;
  /**
   * Human-readable title for explaining or confirming the selected capability.
   */
  name: string;
  /**
   * The outcome the tool is designed to achieve and when it should be selected.
   */
  description: string;
  /**
   * Complete JSON Schema for arguments that must be supplied when calling the tool.
   */
  input_schema: {
    [k: string]: unknown;
  };
  /**
   * Complete JSON Schema for the tool result when the tool declares one.
   */
  output_schema?: {
    [k: string]: unknown;
  };
  /**
   * Built-in Toolsets that provide this tool together with their specialist instruction and related capabilities.
   */
  toolsets: string[];
  /**
   * Whether this capability must be loaded or assigned through one of its Toolsets instead of as an individual tool.
   */
  toolset_only: boolean;
  /**
   * Whether a run must pause for human approval before executing this tool.
   */
  requires_confirmation: boolean;
  /**
   * Whether the conversation must have an active Robot workspace before this tool can be loaded or used.
   */
  requires_workspace: boolean;
}
/**
 * Activate one or more individual tools that are not already callable in the current Denbot conversation after inspecting their schemas and runtime preconditions with tool_get. Tools marked toolset_only must be activated through a declared Toolset. Workspace-dependent tools cannot be loaded until the conversation has an active workspace.
 */
export interface ToolLoad {
  input: ToolToolLoadInput;
  output: ToolToolLoadOutput;
  [k: string]: unknown;
}
export interface ToolToolLoadInput {
  /**
   * Stable tool IDs inspected with tool_get that are not currently callable and are not marked toolset_only.
   *
   * @minItems 1
   */
  tools: [string, ...string[]];
}
export interface ToolToolLoadOutput {
  /**
   * Tools activated for subsequent model steps in this conversation.
   */
  loaded: RobotToolCatalogueItem[];
  /**
   * Recommended continuation now that the selected tools are available.
   */
  next_action: string;
}
/**
 * Create a new Robot (agent) with a specific purpose and behavior. Robots are customizable automations that can help users with specific workflows using direct tools, reusable Toolsets, and instructions.
 */
export interface RobotCreate {
  input: ToolRobotCreateInput;
  output: ToolRobotCreateOutput;
  [k: string]: unknown;
}
export interface ToolRobotCreateInput {
  /**
   * The name of the Robot - should be descriptive and help identify its purpose
   */
  name: string;
  /**
   * Optional human-readable description of the Robot's purpose and capabilities
   */
  description: string;
  /**
   * The directive/system prompt that defines how the Robot behaves, what it helps with, and its personality. This is the primary instruction that guides the Robot's behavior.
   */
  playbook: string;
  /**
   * The language model ID in provider/model_name format. If omitted, the system will use a default model configured for Robots.
   */
  model?: string;
  /**
   * Individual tool names that the Robot can use. Tools already provided by an assigned Toolset are removed.
   */
  tools?: string[];
  /**
   * List of reusable Toolset IDs that the Robot can use.
   */
  toolsets?: string[];
}
export interface ToolRobotCreateOutput {
  /**
   * The unique identifier of the created Robot
   */
  id: string;
  /**
   * The name of the Robot
   */
  name: string;
}
/**
 * List all available Robots. Use this to see what Robots have been created and what they do.
 */
export interface RobotList {
  input: ToolRobotListInput;
  output: ToolRobotListOutput;
  [k: string]: unknown;
}
export interface ToolRobotListInput {
  /**
   * Maximum number of Robots to return (default 20)
   */
  limit?: number;
}
export interface ToolRobotListOutput {
  /**
   * List of Robots
   */
  robots: RobotItem[];
  /**
   * Total number of Robots found
   */
  total: number;
}
export interface RobotItem {
  /**
   * Unique identifier
   */
  id: string;
  /**
   * Robot name
   */
  name: string;
  /**
   * Human-readable description of the Robot's purpose
   */
  description?: string;
}
/**
 * Get details of a specific Robot by its ID. Use this to view the full configuration of a Robot.
 */
export interface RobotGet {
  input: ToolRobotGetInput;
  output: ToolRobotGetOutput;
  [k: string]: unknown;
}
export interface ToolRobotGetInput {
  /**
   * The unique identifier of the Robot to retrieve. Must be a valid XID format (20 character alphanumeric string). Example: cq3pqt0q91s73dq8r000
   */
  id: string;
}
export interface ToolRobotGetOutput {
  /**
   * The unique identifier of the Robot
   */
  id: string;
  /**
   * The name of the Robot
   */
  name: string;
  /**
   * Human-readable description of the Robot's purpose
   */
  description?: string;
  /**
   * The Robot's directive
   */
  playbook: string;
  /**
   * The Robot's language model ID in provider/model_name format
   */
  model: string;
  /**
   * List of individual tool names that the Robot can use
   */
  tools: string[];
  /**
   * List of reusable Toolset IDs that the Robot can use
   */
  toolsets: string[];
}
/**
 * Update a Robot's configuration. You can modify its name, description, playbook (directive), direct tools, or Toolsets. Only provide the fields you want to change.
 */
export interface RobotUpdate {
  input: ToolRobotUpdateInput;
  output: ToolRobotUpdateOutput;
  [k: string]: unknown;
}
export interface ToolRobotUpdateInput {
  /**
   * The unique identifier of the Robot to update. Must be a valid XID format (20 character alphanumeric string).
   */
  id: string;
  /**
   * The new name for the Robot
   */
  name?: string;
  /**
   * The new description for the Robot
   */
  description?: string;
  /**
   * The new directive/system prompt for the Robot
   */
  playbook?: string;
  /**
   * The new language model ID in provider/model_name format.
   */
  model?: string;
  /**
   * The new list of individual tool names that the Robot can use.
   */
  tools?: string[];
  /**
   * The new list of reusable Toolset IDs that the Robot can use.
   */
  toolsets?: string[];
}
export interface ToolRobotUpdateOutput {
  /**
   * The unique identifier of the updated Robot
   */
  id: string;
  /**
   * The Robot's name
   */
  name: string;
}
/**
 * Delete a Robot permanently.
 */
export interface RobotDelete {
  input: ToolRobotDeleteInput;
  output: ToolRobotDeleteOutput;
  [k: string]: unknown;
}
export interface ToolRobotDeleteInput {
  /**
   * The unique identifier of the Robot to delete. Must be a valid XID format (20 character alphanumeric string).
   */
  id: string;
}
export interface ToolRobotDeleteOutput {
  /**
   * Whether the deletion was successful
   */
  success: boolean;
  /**
   * The ID of the deleted Robot
   */
  id: string;
}
/**
 * Search specialised Robots by capability before delegating a task.
 */
export interface RobotSearch {
  input: ToolRobotSearchInput;
  output: ToolRobotSearchOutput;
  [k: string]: unknown;
}
export interface ToolRobotSearchInput {
  /**
   * The task or capability a specialised Robot should handle.
   */
  query: string;
  max_results?: number;
}
export interface ToolRobotSearchOutput {
  robots: RobotSearchResult[];
  next_action: string;
}
export interface RobotSearchResult {
  id: string;
  /**
   * Delegation target name.
   */
  delegate_to: string;
  name: string;
  description: string;
}
/**
 * Search reusable Toolsets for a coherent bundle of capabilities and specialist guidance. Use toolset_get to inspect a candidate before loading it or assigning it to a Robot.
 */
export interface ToolsetSearch {
  input: ToolToolsetSearchInput;
  output: ToolToolsetSearchOutput;
  [k: string]: unknown;
}
export interface ToolToolsetSearchInput {
  /**
   * Capability, task, or Toolset name to search for.
   */
  query: string;
  /**
   * Maximum number of ranked candidates to return; omit to use the default.
   */
  max_results?: number;
}
export interface ToolToolsetSearchOutput {
  /**
   * Ranked Toolset candidates with enough context to choose which bundle to inspect.
   */
  toolsets: RobotToolsetCatalogueItem[];
  /**
   * Recommended next step for inspecting and activating a selected candidate.
   */
  next_action: string;
}
/**
 * A reusable capability bundle returned by discovery; use its ID with toolset_get before loading or assigning it.
 */
export interface RobotToolsetCatalogueItem {
  /**
   * Stable Toolset identifier accepted by toolset_get, toolset_load, and Robot configuration tools.
   */
  id: string;
  /**
   * Human-readable Toolset name that identifies the capability bundle.
   */
  name: string;
  /**
   * The jobs this Toolset is designed to handle, used to decide whether its tools and instruction are worth inspecting.
   */
  description: string;
}
/**
 * Activate reusable Toolsets for the current Robot conversation after inspecting their contents and runtime preconditions with toolset_get. Workspace-dependent Toolsets cannot be loaded until the conversation has an active workspace.
 */
export interface ToolsetLoad {
  input: ToolToolsetLoadInput;
  output: ToolToolsetLoadOutput;
  [k: string]: unknown;
}
export interface ToolToolsetLoadInput {
  /**
   * Toolset IDs inspected with toolset_get.
   *
   * @minItems 1
   */
  toolsets: [string, ...string[]];
}
export interface ToolToolsetLoadOutput {
  /**
   * Toolsets activated for subsequent model steps in this conversation.
   */
  loaded: RobotToolsetCatalogueItem[];
  /**
   * Recommended continuation now that the selected Toolsets and their guidance are active.
   */
  next_action: string;
}
/**
 * Create a reusable custom Toolset from available tool identifiers and optional specialised prompt guidance.
 */
export interface ToolsetCreate {
  input: ToolToolsetCreateInput;
  output: ToolToolsetCreateOutput;
  [k: string]: unknown;
}
export interface ToolToolsetCreateInput {
  name: string;
  description: string;
  instruction?: string;
  tools: string[];
}
export interface ToolToolsetCreateOutput {
  id: string;
  name: string;
}
/**
 * List reusable system, custom, and plugin Toolsets.
 */
export interface ToolsetList {
  input: ToolToolsetListInput;
  output: ToolToolsetListOutput;
  [k: string]: unknown;
}
export interface ToolToolsetListInput {}
export interface ToolToolsetListOutput {
  toolsets: RobotToolsetCatalogueItem[];
}
/**
 * Get a Toolset's complete agent-relevant configuration after finding its ID with toolset_search. Returns its tool IDs, specialist instruction, workspace precondition, and whether this caller may edit it.
 */
export interface ToolsetGet {
  input: ToolToolsetGetInput;
  output: ToolToolsetGetOutput;
  [k: string]: unknown;
}
export interface ToolToolsetGetInput {
  /**
   * Stable Toolset ID returned by toolset_search.
   */
  id: string;
}
export interface ToolToolsetGetOutput {
  /**
   * Stable Toolset identifier accepted by toolset_load and Robot configuration tools.
   */
  id: string;
  /**
   * Human-readable name for explaining or confirming the selected capability bundle.
   */
  name: string;
  /**
   * The jobs this Toolset is designed to handle and when the bundle should be selected.
   */
  description: string;
  /**
   * Tool IDs activated or assigned together by this Toolset; use tool_get when an individual schema is needed.
   */
  tools: string[];
  /**
   * Specialist guidance injected while this Toolset is active, including workflow and safety boundaries.
   */
  instruction: string;
  /**
   * Whether the current caller can update or delete this Toolset through management tools.
   */
  editable: boolean;
  /**
   * Whether the conversation must have an active Robot workspace before this Toolset can be loaded or used.
   */
  requires_workspace: boolean;
}
/**
 * Update an authored custom Toolset. System and plugin Toolsets are read-only.
 */
export interface ToolsetUpdate {
  input: ToolToolsetUpdateInput;
  output: ToolToolsetUpdateOutput;
  [k: string]: unknown;
}
export interface ToolToolsetUpdateInput {
  id: string;
  name?: string;
  description?: string;
  instruction?: string;
  tools?: string[];
}
export interface ToolToolsetUpdateOutput {
  id: string;
  name: string;
}
/**
 * Delete an unused authored custom Toolset. Remove it from Robots first.
 */
export interface ToolsetDelete {
  input: ToolToolsetDeleteInput;
  output: ToolToolsetDeleteOutput;
  [k: string]: unknown;
}
export interface ToolToolsetDeleteInput {
  id: string;
}
export interface ToolToolsetDeleteOutput {
  id: string;
  deleted: boolean;
}
/**
 * List library pages by hierarchy and workflow visibility.
 */
export interface LibraryPageList {
  input: ToolLibraryPageTreeInput;
  output: ToolLibraryPageTreeOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageTreeInput {
  /**
   * Maximum depth to traverse (-1 for unlimited, 0 for root only, 1 for root + children, etc.)
   */
  depth?: number;
  /**
   * Limit pages to one workflow visibility.
   */
  visibility?: "draft" | "review" | "published";
}
export interface ToolLibraryPageTreeOutput {
  /**
   * Pages in hierarchical order for selecting a page to inspect.
   */
  pages: LibraryPageTreeNode[];
  /**
   * Number of pages returned.
   */
  results: number;
}
export interface LibraryPageTreeNode {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the page
   */
  id: string;
  /**
   * URL-friendly identifier for the page
   */
  slug: string;
  /**
   * Display name of the page
   */
  name: string;
  /**
   * Brief description of the page
   */
  description: string;
  /**
   * Slug of the parent page (omitted for root pages)
   */
  parent?: string;
  /**
   * Tags associated with this page
   */
  tags?: string[];
  /**
   * Current publishing workflow state.
   */
  visibility: "draft" | "review" | "published";
}
/**
 * Request a Library page selection from the user. When this tool returns, the user has selected the returned Library page and the task should continue using that page.
 */
export interface LibraryRequestPage {
  input: ToolLibraryRequestPageInput;
  output: ToolLibraryRequestPageOutput;
  [k: string]: unknown;
}
export interface ToolLibraryRequestPageInput {}
export interface ToolLibraryRequestPageOutput {
  /**
   * Unique identifier for the selected Library page.
   */
  id: string;
  /**
   * URL-friendly identifier for the selected Library page.
   */
  slug: string;
  /**
   * Display name for the selected Library page.
   */
  name: string;
  /**
   * Brief description of the selected Library page.
   */
  description?: string;
  /**
   * Absolute frontend URL that opens the selected Library page in a browser. Prefer this in Markdown links when presenting the page to users.
   */
  browser_url?: string;
}
/**
 * Retrieve Library page metadata by ID and indicate how to inspect its document when needed.
 */
export interface LibraryPageGet {
  input: ToolLibraryPageGetInput;
  output: ToolLibraryPageGetOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageGetInput {
  /**
   * The unique identifier of the library page to retrieve
   */
  id: string;
}
export interface ToolLibraryPageGetOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the page
   */
  id: string;
  /**
   * URL-friendly identifier for the page
   */
  slug: string;
  /**
   * Display name of the page
   */
  name: string;
  /**
   * Brief description of the page
   */
  description?: string;
  /**
   * Plain-text page body included for external MCP callers.
   */
  content?: string;
  /**
   * Conditional instruction for a Robot to open this page when it is relevant.
   */
  next_action?: string;
  /**
   * Tags associated with this page
   */
  tags: string[];
  /**
   * Slugs of child pages
   */
  child_pages: string[];
  /**
   * Current publishing workflow state.
   */
  visibility: "draft" | "review" | "published";
}
/**
 * Open a Library page as the active document snapshot and return its initial bounded projection.
 */
export interface LibraryPageOpen {
  input: ToolLibraryPageOpenInput;
  output: RobotDocumentProjection;
  [k: string]: unknown;
}
export interface ToolLibraryPageOpenInput {
  /**
   * Library page ID or slug returned by Library discovery tools.
   */
  id: string;
}
/**
 * Create a new page in the library. A slug will be generated automatically if not provided.
 */
export interface CreateLibraryPage {
  input: ToolLibraryPageCreateInput;
  output: ToolLibraryPageCreateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageCreateInput {
  /**
   * The name/title of the page
   */
  name: string;
  /**
   * The unique slug for this page. If not provided, one will be generated from the name.
   */
  slug?: string;
  /**
   * The content of the page in HTML format
   */
  content?: string;
  /**
   * Slug of the parent page. Only include if you already have a parent slug available. Leave empty to create a root-level page.
   */
  parent?: string;
  /**
   * Visibility of the page (default: published)
   */
  visibility?: "published" | "draft";
  /**
   * Optional external URL if this page references a topic on another website
   */
  url?: string;
  /**
   * Optional tags to categorise this page
   */
  tags?: string[];
}
export interface ToolLibraryPageCreateOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the created page
   */
  id: string;
  /**
   * The slug of the created page
   */
  slug: string;
  /**
   * The name of the created page
   */
  name: string;
}
/**
 * Update an existing page in the library. Only provide the fields you want to change.
 */
export interface UpdateLibraryPage {
  input: ToolLibraryPageUpdateInput;
  output: ToolLibraryPageUpdateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageUpdateInput {
  /**
   * The unique identifier of the page to update
   */
  id: string;
  /**
   * The new name/title of the page
   */
  name?: string;
  /**
   * The new URL slug for the page
   */
  slug?: string;
  /**
   * The new content of the page in HTML format
   */
  content?: string;
  /**
   * New visibility of the page
   */
  visibility?: "published" | "draft";
  /**
   * New external URL reference
   */
  url?: string;
  /**
   * New parent page slug. Provide to move the page to a different parent.
   */
  parent?: string;
  /**
   * New tags to categorise this page
   */
  tags?: string[];
}
export interface ToolLibraryPageUpdateOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the updated page
   */
  id: string;
  /**
   * The slug of the updated page
   */
  slug: string;
  /**
   * The name of the updated page
   */
  name: string;
}
/**
 * Search for pages in the library wiki by relevance.
 */
export interface LibrarySearchPages {
  input: ToolLibrarySearchPagesInput;
  output: ToolLibrarySearchPagesOutput;
  [k: string]: unknown;
}
export interface ToolLibrarySearchPagesInput {
  /**
   * Plain text keyword search string. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces, for example: robot agent automation.
   */
  query: string;
}
export interface ToolLibrarySearchPagesOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of matching library pages
   */
  items: LibraryPageSearchItem[];
}
export interface LibraryPageSearchItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * URL-friendly identifier
   */
  slug: string;
  /**
   * Page name/title
   */
  name: string;
  /**
   * Brief description
   */
  description?: string;
  /**
   * Page content excerpt
   */
  content?: string;
}
/**
 * Get the property schema for a library page's children. This returns the schema fields that define what properties child pages can have.
 */
export interface LibraryPagePropertySchemaGet {
  input: ToolLibraryPagePropertySchemaGetInput;
  output: ToolLibraryPagePropertySchemaGetOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPagePropertySchemaGetInput {
  /**
   * The slug or ID of the library page to get the property schema for
   */
  id: string;
}
export interface ToolLibraryPagePropertySchemaGetOutput {
  /**
   * Whether this page has a property schema defined for its children
   */
  has_schema?: boolean;
  /**
   * The schema fields that define properties for child pages
   */
  fields: PropertySchemaField[];
}
export interface PropertySchemaField {
  /**
   * Unique identifier for this field
   */
  id: string;
  /**
   * Display name of the field
   */
  name: string;
  /**
   * Data type of the field
   */
  type: "text" | "number" | "timestamp" | "boolean";
  /**
   * Sort key for ordering fields
   */
  sort?: string;
}
/**
 * Update the property schema for a library page's children. This defines what properties child pages can have. Fields without an ID will be created as new fields. Fields with an ID will be updated. Existing fields not in the list will be removed.
 */
export interface LibraryPagePropertySchemaUpdate {
  input: ToolLibraryPagePropertySchemaUpdateInput;
  output: ToolLibraryPagePropertySchemaUpdateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPagePropertySchemaUpdateInput {
  /**
   * The slug or ID of the parent library page whose children's schema to update
   */
  id: string;
  /**
   * The schema fields to set. Fields with an ID update existing fields, fields without an ID create new ones. Omitted existing fields will be removed.
   */
  fields: PropertySchemaFieldMutation[];
}
export interface PropertySchemaFieldMutation {
  /**
   * Field ID - if provided, updates an existing field. If omitted, creates a new field.
   */
  id?: string;
  /**
   * Display name of the field
   */
  name: string;
  /**
   * Data type of the field
   */
  type: "text" | "number" | "timestamp" | "boolean";
}
export interface ToolLibraryPagePropertySchemaUpdateOutput {
  /**
   * The updated schema fields
   */
  fields: PropertySchemaFieldResult[];
}
export interface PropertySchemaFieldResult {
  /**
   * Unique identifier for this field
   */
  id: string;
  /**
   * Display name of the field
   */
  name: string;
  /**
   * Data type of the field
   */
  type: "text" | "number" | "timestamp" | "boolean";
}
/**
 * Update property values on a library page. The page must have a property schema defined by its parent. Each property value corresponds to a field in the schema.
 */
export interface LibraryPagePropertiesUpdate {
  input: ToolLibraryPagePropertiesUpdateInput;
  output: ToolLibraryPagePropertiesUpdateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPagePropertiesUpdateInput {
  /**
   * The slug or ID of the library page to update properties on
   */
  id: string;
  /**
   * The property values to set
   */
  properties: PropertyValueMutation[];
}
export interface PropertyValueMutation {
  /**
   * The ID of the schema field this property value is for
   */
  field_id: string;
  /**
   * The value to set for this property
   */
  value: string;
}
export interface ToolLibraryPagePropertiesUpdateOutput {
  /**
   * The updated property values
   */
  properties: PropertyValueResult[];
}
export interface PropertyValueResult {
  /**
   * The ID of the schema field
   */
  field_id: string;
  /**
   * The name of the field
   */
  name: string;
  /**
   * Data type of the field
   */
  type?: "text" | "number" | "timestamp" | "boolean";
  /**
   * The current value of this property
   */
  value: string;
}
/**
 * Get a list of all tags on the site or search for tags by name. Tags are labels used to categorize and organize content. Each tag includes its name and the number of items tagged with it.
 */
export interface TagList {
  input: ToolTagListInput;
  output: ToolTagListOutput;
  [k: string]: unknown;
}
export interface ToolTagListInput {
  /**
   * Optional search query to filter tags by name. If not provided, returns all tags.
   */
  query?: string;
}
export interface ToolTagListOutput {
  /**
   * List of tags
   */
  tags: TagItem[];
}
export interface TagItem {
  /**
   * The tag name
   */
  name: string;
  /**
   * Number of items tagged with this tag
   */
  item_count: number;
}
/**
 * Create or update a link in the shared bookmarks list. Fetches the URL and extracts OpenGraph metadata (title, description) and plain text content. Useful for bookmarking websites, articles, or any web content.
 */
export interface LinkCreate {
  input: ToolLinkCreateInput;
  output: ToolLinkCreateOutput;
  [k: string]: unknown;
}
export interface ToolLinkCreateInput {
  /**
   * The URL to create a bookmark for. Must be a valid HTTP/HTTPS URL.
   */
  url: string;
}
export interface ToolLinkCreateOutput {
  /**
   * Unique identifier for the link
   */
  slug: string;
  /**
   * The bookmarked URL
   */
  url: string;
  /**
   * Title extracted from OpenGraph metadata
   */
  opengraph_title?: string;
  /**
   * Description extracted from OpenGraph metadata
   */
  opengraph_description?: string;
  /**
   * Plain text content extracted from the page
   */
  plain_text?: string;
}
/**
 * Fetch a web page for metadata without adding it to Storyden or opening a navigable document snapshot.
 */
export interface WebFetch {
  input: ToolWebFetchInput;
  output: ToolWebFetchOutput;
  [k: string]: unknown;
}
export interface ToolWebFetchInput {
  /**
   * HTTP or HTTPS page URL to fetch.
   */
  url: string;
}
export interface ToolWebFetchOutput {
  /**
   * Page URL supplied to the safe web fetcher.
   */
  url: string;
  /**
   * Page title when exposed by the fetched document metadata.
   */
  title?: string;
  /**
   * Page summary when exposed by the fetched document metadata.
   */
  description?: string;
  /**
   * Resolved favicon URL when exposed by the fetched document.
   */
  favicon_url?: string;
  /**
   * Resolved representative image URL when exposed by the fetched document.
   */
  image_url?: string;
  /**
   * Plain-text page content included for external MCP callers.
   */
  content?: string;
  /**
   * Conditional instruction for a Robot to open this page when it is relevant.
   */
  next_action?: string;
}
/**
 * Fetch a web page as the active document snapshot without storing it in Storyden's link index.
 */
export interface WebOpen {
  input: ToolWebOpenInput;
  output: RobotDocumentProjection;
  [k: string]: unknown;
}
export interface ToolWebOpenInput {
  /**
   * HTTP or HTTPS page URL to fetch and open.
   */
  url: string;
}
/**
 * Create a new discussion thread in the forum. Threads are discussions organized by category with tags for better discovery.
 */
export interface ThreadCreate {
  input: ToolThreadCreateInput;
  output: ToolThreadCreateOutput;
  [k: string]: unknown;
}
export interface ToolThreadCreateInput {
  /**
   * The title of the thread
   */
  title: string;
  /**
   * The content of the thread in HTML format
   */
  body: string;
  /**
   * The category slug for the thread. Use category_list to see available categories.
   */
  category: string;
  /**
   * Thread visibility (default: published)
   */
  visibility?: "published" | "draft";
  /**
   * Optional URL if this thread is about a specific link
   */
  url?: string;
  /**
   * Optional tags for the thread
   */
  tags?: string[];
}
export interface ToolThreadCreateOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the created thread
   */
  id: string;
  /**
   * The thread slug
   */
  slug: string;
  /**
   * The thread title
   */
  title: string;
  /**
   * Creation timestamp
   */
  created_at?: string;
  /**
   * Thread content as plain text
   */
  content?: string;
  /**
   * Thread visibility
   */
  visibility?: string;
  /**
   * Author handle
   */
  author?: string;
  /**
   * Category name
   */
  category?: string;
  tags?: string[];
  /**
   * Associated URL if present
   */
  url?: string;
}
/**
 * List and search discussion threads with filtering and pagination.
 */
export interface ThreadList {
  input: ToolThreadListInput;
  output: ToolThreadListOutput;
  [k: string]: unknown;
}
export interface ToolThreadListInput {
  /**
   * Search query to filter threads
   */
  query?: string;
  /**
   * Filter by visibility
   */
  visibility?: "draft" | "published";
  /**
   * Page number (default: 1)
   */
  page?: number;
}
export interface ToolThreadListOutput {
  threads: ThreadSummary[];
  current_page: number;
  total_pages: number;
  total_results: number;
  /**
   * Next page number if available
   */
  next_page?: number;
}
export interface ThreadSummary {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the thread
   */
  id: string;
  slug: string;
  title: string;
  /**
   * Short excerpt of the thread
   */
  excerpt: string;
  /**
   * Author handle
   */
  author: string;
  /**
   * Category name
   */
  category: string;
}
/**
 * Retrieve discussion thread metadata by ID and indicate how to inspect its document when needed.
 */
export interface ThreadGet {
  input: ToolThreadGetInput;
  output: ToolThreadGetOutput;
  [k: string]: unknown;
}
export interface ToolThreadGetInput {
  /**
   * The unique identifier of the thread to retrieve
   */
  id: string;
  /**
   * Page number for replies (default: 1)
   */
  page?: number;
}
export interface ToolThreadGetOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the thread
   */
  id: string;
  slug: string;
  title: string;
  /**
   * Plain-text thread body included for external MCP callers.
   */
  content?: string;
  /**
   * Conditional instruction for a Robot to open this thread when it is relevant.
   */
  next_action?: string;
  visibility: string;
  /**
   * Author handle
   */
  author: string;
  /**
   * Category name
   */
  category: string;
  tags: string[];
  /**
   * Creation timestamp
   */
  created_at: string;
  /**
   * Associated URL if present
   */
  url?: string;
}
/**
 * Open a discussion thread and all currently visible replies as one active document snapshot, then return its initial bounded projection.
 */
export interface ThreadOpen {
  input: ToolThreadOpenInput;
  output: RobotDocumentProjection;
  [k: string]: unknown;
}
export interface ToolThreadOpenInput {
  /**
   * Thread ID or slug returned by discussion discovery tools.
   */
  id: string;
}
/**
 * Update an existing thread's properties
 */
export interface ThreadUpdate {
  input: ToolThreadUpdateInput;
  output: ToolThreadUpdateOutput;
  [k: string]: unknown;
}
export interface ToolThreadUpdateInput {
  /**
   * The unique identifier of the thread to update
   */
  id: string;
  /**
   * New title for the thread
   */
  title?: string;
  /**
   * New content for the thread in HTML format
   */
  body?: string;
  /**
   * New visibility: published or draft
   */
  visibility?: "published" | "draft";
  /**
   * New tags for the thread
   */
  tags?: string[];
}
export interface ToolThreadUpdateOutput {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the thread
   */
  id: string;
  /**
   * The thread slug
   */
  slug: string;
  /**
   * The thread title
   */
  title: string;
  /**
   * Creation timestamp
   */
  created_at?: string;
  /**
   * Thread content as plain text
   */
  content?: string;
  /**
   * Thread visibility
   */
  visibility?: string;
  /**
   * Author handle
   */
  author?: string;
  /**
   * Category name
   */
  category?: string;
  tags?: string[];
  /**
   * Associated URL if present
   */
  url?: string;
}
/**
 * Add a reply to an existing thread
 */
export interface ThreadReply {
  input: ToolThreadReplyInput;
  output: ToolThreadReplyOutput;
  [k: string]: unknown;
}
export interface ToolThreadReplyInput {
  /**
   * The unique identifier of the thread to reply to
   */
  id: string;
  /**
   * The reply content in HTML format
   */
  body: string;
}
export interface ToolThreadReplyOutput {
  /**
   * Browser URL for the parent thread. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier for the reply
   */
  id: string;
  /**
   * Author handle
   */
  author: string;
  /**
   * Reply content as plain text
   */
  content: string;
  /**
   * Creation timestamp
   */
  created_at: string;
  /**
   * Update timestamp
   */
  updated_at: string;
}
/**
 * List all thread categories with their names and descriptions
 */
export interface CategoryList {
  input: ToolCategoryListInput;
  output: ToolCategoryListOutput;
  [k: string]: unknown;
}
export interface ToolCategoryListInput {}
export interface ToolCategoryListOutput {
  categories: CategoryItem[];
}
export interface CategoryItem {
  /**
   * Category slug
   */
  slug: string;
  /**
   * Category name
   */
  name: string;
  /**
   * Category description
   */
  description?: string;
}
/**
 * Search published forum threads by relevance.
 */
export interface ThreadSearch {
  input: ToolThreadSearchInput;
  output: ToolThreadSearchOutput;
  [k: string]: unknown;
}
export interface ToolThreadSearchInput {
  /**
   * Plain text keyword search string. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces, for example: robot agent automation.
   */
  query: string;
  /**
   * Filter by author handles (usernames). Do not use '@' prefix.
   */
  authors?: string[];
  /**
   * Filter by category names. Category names are case-insensitive.
   */
  categories?: string[];
  /**
   * Filter by tag names. Tags are case-sensitive.
   */
  tags?: string[];
  /**
   * Maximum number of results to return (default 10, max 100)
   */
  max_results?: number;
}
export interface ToolThreadSearchOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of matching threads
   */
  items: ThreadSearchItem[];
}
export interface ThreadSearchItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * URL friendly identifier
   */
  slug: string;
  /**
   * Thread title
   */
  name: string;
  /**
   * Brief excerpt from the thread
   */
  description?: string;
}
/**
 * Search published forum replies by relevance.
 */
export interface ReplySearch {
  input: ToolReplySearchInput;
  output: ToolReplySearchOutput;
  [k: string]: unknown;
}
export interface ToolReplySearchInput {
  /**
   * Plain text keyword search string. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces, for example: robot agent automation.
   */
  query: string;
  /**
   * Filter by author handles (usernames). Do not use '@' prefix.
   */
  authors?: string[];
  /**
   * Filter by tag names. Tags are case-sensitive.
   */
  tags?: string[];
  /**
   * Maximum number of results to return (default 10, max 100)
   */
  max_results?: number;
}
export interface ToolReplySearchOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of matching replies
   */
  items: ReplySearchItem[];
}
export interface ReplySearchItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * URL friendly identifier
   */
  slug: string;
  /**
   * Title of the thread this reply belongs to
   */
  name: string;
  /**
   * Brief excerpt from the reply
   */
  description?: string;
}
/**
 * Search published forum threads and replies together in a single query.
 */
export interface PostSearch {
  input: ToolPostSearchInput;
  output: ToolPostSearchOutput;
  [k: string]: unknown;
}
export interface ToolPostSearchInput {
  /**
   * Plain text keyword search string. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces, for example: robot agent automation.
   */
  query: string;
  /**
   * Filter by author handles (usernames). Do not use '@' prefix.
   */
  authors?: string[];
  /**
   * Filter by category names (applies to threads). Category names are case-insensitive.
   */
  categories?: string[];
  /**
   * Filter by tag names. Tags are case-sensitive.
   */
  tags?: string[];
  /**
   * Maximum number of results to return (default 10, max 100)
   */
  max_results?: number;
}
export interface ToolPostSearchOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of matching threads and replies
   */
  items: PostSearchItem[];
}
export interface PostSearchItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * Content type: thread or reply
   */
  kind: string;
  /**
   * URL friendly identifier
   */
  slug: string;
  /**
   * Thread title or reply excerpt
   */
  name: string;
  /**
   * Brief excerpt
   */
  description?: string;
}
/**
 * Search community members by name or handle.
 */
export interface MemberSearch {
  input: ToolMemberSearchInput;
  output: ToolMemberSearchOutput;
  [k: string]: unknown;
}
export interface ToolMemberSearchInput {
  /**
   * Plain text keyword search string matching name or handle. Boolean operators like OR, AND, and NOT, quoted phrases, and other special query syntax are not supported. Use simple words separated by spaces.
   */
  query: string;
  /**
   * Maximum number of results to return (default 10, max 100)
   */
  max_results?: number;
}
export interface ToolMemberSearchOutput {
  /**
   * Total number of results found
   */
  results: number;
  /**
   * List of matching members
   */
  items: MemberSearchItem[];
}
export interface MemberSearchItem {
  /**
   * Browser URL for this resource. Always present this as a Markdown link when showing results to the user.
   */
  browser_url: string;
  /**
   * Unique identifier
   */
  id: string;
  /**
   * The member's unique handle (username)
   */
  handle: string;
  /**
   * The member's display name
   */
  name: string;
  /**
   * Brief member bio or description
   */
  bio?: string;
}

export type ToolName = "content_search" | "document_get" | "document_search" | "document_list" | "document_close" | "tool_search" | "tool_get" | "tool_load" | "robot_create" | "robot_list" | "robot_get" | "robot_update" | "robot_delete" | "robot_search" | "toolset_search" | "toolset_load" | "toolset_create" | "toolset_list" | "toolset_get" | "toolset_update" | "toolset_delete" | "library_page_list" | "library_request_page" | "library_page_get" | "library_page_open" | "create_library_page" | "update_library_page" | "library_search_pages" | "library_page_property_schema_get" | "library_page_property_schema_update" | "library_page_properties_update" | "tag_list" | "link_create" | "web_fetch" | "web_open" | "thread_create" | "thread_list" | "thread_get" | "thread_open" | "thread_update" | "thread_reply" | "category_list" | "thread_search" | "reply_search" | "post_search" | "member_search";

export const TOOL_NAMES = ["content_search", "document_get", "document_search", "document_list", "document_close", "tool_search", "tool_get", "tool_load", "robot_create", "robot_list", "robot_get", "robot_update", "robot_delete", "robot_search", "toolset_search", "toolset_load", "toolset_create", "toolset_list", "toolset_get", "toolset_update", "toolset_delete", "library_page_list", "library_request_page", "library_page_get", "library_page_open", "create_library_page", "update_library_page", "library_search_pages", "library_page_property_schema_get", "library_page_property_schema_update", "library_page_properties_update", "tag_list", "link_create", "web_fetch", "web_open", "thread_create", "thread_list", "thread_get", "thread_open", "thread_update", "thread_reply", "category_list", "thread_search", "reply_search", "post_search", "member_search"] as const;

export type ToolInputMap = {
  "content_search": ToolContentSearchInput;
  "document_get": ToolDocumentGetInput;
  "document_search": ToolDocumentSearchInput;
  "document_list": ToolDocumentListInput;
  "document_close": ToolDocumentCloseInput;
  "tool_search": ToolToolSearchInput;
  "tool_get": ToolToolGetInput;
  "tool_load": ToolToolLoadInput;
  "robot_create": ToolRobotCreateInput;
  "robot_list": ToolRobotListInput;
  "robot_get": ToolRobotGetInput;
  "robot_update": ToolRobotUpdateInput;
  "robot_delete": ToolRobotDeleteInput;
  "robot_search": ToolRobotSearchInput;
  "toolset_search": ToolToolsetSearchInput;
  "toolset_load": ToolToolsetLoadInput;
  "toolset_create": ToolToolsetCreateInput;
  "toolset_list": ToolToolsetListInput;
  "toolset_get": ToolToolsetGetInput;
  "toolset_update": ToolToolsetUpdateInput;
  "toolset_delete": ToolToolsetDeleteInput;
  "library_page_list": ToolLibraryPageTreeInput;
  "library_request_page": ToolLibraryRequestPageInput;
  "library_page_get": ToolLibraryPageGetInput;
  "library_page_open": ToolLibraryPageOpenInput;
  "create_library_page": ToolLibraryPageCreateInput;
  "update_library_page": ToolLibraryPageUpdateInput;
  "library_search_pages": ToolLibrarySearchPagesInput;
  "library_page_property_schema_get": ToolLibraryPagePropertySchemaGetInput;
  "library_page_property_schema_update": ToolLibraryPagePropertySchemaUpdateInput;
  "library_page_properties_update": ToolLibraryPagePropertiesUpdateInput;
  "tag_list": ToolTagListInput;
  "link_create": ToolLinkCreateInput;
  "web_fetch": ToolWebFetchInput;
  "web_open": ToolWebOpenInput;
  "thread_create": ToolThreadCreateInput;
  "thread_list": ToolThreadListInput;
  "thread_get": ToolThreadGetInput;
  "thread_open": ToolThreadOpenInput;
  "thread_update": ToolThreadUpdateInput;
  "thread_reply": ToolThreadReplyInput;
  "category_list": ToolCategoryListInput;
  "thread_search": ToolThreadSearchInput;
  "reply_search": ToolReplySearchInput;
  "post_search": ToolPostSearchInput;
  "member_search": ToolMemberSearchInput;
};

export type ToolOutputMap = {
  "content_search": ToolContentSearchOutput;
  "document_get": RobotDocumentProjection;
  "document_search": ToolDocumentSearchOutput;
  "document_list": ToolDocumentListOutput;
  "document_close": ToolDocumentCloseOutput;
  "tool_search": ToolToolSearchOutput;
  "tool_get": ToolToolGetOutput;
  "tool_load": ToolToolLoadOutput;
  "robot_create": ToolRobotCreateOutput;
  "robot_list": ToolRobotListOutput;
  "robot_get": ToolRobotGetOutput;
  "robot_update": ToolRobotUpdateOutput;
  "robot_delete": ToolRobotDeleteOutput;
  "robot_search": ToolRobotSearchOutput;
  "toolset_search": ToolToolsetSearchOutput;
  "toolset_load": ToolToolsetLoadOutput;
  "toolset_create": ToolToolsetCreateOutput;
  "toolset_list": ToolToolsetListOutput;
  "toolset_get": ToolToolsetGetOutput;
  "toolset_update": ToolToolsetUpdateOutput;
  "toolset_delete": ToolToolsetDeleteOutput;
  "library_page_list": ToolLibraryPageTreeOutput;
  "library_request_page": ToolLibraryRequestPageOutput;
  "library_page_get": ToolLibraryPageGetOutput;
  "library_page_open": RobotDocumentProjection;
  "create_library_page": ToolLibraryPageCreateOutput;
  "update_library_page": ToolLibraryPageUpdateOutput;
  "library_search_pages": ToolLibrarySearchPagesOutput;
  "library_page_property_schema_get": ToolLibraryPagePropertySchemaGetOutput;
  "library_page_property_schema_update": ToolLibraryPagePropertySchemaUpdateOutput;
  "library_page_properties_update": ToolLibraryPagePropertiesUpdateOutput;
  "tag_list": ToolTagListOutput;
  "link_create": ToolLinkCreateOutput;
  "web_fetch": ToolWebFetchOutput;
  "web_open": RobotDocumentProjection;
  "thread_create": ToolThreadCreateOutput;
  "thread_list": ToolThreadListOutput;
  "thread_get": ToolThreadGetOutput;
  "thread_open": RobotDocumentProjection;
  "thread_update": ToolThreadUpdateOutput;
  "thread_reply": ToolThreadReplyOutput;
  "category_list": ToolCategoryListOutput;
  "thread_search": ToolThreadSearchOutput;
  "reply_search": ToolReplySearchOutput;
  "post_search": ToolPostSearchOutput;
  "member_search": ToolMemberSearchOutput;
};
export type StorydenTools = {
  "content_search": {
    input: ToolContentSearchInput;
    output: ToolContentSearchOutput;
  };
  "document_get": {
    input: ToolDocumentGetInput;
    output: RobotDocumentProjection;
  };
  "document_search": {
    input: ToolDocumentSearchInput;
    output: ToolDocumentSearchOutput;
  };
  "document_list": {
    input: ToolDocumentListInput;
    output: ToolDocumentListOutput;
  };
  "document_close": {
    input: ToolDocumentCloseInput;
    output: ToolDocumentCloseOutput;
  };
  "tool_search": {
    input: ToolToolSearchInput;
    output: ToolToolSearchOutput;
  };
  "tool_get": {
    input: ToolToolGetInput;
    output: ToolToolGetOutput;
  };
  "tool_load": {
    input: ToolToolLoadInput;
    output: ToolToolLoadOutput;
  };
  "robot_create": {
    input: ToolRobotCreateInput;
    output: ToolRobotCreateOutput;
  };
  "robot_list": {
    input: ToolRobotListInput;
    output: ToolRobotListOutput;
  };
  "robot_get": {
    input: ToolRobotGetInput;
    output: ToolRobotGetOutput;
  };
  "robot_update": {
    input: ToolRobotUpdateInput;
    output: ToolRobotUpdateOutput;
  };
  "robot_delete": {
    input: ToolRobotDeleteInput;
    output: ToolRobotDeleteOutput;
  };
  "robot_search": {
    input: ToolRobotSearchInput;
    output: ToolRobotSearchOutput;
  };
  "toolset_search": {
    input: ToolToolsetSearchInput;
    output: ToolToolsetSearchOutput;
  };
  "toolset_load": {
    input: ToolToolsetLoadInput;
    output: ToolToolsetLoadOutput;
  };
  "toolset_create": {
    input: ToolToolsetCreateInput;
    output: ToolToolsetCreateOutput;
  };
  "toolset_list": {
    input: ToolToolsetListInput;
    output: ToolToolsetListOutput;
  };
  "toolset_get": {
    input: ToolToolsetGetInput;
    output: ToolToolsetGetOutput;
  };
  "toolset_update": {
    input: ToolToolsetUpdateInput;
    output: ToolToolsetUpdateOutput;
  };
  "toolset_delete": {
    input: ToolToolsetDeleteInput;
    output: ToolToolsetDeleteOutput;
  };
  "library_page_list": {
    input: ToolLibraryPageTreeInput;
    output: ToolLibraryPageTreeOutput;
  };
  "library_request_page": {
    input: ToolLibraryRequestPageInput;
    output: ToolLibraryRequestPageOutput;
  };
  "library_page_get": {
    input: ToolLibraryPageGetInput;
    output: ToolLibraryPageGetOutput;
  };
  "library_page_open": {
    input: ToolLibraryPageOpenInput;
    output: RobotDocumentProjection;
  };
  "create_library_page": {
    input: ToolLibraryPageCreateInput;
    output: ToolLibraryPageCreateOutput;
  };
  "update_library_page": {
    input: ToolLibraryPageUpdateInput;
    output: ToolLibraryPageUpdateOutput;
  };
  "library_search_pages": {
    input: ToolLibrarySearchPagesInput;
    output: ToolLibrarySearchPagesOutput;
  };
  "library_page_property_schema_get": {
    input: ToolLibraryPagePropertySchemaGetInput;
    output: ToolLibraryPagePropertySchemaGetOutput;
  };
  "library_page_property_schema_update": {
    input: ToolLibraryPagePropertySchemaUpdateInput;
    output: ToolLibraryPagePropertySchemaUpdateOutput;
  };
  "library_page_properties_update": {
    input: ToolLibraryPagePropertiesUpdateInput;
    output: ToolLibraryPagePropertiesUpdateOutput;
  };
  "tag_list": {
    input: ToolTagListInput;
    output: ToolTagListOutput;
  };
  "link_create": {
    input: ToolLinkCreateInput;
    output: ToolLinkCreateOutput;
  };
  "web_fetch": {
    input: ToolWebFetchInput;
    output: ToolWebFetchOutput;
  };
  "web_open": {
    input: ToolWebOpenInput;
    output: RobotDocumentProjection;
  };
  "thread_create": {
    input: ToolThreadCreateInput;
    output: ToolThreadCreateOutput;
  };
  "thread_list": {
    input: ToolThreadListInput;
    output: ToolThreadListOutput;
  };
  "thread_get": {
    input: ToolThreadGetInput;
    output: ToolThreadGetOutput;
  };
  "thread_open": {
    input: ToolThreadOpenInput;
    output: RobotDocumentProjection;
  };
  "thread_update": {
    input: ToolThreadUpdateInput;
    output: ToolThreadUpdateOutput;
  };
  "thread_reply": {
    input: ToolThreadReplyInput;
    output: ToolThreadReplyOutput;
  };
  "category_list": {
    input: ToolCategoryListInput;
    output: ToolCategoryListOutput;
  };
  "thread_search": {
    input: ToolThreadSearchInput;
    output: ToolThreadSearchOutput;
  };
  "reply_search": {
    input: ToolReplySearchInput;
    output: ToolReplySearchOutput;
  };
  "post_search": {
    input: ToolPostSearchInput;
    output: ToolPostSearchOutput;
  };
  "member_search": {
    input: ToolMemberSearchInput;
    output: ToolMemberSearchOutput;
  };
};
