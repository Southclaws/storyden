# Storyden Tool Schema Guidelines

These guidelines define how **LLM-facing tools** must be designed, named, and evolved in Storyden.
They exist to ensure **stable tool selection, low hallucination rates, and long-term compatibility** across Codex, Claude, Cursor, and similar agents.

These rules apply to **all new tools** and **any modifications to existing tools**.

---

## 1. Tool Naming

### 1.1 Naming format

- Tools MUST be named using **noun_verb** ordering.
- Use **snake_case** consistently.

**Correct**

```
robot_create
robot_list
library_page_get
thread_reply
```

**Incorrect**

```
create_robot
getLibraryPage
robotCreate
```

### 1.2 Resource grouping

- Tool names MUST group by resource prefix.
- All tools operating on the same resource MUST share the same prefix.

**Example**

```
thread_create
thread_list
thread_get
thread_update
thread_reply
```

This grouping is intentional and MUST be preserved.

---

### 1.3 Tool Description Guidelines

Tool descriptions are the **primary semantic signal** used by modern language models when selecting and executing tools.
Poor descriptions cause mis-selection more often than poor naming or schema design.

These rules MUST be followed for all tool descriptions.

---

#### 1. Describe the Intent, Not the Mechanics

Descriptions MUST explain **what the tool does**, not how it is implemented.

**Correct**

> “Create a new forum thread in the specified category.”

**Incorrect**

> “Inserts a new thread row into the database.”

---

#### 2. Use a Single, Literal Sentence

- Descriptions MUST be **one sentence**
- Avoid conjunctions (`and`, `or`) unless strictly necessary
- No marketing language, no qualifiers

**Correct**

> “Retrieve a forum thread by its ID.”

**Incorrect**

> “Retrieve and display a forum thread with replies and metadata.”

---

#### 3. Start With a Verb Phrase (Even if the Tool Name Is Noun-Verb)

This helps intent matching without renaming tools.

**Correct**

> “Update an existing library page.”

**Incorrect**

> “Library page update operation.”

---

#### 4. State the Primary Object Explicitly

Always name the resource being operated on.

**Correct**

> “List robots available to the current user.”

**Incorrect**

> “List available items.”

---

#### 5. Mention Required Preconditions Only if Non-Obvious

Only mention constraints if violating them would cause confusion.

**Correct**

> “Reply to an existing forum thread.”

**Incorrect**

> “Reply to an existing forum thread if the user is authenticated and the thread is not locked.”

Auth, permissions, and validation are assumed unless unusual.

---

#### 6. Do NOT Restate Schema Information

Descriptions MUST NOT repeat:

- enum values
- required fields
- min/max ranges
- validation rules

That information already exists in the schema.

**Correct**

```
description: "Filter results by content type."
enum: ["thread","node","profile"]
```

**Incorrect**

```
description: "Filter by content type such as thread, node, or profile."
```

---

#### 7. Avoid Output Descriptions Unless Necessary

Models infer output structure from the schema.

Only describe outputs if:

- there is a non-obvious side effect
- ordering is important
- the result is not directly “the thing you asked for”

**Correct**

> “Search the knowledge base and return matching content ordered by relevance.”

**Incorrect**

> “Returns an array of objects containing id, slug, name, and description.”

---

#### 8. Avoid Ambiguous Verbs

Prefer concrete verbs over vague ones.

**Prefer**

- create
- retrieve
- update
- delete
- list
- search
- reply
- switch

**Avoid**

- handle
- process
- manage
- operate
- execute

---

#### 9. Never Describe Error Handling or Partial Success

Descriptions MUST assume:

- valid input
- successful execution

Error cases are handled by the runtime, not the description.

**Incorrect**

> “Returns an error if the thread does not exist.”

---

#### 10. Tool Descriptions Are Not User-Facing Copy

Descriptions MUST NOT:

- explain UI behavior
- mention frontend concepts
- reference buttons, pages, or navigation

**Incorrect**

> “Gets the thread so it can be shown on the forum page.”

---

### Reference Examples (Storyden)

**Good**

```
thread_get:
  description: "Retrieve a forum thread by its ID."
```

