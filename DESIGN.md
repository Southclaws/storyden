---
name: Storyden
description: Refined, restrained community software shaped around the Quiet Campfire.
colors:
  bg-site: "#fcfcfc"
  bg-default: "#f9f9f9"
  bg-subtle: "#f0f0f0"
  bg-muted: "#e8e8e8"
  fg-default: "#202020"
  fg-subtle: "#646464"
  fg-muted: "#838383"
  border-default: "#e0e0e0"
  border-muted: "#cecece"
  accent-default: "var(--accent-colour-flat-fill-400)"
  accent-subtle: "var(--accent-colour-flat-fill-200)"
typography:
  display:
    fontFamily: "var(--font-inter-display), Inter, sans-serif"
    fontSize: "2.488rem"
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: "-0.02em"
  headline:
    fontFamily: "var(--font-inter-display), Inter, sans-serif"
    fontSize: "2.074rem"
    fontWeight: 600
    lineHeight: 1.25
  title:
    fontFamily: "var(--font-inter-display), Inter, sans-serif"
    fontSize: "1.44rem"
    fontWeight: 600
    lineHeight: 1.25
  body:
    fontFamily: "var(--font-inter), Inter, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "var(--font-inter), Inter, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 600
    lineHeight: 1.25
rounded:
  l1: "0.125rem"
  l2: "0.25rem"
  l3: "0.375rem"
  card: "0.5rem"
  panel: "0.75rem"
  full: "9999px"
spacing:
  "1": "0.25rem"
  "2": "0.5rem"
  "3": "0.75rem"
  "4": "1rem"
  "6": "1.5rem"
  "8": "2rem"
components:
  button-solid:
    backgroundColor: "{colors.accent-default}"
    textColor: "{colors.bg-default}"
    typography: "{typography.label}"
    rounded: "{rounded.l2}"
    padding: "0 1rem"
    height: "2.5rem"
  button-outline:
    textColor: "{colors.fg-subtle}"
    typography: "{typography.label}"
    rounded: "{rounded.l2}"
    padding: "0 1rem"
    height: "2.5rem"
  input-outline:
    backgroundColor: "transparent"
    textColor: "{colors.fg-default}"
    typography: "{typography.body}"
    rounded: "{rounded.l2}"
    padding: "0 0.75rem"
    height: "2.5rem"
  badge-subtle:
    backgroundColor: "{colors.bg-subtle}"
    textColor: "{colors.fg-default}"
    typography: "{typography.label}"
    rounded: "{rounded.full}"
    padding: "0 0.625rem"
    height: "1.5rem"
  rich-card:
    backgroundColor: "{colors.bg-default}"
    textColor: "{colors.fg-default}"
    rounded: "{rounded.card}"
    padding: "0.5rem"
---

# Design System: Storyden

## Overview

**Creative North Star: "The Quiet Campfire"**

Storyden is a calm gathering place built around participation, conversation, and communal memory. Its interface should feel quietly welcoming and carefully made: present enough to guide people, restrained enough to let their community, content, and chosen accent color remain the protagonists.

The system is refined and modern without becoming corporate or precious. Compact type, modest controls, and close spatial relationships provide useful information density; neutral surfaces, clear hierarchy, and selective depth keep that density readable. The result should feel premium through consistency and craft, not decoration.

Reference influences remain Linear's precision and restraint, Luma's selective moments of fresh color, and Discord's community and moderation ergonomics. They are signals to interpret through Storyden's white-label system, not visual templates to reproduce.

**Key Characteristics:**

- Calm, premium, and community-native rather than corporate SaaS.
- Refined, restrained, and compact, with deliberate information density.
- White-label by design: neutral foundations carry administrator-controlled accents and imagery.
- Modern and minimal without losing the warmth or utility of forum culture.
- Accessible, responsive, and respectful of reduced-motion preferences.

## Colors

Storyden uses semantic color roles as its public vocabulary. Neutral layers establish hierarchy while the administrator-controlled accent ramp gives each community its own identity.

### Primary

- **`accent.default`**: the dynamic community accent used for primary actions, selected states, and focused emphasis.
- **`accent.subtle`**: a low-intensity accent layer for quiet selection, hover, and contextual emphasis.

### Neutral

