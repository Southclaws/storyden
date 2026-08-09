---
name: Storyden
description: Compact, calm community software shaped around the Quiet Campfire.
colors:
  background-canvas: "#fcfcfc"
  background-surface: "#f9f9f9"
  background-inset: "#f0f0f0"
  background-overlay: "#ffffff"
  background-control: "#fcfcfc"
  text-default: "#202020"
  text-muted: "#646464"
  border-default: "#e0e0e0"
  accent-default: "var(--accent-colour-flat-fill-400)"
typography:
  page-heading:
    fontFamily: "var(--font-inter-display), Inter, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 600
    lineHeight: 1.5
  section-heading:
    fontFamily: "var(--font-inter-display), Inter, sans-serif"
    fontSize: "1rem"
    fontWeight: 700
    lineHeight: 1.5
  body:
    fontFamily: "var(--font-inter), Inter, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  supporting:
    fontFamily: "var(--font-inter), Inter, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.25
  metadata:
    fontFamily: "var(--font-inter), Inter, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 500
    lineHeight: 1.25
rounded:
  control: "0.25rem"
  panel: "0.5rem"
  overlay: "0.75rem"
  pill: "9999px"
spacing:
  "1": "0.25rem"
  "2": "0.5rem"
  "3": "0.75rem"
  "4": "1rem"
  "6": "1.5rem"
  "8": "2rem"
components:
  button-subtle:
    backgroundColor: "{colors.background-inset}"
    textColor: "{colors.text-default}"
    rounded: "{rounded.control}"
    height: "1.5rem"
  button-solid:
    backgroundColor: "{colors.accent-default}"
    rounded: "{rounded.control}"
    height: "1.5rem"
  surface:
    backgroundColor: "{colors.background-surface}"
    textColor: "{colors.text-default}"
    rounded: "{rounded.panel}"
  input:
    backgroundColor: "{colors.background-control}"
    textColor: "{colors.text-default}"
    rounded: "{rounded.control}"
    height: "1.5rem"
---

# Design System: Storyden

## Overview

**Creative north star: The Quiet Campfire.**

Storyden is calm community software. The interface should guide participation,
curation, moderation, and administration without competing with the community's
identity or content. It should feel compact and precise, but never cramped;
neutral and themeable, but never anonymous; modern, but still recognisably a
place where people talk and build shared knowledge.

Storybook is the directional authority for implemented UI. It is not only a
catalogue of component capabilities. Before composing a screen, read
`Foundations/Principles/Storyden`, then the relevant component and composition
stories. This document defines durable philosophy; Storybook shows that
philosophy rendered in the current system.

When making a UI decision, use this order:

1. Identify the user's task, context, and return path.
2. Choose the screen width and navigation model.
3. Establish page, section, supporting-text, and metadata hierarchy.
4. Assign structural layers: Canvas, Surface, Inset, or Overlay.
5. Choose action emphasis from workflow priority and consequence.
6. Verify loading, empty, disabled, error, mobile, keyboard, and dark states.

The default answer is usually the quietest existing component that communicates
the role clearly. New visual components are design-system decisions, not local
conveniences.

## Colors

### Structural layers

Background tokens describe **structural role**, not appearance or arbitrary
elevation:

- **`background.canvas`** is the root application background. Page headings,
  settings forms, editors, and top-level workflows may sit directly on it.
- **`background.surface`** is an independent content object: a post, card,
  report, panel, or substantial list item. Every Surface has a
  `border.default` boundary; fill alone is not sufficient separation from
  Canvas.
- **`background.inset`** is secondary or embedded content: quoted posts,
  previews, metadata blocks, subcategory rows, code, summaries, and dense
  supporting regions. It contrasts with its immediate surroundings.
- **`background.overlay`** is temporary content above document flow: menus,
  popovers, dialogs, tooltips, command palettes, and date pickers. It uses a
  95% opaque fill with the shared 4px `blurs.subtle` backdrop treatment.
  Full overlays use `shadows.overlay`; compact transient content such as
  tooltips and editor bubbles uses `shadows.floating`.

Storyden has two ordinary content layers: Surface and Inset. Avoid visually
identical Surface-in-Surface nesting; subordinate content belongs on Inset.
Avoid Inset-in-Inset. An Inset may sit directly on Canvas when its meaning is
clearly supporting rather than independent.

This is a hard rule: all Cards are Surfaces and all Surfaces have borders.
Component recipes enforce the boundary so call sites do not decide whether a
Card needs one. A media-backed Card may omit the redundant border between its
cover and content, but the content region still retains the Surface boundary.