```
library_page_list:
  description: "List library pages starting from a parent page."
```

```
search:
  description: "Search the Storyden knowledge base for relevant content."
```

**Bad**

```
search:
  description: "Search for posts, threads, replies, profiles, and other content types using a full-text search engine."
```

---

### Summary Rules (For Agents)

When writing or auditing a tool description:

- [ ] One sentence
- [ ] Starts with a verb phrase
- [ ] Names the primary resource
- [ ] States intent, not implementation
- [ ] Does not repeat schema information
- [ ] Does not mention errors or UI
- [ ] Uses concrete verbs

If any rule is violated, the description MUST be rewritten.

---

## 2. Parameter Naming

### 2.1 Array parameters use singular names

Array parameters MUST use **singular** names, not plural.

**Correct**

```yaml
kind:
  type: array
  items:
    type: string
```

**Incorrect**

```yaml
kinds:
  type: array
  items:
    type: string
```

### 2.2 Empty input schemas

Tools with no input parameters MUST still have an explicit empty `properties` object.

**Correct**

```yaml
ToolFooInput:
  type: object
  properties: {}
  additionalProperties: false
```

**Incorrect**

```yaml
ToolFooInput:
  type: object
  additionalProperties: false
```

---

## 3. Identifiers: IDs vs Slugs

### 3.1 Input arguments

- **IDs are REQUIRED** for all tool input arguments that reference existing resources.

**Correct**

```json
{ "thread_id": "abc123" }
```

**Incorrect**

```json
{ "slug": "my-thread-title" }
```

### 3.2 Output fields

- Tool outputs SHOULD include:

  - `id`
  - `slug` (if applicable)

This allows both reliable follow-up calls (ID) and human-readable output (slug).

### 3.3 XID format

Storyden uses **XID** as the identifier format:

- 20 characters
- Base32 alphabet: `0-9a-v`
- Sortable by creation time

**Validation pattern**

```yaml
pattern: "^[0-9a-v]{20}$"
```

**Example**

```
cq3pqt0q91s73dq8r000
```

### 3.4 Mark system (backend implementation detail)

The backend's `mark.QueryKey` system accepts multiple input formats:

- Pure XID: `cq3pqt0q91s73dq8r000`
- Mark (XID + slug): `cq3pqt0q91s73dq8r000-my-page-title`
- Pure slug: `my-page-title`

Tool schemas should only accept IDs, but the backend is flexible enough to parse marks if needed.

---

## 4. SDRs (Storyden References)

- SDRs are **NOT exposed** to LLM-facing tools.
- SDRs are internal to backend/frontend graph resolution.
- Tools MUST NOT accept or return SDRs.

This is intentional and MUST NOT be changed without a separate design review.

---

## 5. Tool Descriptions

### 5.1 Descriptions are required

Every tool MUST have:

- A short, literal description of what it does.
- No redundancy with schema constraints.

**Do NOT restate enums, ranges, or validation rules already expressed in schema.**

**Correct**

```
description: "Filter by content types."
enum: ["thread","node","profile"]
```

**Incorrect**

```
description: "Filter by content types: thread, node, profile."
```

---

## 6. Input Schema Rules

### 6.1 Explicit schemas

- All inputs MUST be explicitly typed.
- `any`, free-form objects, or unbounded maps are NOT allowed.

### 6.2 Fail-fast validation

- Invalid input MUST cause the tool to fail.
- Partial success is NOT allowed for LLM-facing tools.
- Do NOT return “warnings” or “partial results”.

This prevents ambiguity in tool execution.

---

## 7. Output Schema Rules

### 7.1 Stable shape

- Tool outputs MUST have a deterministic, stable structure.
- Avoid polymorphic or shape-shifting responses.

### 7.2 Lists

List outputs MUST include:

- `<resource name plural>: []`
- A count or pagination indicator when applicable.

For example:

```json
{
    "threads": [ ... ]
}
```

### 7.3 Search-style tools

Search tools MUST:

- Return results in **descending relevance order**
- Allow pagination via a 1-indexed `page` parameter
- NOT expose internal engine details