- **`bg.site`**: the outermost application canvas.
- **`bg.default`**: the default control, card, and contained-surface layer.
- **`bg.subtle`**: a gentle step above the default surface for badges, groups, and separation.
- **`bg.muted`**: the strongest routine neutral layer for hover and structural emphasis.
- **`fg.default`**: primary text and icons.
- **`fg.subtle`**: supporting copy and secondary controls.
- **`fg.muted`**: metadata, timestamps, and low-priority information.
- **`border.default`**: ordinary control and surface boundaries.
- **`border.muted`**: boundaries that need slightly greater definition.

### Named Rules

**The Community Owns the Accent Rule.** Use semantic `accent.*` roles rather than a fixed Storyden hue. The neutral system must let the host community's configured color remain the primary brand expression.

**The Quiet Canvas Rule.** Most screen area belongs to `bg.site`, `bg.default`, and neutral foreground roles. Accent is purposeful punctuation, not a wash over the interface.

## Typography

**Display Font:** Inter Display (with Inter and sans-serif fallbacks)  
**Body Font:** Inter (with sans-serif fallback)  
**Label/Mono Font:** Inter for interface labels; the system monospace stack for code and technical values

**Character:** Inter keeps dense community interfaces clear and familiar, while Inter Display gives larger headings slightly more authority without introducing a separate stylistic voice. Compact cards intentionally keep titles close to body scale so lists remain scannable.

### Hierarchy

- **Display** (600, `2.488rem`, 1.25): rare page-level statements and high-emphasis empty or onboarding states.
- **Headline** (600, `2.074rem`, 1.25): major page or section introductions.
- **Title** (600, `1.44rem`, 1.25): prominent section and panel titles.
- **Body** (400, `1rem`, 1.5): discussion, descriptions, settings copy, and general reading; long prose should stay near the `65ch` prose measure.
- **Label** (600, `0.875rem`, 1.25): buttons, navigation, compact controls, and metadata labels. Extra-compact contexts may step down to `0.75rem`.

### Named Rules

**The Density Without Crowding Rule.** Prefer compact type and short vertical rhythms, but preserve line height, grouping, and contrast so the interface scans as structured information rather than compressed noise.

## Layout

The application shell centers a main content column beside a persistent desktop navigation sidebar. Desktop navigation is approximately `16rem` wide, widening to `18rem` on large screens, with a `1.5rem` to `2rem` gap before the content column. Viewport padding grows from `0.5rem` on mobile to `1rem` on desktop and `1.5rem` on wide screens.

The system uses breakpoints at `640px`, `768px`, `1024px`, `1280px`, and `1536px`. Below `768px`, the sidebar becomes a modal drawer and a compact command bar remains fixed near the bottom safe area. Responsive rich cards use container queries: a row with optional media on wide containers becomes a stacked composition on narrow containers. Card grids move from one to two and then three columns as their container permits.

Spacing follows a quarter-rem base scale. Routine component gaps cluster around `0.25rem` to `0.75rem`; component padding commonly uses `0.5rem` to `1rem`; larger shell and section separation uses `1.5rem` to `2rem`. Information density should come from consistent alignment, shared baselines, and compact grouping—not from removing meaningful separation.

## Elevation & Depth

Storyden is quietly layered. Tonal separation and hairline borders do most of the structural work. Small ambient shadows lift cards just enough to distinguish them from the site canvas, while larger shadows are reserved for floating menus, popovers, dialogs, and tooltips. Navigation chrome may use translucent `bg.site` surfaces with subtle blur so it stays available without becoming visually heavy.

### Shadow Vocabulary

- **Subtle surface** (`0 1px 2px rgb(0 0 0 / 2.4%), 0 0 1px rgb(0 0 0 / 6%)`): resting rich cards and quiet contained surfaces.
- **Small float** (`0 2px 4px rgb(0 0 0 / 6%), 0 0 1px rgb(0 0 0 / 19%)`): tooltips and small transient surfaces.
- **Large float** (`0 8px 16px rgb(0 0 0 / 6%), 0 0 1px rgb(0 0 0 / 19%)`): menus, popovers, and modal layers.

### Named Rules

**The Quietly Layered Rule.** Establish hierarchy with tone and borders first. Add shadow only when a surface truly sits above another surface or must remain legible while floating over content.