Controls use the **`background.control*`** family. A Quick Share composer,
input, select, or editable field is a control, not a content Surface. Use
`background.controlInset` only when a control itself belongs in an Inset region.

### Foregrounds and borders

- **`text.default`** is primary UI copy and content labels.
- **`text.muted`** is supporting text, metadata, secondary navigation, and
  subdued icons.
- **`text.subtle`** is lower-emphasis chrome only. Do not use it for compact
  prose that must meet normal-text contrast.
- **`text.disabled`** communicates unavailable controls and must remain
  readable in both themes.
- **`border.default`** separates ordinary regions and Surface edges.
- **`border.muted`** is the quietest structural separator.
- **`border.strong`** is reserved for focus, selection, or boundaries that must
  remain legible against neighbouring layers.

### Meaning and interaction

Structural hierarchy and semantic meaning are separate:

- Accent communicates the host community, focus, deliberate selection, and a
  primary action. It must not become a decorative wash over the interface.
- Success, warning, danger, and information tokens communicate status. Alerts
  use Inset structure, then layer status tokens over it.
- A Button `variant` communicates priority; `intent` communicates meaning. Use
  `intent="success"`, `intent="warning"`, or `intent="destructive"` rather
  than naming a colour palette at a call site. A destructive action may be a
  quiet outline before confirmation, then solid for the final irreversible
  command. `colorPalette` is a low-level theming escape hatch, not the product
  API for standard actions.
- Badges use colour only when the colour has stable meaning: category, state,
  role, or classification. Generated tag colours are valid when paired with
  readable foregrounds.

Never select raw neutral ramp values at a call site merely because a grey looks
right. If the existing semantic role is insufficient, define the missing role.

## Typography

Storyden has separate systems for product UI and community-authored content.

### Product UI

- **`PageHeading`** is the screen identity and always renders an `h1`. Use one
  per screen when the identity is not already unambiguous from the content.
- **`SectionHeading`** introduces a meaningful subdivision and always renders
  an `h2`. It is intentionally smaller, heavier, and muted.
- **`Text` `body`** is primary UI prose and explanations.
- **`Text` `supporting`** sits below headings, labels, or primary content to
  explain purpose or consequence.
- **`Text` `metadata`** is for timestamps, counts, identifiers, compact status,
  and other low-priority annotations.

Do not rebuild these roles with `styled.p`, `fontSize`, and `color` at call
sites. Semantic HTML may vary through `Text as`, but the visual role stays
constrained.

Surface titles belong to the Surface recipe because their visual size depends
on that object, not on document heading level. A Rich Card may contain an `h1`
or `h2` based on its place in the semantic tree without looking like a page
heading.

### Community content

Rich posts and Library prose use the global `.typography` system. Their type
scale may be more expressive because it serves reading and authored hierarchy.
Do not use product UI heading components inside rendered rich text, and do not
use rich-text styles to lay out settings or tools.

### Density

The canonical compact scale is intentional: controls default to `sm`, ordinary
UI copy is 16px, supporting copy is 14px, metadata and compact actions are 12px.
Increase size only when hierarchy or touch ergonomics genuinely requires it.
Information density comes from consistent alignment and close semantic groups,
not by collapsing line height or removing necessary separation.

## Layout

### Page widths

Use `PageLayout` to declare the screen's workspace:

- **`content`** is the default for feeds, settings, forms, and reading-focused
  pages.
- **`wide`** is for dashboards, tables, Robots, multi-column tools, and screens
  that benefit from comparison or parallel work.
- **`full`** is for canvas-like editors and data tools. Full width is not
  permission for full-width prose; constrain readable text internally.

Width is a screen-level decision. Do not make a child card compensate for an
incorrect page width.

### Page headers

`PageHeader` is the canonical top-of-page composition for return or hierarchy
navigation, `PageHeading`, supporting text, and a small action group.

- Breadcrumbs communicate **where the page is** in an information hierarchy.
- `BackAction` communicates **where the user can return**. It is especially
  important on mobile screens where contextual desktop navigation is absent.
- Omit a page heading only when the primary content makes the screen identity
  genuinely obvious. Do not omit it merely to save vertical space.
- Page-level Save/Create actions align with the heading when they affect the
  whole page. Section actions align with their `SectionHeading`.

Settings screens sit flush on Canvas. Use a stack with gap 4 between sections;
within each heading group, use gap 1 between heading and supporting text.
Do not wrap the entire settings screen or each arbitrary section in CardBox.