Score values are optional; ordering is authoritative.

---

## 8. Pagination

- Prefer **page-based pagination** unless cursor semantics already exist.
- Page size may be fixed and hidden from the model.
- Expose:

  - `page` (input)
  - `has_more` or `total_pages` (output)

---

## 9. Consistency Over Time

### 9.1 Backwards compatibility

- Tool names MUST NOT change once published.
- Argument names MUST NOT be renamed.
- Fields may only be added, never removed.

### 9.2 Stable affordances

Long conversations assume tool stability.
Breaking changes degrade model reliability.

---

## 10. When Adding a New Tool

Before adding a new tool, confirm:

- [ ] Name follows `noun_verb`
- [ ] Resource prefix matches existing tools
- [ ] Input uses IDs, not slugs
- [ ] Description is minimal and non-redundant
- [ ] Schema is explicit and fail-fast
- [ ] Output shape is stable and predictable
- [ ] Tool intent does not overlap with an existing tool

If any box cannot be checked, STOP and redesign.

---

## 11. Philosophy (Non-Negotiable)

- Tools are designed for **automation reliability**, not human aesthetics.
- Boring, explicit schemas outperform clever abstractions.
- Ambiguity is the primary enemy of LLM tool usage.
- Human DX is handled elsewhere; tool schemas are contracts.

---

## 12. File Locations and Code Generation

### 12.1 Source of truth

The tool schema is defined in:

```
api/robots.yaml           # Main schema file with references
api/robots/tools/*.yaml   # Individual tool definitions
```

`robots.yaml` is the main entry point, but tool definitions can be split into separate files in `api/robots/tools/` for better maintainability.

### 12.2 Split tool files

**ONE TOOL = ONE FILE**

Each tool and its related schemas (Input, Output, supporting types) should live in a single YAML file.

**Example structure:**

```
api/robots/tools/search.yaml    # Contains ToolSearch, ToolSearchInput, ToolSearchOutput, SearchedItem
api/robots/tools/thread.yaml    # Contains all thread-related tools
```

The main `robots.yaml` file references these with:

```yaml
definitions:
  ToolSearch:
    $ref: "./robots/tools/search.yaml"
```

The referenced file contains all related schemas:

```yaml
x-storyden-role: USE_ROBOTS
title: search
description: "Search the Storyden knowledge base..."
type: object
required:
  - input
  - output
properties:
  input:
    $ref: "#/definitions/ToolSearchInput"
  output:
    $ref: "#/definitions/ToolSearchOutput"
definitions:
  ToolSearchInput:
    # ... input schema ...
  ToolSearchOutput:
    # ... output schema ...
  SearchedItem:
    # ... supporting type ...
```

### 12.3 Permission annotations

Tools can specify required permissions using the `x-storyden-role` extension:

```yaml
ToolSearch:
  x-storyden-role: USE_ROBOTS
  title: search
  # ... rest of tool definition
```

This permission is:
- Extracted at runtime by the Go bindings
- Enforced automatically before tool execution
- Type-checked against `rbac.Permission` enum

Valid permission values must exist in `app/resources/rbac/permission.go`.

### 12.4 Code generation

After modifying tool schemas, regenerate the code:

```bash
go generate ./api/...
```

This:
1. Dereferences all `$ref` paths and inlines them into `mcp/robots.json`
2. Generates `mcp/mcp_schema.go` with type-safe Go structs
3. Extracts `x-storyden-role` annotations into `ToolDefinition.RequiredPermission`

### 12.5 Related files

- `api/robots/tools/CLAUDE.md` - These guidelines (you are here)
- `api/robots.yaml` - Main tool schema with references
- `api/robots/tools/*.yaml` - Individual tool definitions
- `mcp/robots.json` - Dereferenced schema (generated, do not edit)
- `mcp/mcp_schema.go` - Generated Go code (do not edit)
- `mcp/bindings.go` - Go bindings with permission extraction
- `internal/tools/schemaderef/` - Schema dereferencing tool

---

**This document is authoritative.**
All agents modifying or auditing tools MUST follow these rules exactly.