## Shapes

Controls use gently rounded corners, typically the semantic `l2` radius (`0.25rem`), while compact nested elements may use `l1` (`0.125rem`). Content cards use a soft `0.5rem` radius, floating mobile navigation uses `0.75rem`, and pills or badges use a full radius. Borders are normally one pixel and neutral; shape changes should communicate containment or state rather than decoration.

The silhouette stays simple and rectangular. Rounded forms soften dense information but should not turn every section into an isolated card. Media is clipped to its containing card and may use a smaller nested radius when stacked inside it.

## Components

Components are refined and restrained, but compact enough to support real community and administrative information density. Their states should be obvious through tone, border, and focus treatment rather than exaggerated scale or ornament.

### Buttons

- **Shape:** gently rounded (`l2`, `0.25rem`) with compact heights from `1.5rem` to `4rem`; `2rem` and `2.5rem` are the common interface sizes.
- **Primary:** `accent.default` background with `bg.default` text, semibold label type, and horizontal padding proportional to size.
- **Hover / Focus:** the primary fill softens on hover; keyboard focus uses a two-pixel accent outline with a two-pixel offset.
- **Outline / Ghost / Subtle:** outline uses a neutral or palette-aware hairline; ghost relies on text until hover; subtle uses a muted translucent surface. All disabled states lower contrast and remove pointer affordance.

### Chips

- **Style:** badges use a full-radius silhouette, compact label text, and either `bg.subtle` with `border.subtle`, a solid palette fill, or a stronger outline.
- **State:** use semantic palettes for category and status meaning; keep the default badge quiet enough to coexist with dense metadata.

### Cards / Containers

- **Corner Style:** soft (`0.5rem`) on rich cards; smaller nested radii for media and child rows.
- **Background:** `bg.default` at rest, with semantic emphasized and accent-backed variants for meaningful states.
- **Shadow Strategy:** the subtle surface shadow is the default; avoid stacking shadows on nested children.
- **Border:** use borders only where a boundary or state needs more definition than tone and shadow provide.
- **Internal Padding:** rich cards use a compact `0.5rem` edge rhythm and slot-based layout rather than broad generic padding.

### Inputs / Fields

- **Style:** transparent background, one-pixel `border.default`, `l2` radius, and compact horizontal padding.
- **Focus:** preserve visible keyboard focus and use semantic accent or error treatment rather than removing the boundary without replacement.
- **Error / Disabled:** error uses the semantic error border; disabled fields lower opacity and remove pointer affordance without hiding their value.

### Navigation

Desktop navigation is a persistent left-hand information tree with a quiet frosted top bar. Navigation anchors use compact ghost-button styling, muted default color, and clear hover or current states. Mobile uses a bottom command bar plus a modal navigation drawer with focus trapping, inert background content, safe-area spacing, a scrim, and reduced-motion support.

### Rich Card

The rich card is Storyden's signature reusable content surface. Its named slots coordinate title, body, metadata, controls, and optional media. Row, responsive, box, and fill shapes adapt the same content hierarchy to feeds, lists, and grids; container queries preserve the composition when the card is embedded in narrower contexts.

## Do's and Don'ts

### Do:

- **Do** use the existing semantic roles—`bg.*`, `fg.*`, `border.*`, and `accent.*`—so white-label themes remain coherent.
- **Do** make dense screens legible through alignment, hierarchy, and consistent compact spacing.
- **Do** keep community content and host branding visually stronger than the application chrome.
- **Do** preserve keyboard focus, sufficient contrast, reduced-motion behavior, and responsive access to every workflow.
- **Do** favor restrained motion, precise typography, and carefully resolved interaction states as the source of polish.

### Don't:

- **Don't** bypass semantic tokens with one-off raw colors or fixed brand hues.
- **Don't** reproduce legacy PHP forum styling, generic corporate SaaS dashboards, or dense AWS Console-like control surfaces.
- **Don't** turn every group into a floating card or use shadow where a tonal layer or border is sufficient.
- **Don't** sacrifice useful information density by inflating controls, type, or whitespace indiscriminately.
- **Don't** let decorative effects overpower discussion, knowledge, or the identity of the host community.