### Responsive navigation

Desktop uses the persistent left sidebar as the canonical navigation surface;
there is no desktop top bar. The sidebar may swap to contextual navigation via
parallel routes for Settings, Admin, or Robots.

Mobile keeps the global site tree in a fixed top navigation and top-entering
drawer. Contextual navigation remains visible in the page: `SectionNavigation`
for route-heavy Settings/Admin, Tabs for a small set of peer views within one
product feature such as Robots.

Never render the heavy category/Library navigation tree twice to solve a
responsive layout. The same DOM navigation pane is repositioned with CSS.

### Composition

Use spacing to communicate relationship:

- gap 1: title/supporting copy, badges in a row, compact metadata.
- gap 2: controls within one command group or card internals.
- gap 3: related fields and repeated compact objects.
- gap 4: page sections and major form groups.
- gap 6 or 8: only for genuinely separate screen regions.

Prefer unframed full-width sections over decorative page cards. Cards are for
objects, not for dividing every part of a page.

## Elevation & Depth

Tone and border establish ordinary hierarchy. Shadows are reserved for spatial
elevation:

- **`shadow.surface`** is optional and quiet. Most Surfaces rely on fill and
  border instead.
- **`shadow.floating`** separates tooltips and compact transient surfaces.
- **`shadow.overlay`** separates menus, popovers, and dialogs where overlapping
  content requires a stronger boundary.

An Overlay must use the canonical closed component defaults for portalling,
placement, collision handling, viewport fitting, and click-away behavior.
Do not add Floating UI parameters independently at every call site.

Motion should explain a state change, not decorate it. Keep transitions short
and restrained. Preserve reduced-motion behavior. Drag overlays disappear on
drop without decorative return animations; the final object should remain
stable in its new position.

## Shapes

Storyden uses modest radii:

- **`radii.control`** for buttons, inputs, compact rows, and controls.
- **`radii.panel`** for Surfaces, cards, and structured content objects.
- **`radii.overlay`** for menus, dialogs, and temporary floating content.
- **`radii.pill`** for badges and truly pill-shaped state labels.

Do not use large rounded containers to make ordinary page sections feel like
marketing cards. Media may run edge-to-edge at the top of a Surface; clip it to
the outer panel radius and keep the content boundary coherent.

Borders should normally be one pixel. Remove a redundant edge only when an
adjacent element already provides the same boundary, such as a cover image
meeting a card content area.

## Components

### Buttons and icon buttons

Buttons default to `size="sm"` and `variant="subtle"`.

- **`subtle`**: routine product commands such as Save, Create, Connect, Edit,
  Share, and low-risk workflow actions. This is the default.
- **`solid`**: the single decisive completion action in a workflow, or the
  confirmed destructive action with `intent="destructive"`. Avoid multiple
  competing solid buttons in one action group.
- **`outline`**: a visible secondary peer, alternate route, upload/browse
  action, or control that needs a persistent boundary against its surroundings.
- **`ghost`**: actions that should recede into page or card chrome, including
  Cancel, Close, pagination, and contextual object tools.
- **`plain`**: only for dense navigation/tree rows where a visible container
  would determine row height. It keeps the accessible hit area while changing
  foreground colour on hover.

Choose meaning independently through `intent`:

- **`success`**: confirms or completes a positive outcome. Do not use it as a
  substitute for the ordinary accent action.
- **`warning`**: proceeds with caution or enters a consequential moderation
  flow that is not itself destructive.
- **`destructive`**: removes, revokes, purges, suspends, or permanently changes
  data. Use a quiet variant to enter the flow and `solid` only for its final
  confirmation.

Do not pass raw red, green, orange, or amber `colorPalette` values to Buttons.
If these intents are insufficient, define the missing semantic intent before
introducing another product action colour.

Icon-only buttons require one accessible label prop; the component supplies the
tooltip/title behavior. Use familiar icons and keep their rendered size aligned
with the control size. More-menu buttons should be quiet and only appear when
the commands cannot be made clearer or more direct.

### Surfaces and cards

Use the Surface family according to content role:

- **`Card`/Rich Card** is the canonical independent object when content has a
  title, description, metadata, media, actions, or supporting rows. Its Surface
  border is mandatory.
- **`CardGrid`** and **`CardRows`** arrange repeated Surface objects and own
  equal-height/grid heuristics. Do not recreate card grids locally.
- **`CardBox`** is a compatibility Surface for compact arbitrary independent
  content. Keep its padding and radius defaults, but do not use it as a generic
  wrapper around a page, settings screen, or layout section.
