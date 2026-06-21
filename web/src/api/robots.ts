/* eslint-disable */
/**
 * This file was automatically generated from mcp/schema.yaml
 * DO NOT MODIFY IT BY HAND. Run: node mcp/generate-ts.js
 */

export interface MCPTools {
  ToolContentSearch?: ContentSearch;
  ToolRobotSwitch?: RobotSwitch;
  ToolSystemRobotToolCatalog?: SystemRobotToolCatalog;
  ToolRobotCreate?: RobotCreate;
  ToolRobotList?: RobotList;
  ToolRobotGet?: RobotGet;
  ToolRobotUpdate?: RobotUpdate;
  ToolRobotDelete?: RobotDelete;
  ToolLibraryPageTree?: LibraryPageList;
  ToolLibraryRequestPage?: LibraryRequestPage;
  ToolLibraryPageGet?: GetLibraryPage;
  ToolLibraryPageCreate?: CreateLibraryPage;
  ToolLibraryPageUpdate?: UpdateLibraryPage;
  ToolLibrarySearchPages?: LibrarySearchPages;
  ToolLibraryPagePropertySchemaGet?: LibraryPagePropertySchemaGet;
  ToolLibraryPagePropertySchemaUpdate?: LibraryPagePropertySchemaUpdate;
  ToolLibraryPagePropertiesUpdate?: LibraryPagePropertiesUpdate;
  ToolTagList?: TagList;
  ToolLinkCreate?: LinkCreate;
  ToolThreadCreate?: ThreadCreate;
  ToolThreadList?: ThreadList;
  ToolThreadGet?: ThreadGet;
  ToolThreadUpdate?: ThreadUpdate;
  ToolThreadReply?: ThreadReply;
  ToolCategoryList?: CategoryList;
  ToolThreadSearch?: ThreadSearch;
  ToolReplySearch?: ReplySearch;
  ToolPostSearch?: PostSearch;
  ToolMemberSearch?: MemberSearch;
  DatagraphItemRef?: DatagraphItemRef;
  RobotChatContext?: RobotChatContext;
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
 * Switch the current conversation to a different Robot (agent). Use this when the user wants to talk to a different specialized agent or when a different agent would be better suited to help with the user's request.
 */
export interface RobotSwitch {
  input: ToolRobotSwitchInput;
  output: ToolRobotSwitchOutput;
  [k: string]: unknown;
}
export interface ToolRobotSwitchInput {
  /**
   * The ID of the Robot (agent) to switch to. Must be a valid XID format (20 character alphanumeric string). Use robot_list to see available Robot IDs.
   */
  robot_id: string;
}
export interface ToolRobotSwitchOutput {
  /**
   * Whether the agent switch was successful
   */
  success: boolean;
  /**
   * The ID of the robot that was switched to
   */
  robot_id: string;
}
/**
 * Returns a list of tools that may be assigned to a Robot during Robot creation or update. Use this when selecting capabilities assigned to a Robot. This catalog describes Robot capabilities, not this conversation's capabilities.
 *
 */
export interface SystemRobotToolCatalog {
  input: ToolSystemRobotToolCatalogInput;
  output: ToolSystemRobotToolCatalogOutput;
  [k: string]: unknown;
}
export interface ToolSystemRobotToolCatalogInput {}
/**
 * List of the tools available in the catalog that you can assign to Robots, with their names and descriptions
 */
export interface ToolSystemRobotToolCatalogOutput {
  tools?: ToolInfo[];
  [k: string]: unknown;
}
export interface ToolInfo {
  /**
   * The tool name identifier
   */
  name: string;
  /**
   * Brief description of what the tool does
   */
  description: string;
}
/**
 * Create a new Robot (agent) with a specific purpose and behavior. Robots are customizable automations that can help users with specific workflows using tailored tools and instructions.
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
   * List of tool names that the Robot can use.
   */
  tools?: string[];
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
   * List of tool names that the Robot can use
   */
  tools: string[];
}
/**
 * Update a Robot's configuration. You can modify its name, description, playbook (directive), or tools. Only provide the fields you want to change.
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
   * The new list of tool names that the Robot can use.
   */
  tools?: string[];
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
 * Get the full tree structure of pages in the library. Returns a hierarchical view of all wiki pages showing their parent-child relationships.
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
}
export interface ToolLibraryPageTreeOutput {
  /**
   * List of pages in tree structure
   */
  pages: LibraryPageTreeNode[];
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
 * Get detailed information about a specific library page by its slug.
 */
export interface GetLibraryPage {
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
   * Tags associated with this page
   */
  tags: string[];
  /**
   * Slugs of child pages
   */
  child_pages: string[];
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
 * Get a specific thread with its content. Returns the thread details including author, content, and category information.
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
   * Thread content as plain text
   */
  content: string;
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
export interface DatagraphItemRef {
  /**
   * Unique identifier of the datagraph item
   */
  id: string;
  /**
   * URL-friendly slug of the datagraph item
   */
  slug: string;
  /**
   * Type of datagraph item (thread, node, profile, collection, etc.)
   */
  kind: "post" | "thread" | "reply" | "node" | "collection" | "profile" | "event";
}
export interface RobotChatContext {
  datagraph_item?: DatagraphItemRef1;
  /**
   * Human-readable page type if not viewing a specific datagraph item. Examples: 'Index page', 'Settings page', 'Admin page', 'Search page'. This is free-form text since the backend doesn't know about frontend routes.
   */
  page_type?: string;
}
/**
 * Optional reference to a datagraph item if the user is viewing one (e.g., a thread, library page, profile)
 */
export interface DatagraphItemRef1 {
  /**
   * Unique identifier of the datagraph item
   */
  id: string;
  /**
   * URL-friendly slug of the datagraph item
   */
  slug: string;
  /**
   * Type of datagraph item (thread, node, profile, collection, etc.)
   */
  kind: "post" | "thread" | "reply" | "node" | "collection" | "profile" | "event";
}

export type ToolName = "content_search" | "robot_switch" | "system_robot_tool_catalog" | "robot_create" | "robot_list" | "robot_get" | "robot_update" | "robot_delete" | "library_page_list" | "library_request_page" | "get_library_page" | "create_library_page" | "update_library_page" | "library_search_pages" | "library_page_property_schema_get" | "library_page_property_schema_update" | "library_page_properties_update" | "tag_list" | "link_create" | "thread_create" | "thread_list" | "thread_get" | "thread_update" | "thread_reply" | "category_list" | "thread_search" | "reply_search" | "post_search" | "member_search";

export const TOOL_NAMES = ["content_search", "robot_switch", "system_robot_tool_catalog", "robot_create", "robot_list", "robot_get", "robot_update", "robot_delete", "library_page_list", "library_request_page", "get_library_page", "create_library_page", "update_library_page", "library_search_pages", "library_page_property_schema_get", "library_page_property_schema_update", "library_page_properties_update", "tag_list", "link_create", "thread_create", "thread_list", "thread_get", "thread_update", "thread_reply", "category_list", "thread_search", "reply_search", "post_search", "member_search"] as const;

export type ToolInputMap = {
  "content_search": ToolContentSearchInput;
  "robot_switch": ToolRobotSwitchInput;
  "system_robot_tool_catalog": ToolSystemRobotToolCatalogInput;
  "robot_create": ToolRobotCreateInput;
  "robot_list": ToolRobotListInput;
  "robot_get": ToolRobotGetInput;
  "robot_update": ToolRobotUpdateInput;
  "robot_delete": ToolRobotDeleteInput;
  "library_page_list": ToolLibraryPageTreeInput;
  "library_request_page": ToolLibraryRequestPageInput;
  "get_library_page": ToolLibraryPageGetInput;
  "create_library_page": ToolLibraryPageCreateInput;
  "update_library_page": ToolLibraryPageUpdateInput;
  "library_search_pages": ToolLibrarySearchPagesInput;
  "library_page_property_schema_get": ToolLibraryPagePropertySchemaGetInput;
  "library_page_property_schema_update": ToolLibraryPagePropertySchemaUpdateInput;
  "library_page_properties_update": ToolLibraryPagePropertiesUpdateInput;
  "tag_list": ToolTagListInput;
  "link_create": ToolLinkCreateInput;
  "thread_create": ToolThreadCreateInput;
  "thread_list": ToolThreadListInput;
  "thread_get": ToolThreadGetInput;
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
  "robot_switch": ToolRobotSwitchOutput;
  "system_robot_tool_catalog": ToolSystemRobotToolCatalogOutput;
  "robot_create": ToolRobotCreateOutput;
  "robot_list": ToolRobotListOutput;
  "robot_get": ToolRobotGetOutput;
  "robot_update": ToolRobotUpdateOutput;
  "robot_delete": ToolRobotDeleteOutput;
  "library_page_list": ToolLibraryPageTreeOutput;
  "library_request_page": ToolLibraryRequestPageOutput;
  "get_library_page": ToolLibraryPageGetOutput;
  "create_library_page": ToolLibraryPageCreateOutput;
  "update_library_page": ToolLibraryPageUpdateOutput;
  "library_search_pages": ToolLibrarySearchPagesOutput;
  "library_page_property_schema_get": ToolLibraryPagePropertySchemaGetOutput;
  "library_page_property_schema_update": ToolLibraryPagePropertySchemaUpdateOutput;
  "library_page_properties_update": ToolLibraryPagePropertiesUpdateOutput;
  "tag_list": ToolTagListOutput;
  "link_create": ToolLinkCreateOutput;
  "thread_create": ToolThreadCreateOutput;
  "thread_list": ToolThreadListOutput;
  "thread_get": ToolThreadGetOutput;
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
  "robot_switch": {
    input: ToolRobotSwitchInput;
    output: ToolRobotSwitchOutput;
  };
  "system_robot_tool_catalog": {
    input: ToolSystemRobotToolCatalogInput;
    output: ToolSystemRobotToolCatalogOutput;
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
  "library_page_list": {
    input: ToolLibraryPageTreeInput;
    output: ToolLibraryPageTreeOutput;
  };
  "library_request_page": {
    input: ToolLibraryRequestPageInput;
    output: ToolLibraryRequestPageOutput;
  };
  "get_library_page": {
    input: ToolLibraryPageGetInput;
    output: ToolLibraryPageGetOutput;
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