- Subordinate content inside a Surface uses Inset styling, not another
  visually identical Surface.

Category cards, thread cards, Library cards, reports, and similar content may
share Surface fundamentals but retain domain-specific slots and interaction.
Preserve media, metadata, controls, and semantic article structure rather than
flattening every object into one generic card API.

### Forms and controls

Inputs, Selects, Comboboxes, Number Inputs, Pin Inputs, Checkboxes, Radios,
Switches, and Sliders share `sm`, `md`, and `lg`; `sm` is the default. Controls
that appear in one row must align in height.

Use `FormControl` for label, control, helper text, and error text. Helper and
error copy use the metadata text role. React Hook Form adapters live beside the
component as `.field.tsx`; do not create feature-local typed wrappers.

Use an outline control by default. Use an Inset control when embedded in an
Inset region, such as a Library directory block. Use ghost inputs only for
deliberate inline editing or composer patterns. Prefixes, suffixes, and icon
triggers must inherit the parent control's size and padding.

Disabled controls must remain visibly and textually distinguishable in both
themes. Enter inside a staged-list editor should perform that local add action,
not accidentally submit the enclosing settings form.

### Navigation

- **Tabs** switch a small, stable set of peer views inside one feature.
- **SectionNavigation** navigates route-heavy settings/admin sections on small
  screens and mirrors the desktop contextual sidebar.
- **Breadcrumbs** describe hierarchy; they are not a browser-history control.
- **DragTree** is the interactive sortable hierarchy for category and Library
  navigation. Matching/selection logic is supplied by the domain; the component
  remains route-agnostic.
- **TreeView** is the read-only hierarchy and should share DragTree's visual
  language without drag behavior.

Root tree items are flush with the sidebar; depth begins visually at one.
Selected items use neutral emphasis, not a gratuitous accent fill. Tree rows and
ordinary sidebar links share the same muted resting foreground.

### Overlays and feedback

Prefer closed Menu, Popover, Tooltip, and Date Picker components with canonical
positioning defaults. Use compound primitives only for genuinely custom
composition. Menus contain contextual commands, not large forms or site
navigation.

- Alert: persistent task-relevant information, structurally Inset plus status.
- Admonition: the preferred in-flow alternative to a toast for task-level
  errors, warnings, and meaningful success feedback. It remains near the
  relevant content and may be dismissed once acknowledged.
- Toast: avoid by default. Reserve it for non-critical, transient confirmation
  with no action or recovery. Do not use it for errors, warnings, or information
  a user may need time to find, read, or revisit.
- Helper/error text: field-level instruction or validation.
- Badge: compact metadata or classification, not an action.

This follows [GitHub Primer's accessible notifications
guidance](https://primer.style/accessibility/toasts/): toast interfaces carry
inherent accessibility and usability risks, while in-flow messages are easier
to discover and revisit.

### Block editors and draggable UI

The shared Block Editor UI owns theme-ready gutters, handles, contextual menus,
and interaction chrome. Domain editors own block schemas, permissions,
persistence, and drag payloads understood by the global drag/drop provider.

Edit mode should not shift the read layout. Sidebar edit overlays sit over the
existing rows rather than wrapping them in larger containers. Custom drag
overlays preserve the source object's proportions, follow the pointer, and do
not animate back after drop. Interaction tests must cover dragging, opening the
menu, clicking controls inside it, clicking inert menu chrome, and click-away.

## Do's and Don'ts

### Do

- Use semantic layer, text, border, status, and control tokens.
- Consult Storybook principles and realistic compositions before building.
- Keep community content and host branding stronger than application chrome.
- Prefer compact defaults, stable dimensions, and aligned controls.
- Preserve keyboard focus, accessible names, contrast, reduced motion, and
  responsive access to every workflow.
- Render one expensive navigation tree and reposition it responsively.
- Test interaction-heavy components from the outside, including slow requests,
  hydration, drag/drop, focus, click-away, and viewport edges.

### Don't

- Do not invent one-off backgrounds, wrappers, card treatments, or text roles.
- Do not turn every page section into a floating card or nest identical
  Surfaces.
- Do not use accent colour merely to make an active state more obvious when a
  neutral selected treatment is sufficient.
- Do not inflate controls or headings to create hierarchy that composition
  should provide.
- Do not put settings navigation in horizontally overflowing Tabs.
- Do not duplicate server-rendered navigation content for desktop and mobile.
- Do not treat Storybook as a prop playground; stories must explain the approved
  product context, decision, and boundaries of each component.
