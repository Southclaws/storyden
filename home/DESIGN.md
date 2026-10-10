---
name: Storyden Marketing
description: Woodland storytelling for the Storyden website, blog, and documentation shell.
colors:
  forest: "#285f39"
  forest-brand: "#307343"
  moon: "#104059"
  fire: "#d68e4d"
  newspaper: "#d8dbcd"
  paper: "#f2f0ef"
  ink: "#131b1a"
typography:
  display:
    fontFamily: "var(--font-joie), sans-serif"
    fontSize: "clamp(2.3rem, 5vw, 4.6rem)"
    fontWeight: 700
    lineHeight: 1.08
    letterSpacing: "-0.035em"
  headline:
    fontFamily: "var(--font-joie), sans-serif"
    fontSize: "clamp(2.25rem, 4.5vw, 4.2rem)"
    fontWeight: 700
    lineHeight: 1.08
    letterSpacing: "-0.035em"
  title:
    fontFamily: "var(--font-joie), sans-serif"
    fontSize: "clamp(1.5rem, 2.4vw, 2.3rem)"
    fontWeight: 700
    lineHeight: 1.15
    letterSpacing: "-0.025em"
  body:
    fontFamily: "var(--font-worksans), sans-serif"
    fontSize: "1.06rem"
    fontWeight: 400
    lineHeight: 1.7
  caption:
    fontFamily: "var(--font-worksans), sans-serif"
    fontSize: "0.78rem"
    fontWeight: 400
    lineHeight: 1.55
rounded:
  label: "4px"
  action: "8px"
  frame: "12px"
spacing:
  compact: "12px"
  related: "24px"
  group: "48px"
  chapter: "clamp(80px, 9vw, 144px)"
components:
  button-paper:
    backgroundColor: "{colors.newspaper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.action}"
    padding: "14px 20px"
  button-forest:
    backgroundColor: "{colors.forest}"
    textColor: "#ffffff"
    rounded: "{rounded.action}"
    padding: "14px 20px"
  optional-label:
    rounded: "{rounded.label}"
    padding: "4px 8px"
---

# Design System: Storyden Marketing

## Overview

**Creative North Star: "A lasting home"**

This record applies only to `home/`: Storyden's marketing website, blog, and documentation shell. It captures an extension of the existing woodland identity, not a rebrand. [Root PRODUCT.md](../PRODUCT.md) remains the product truth; [root DESIGN.md](../DESIGN.md) governs the separate `web/` product and the hierarchy represented inside product demonstrations. Neither document is replaced by this marketing record.

Warm woodland imagery, expressive Joie headings, and readable Work Sans prose invite people to imagine their community's home. Broad colour chapters support storytelling; compact demonstrations make the discussion and Library model concrete. Existing brand marks, original artwork, and the established hero remain authoritative. The approved homepage composition and retained-content inventory belong in the [surface brief](.impeccable/surfaces/home-src-app-home-page-tsx.md), not in reusable global rules.

**Key Characteristics:**

- Woodland imagery and spacious, full-width colour chapters.
- Expressive marketing typography surrounding quiet product examples.
- Direct invitations, visible optionality, and readable community stories.

Evidence: `src/app/marketing.css`, `panda.config.ts`, `src/app/globals.css`, `src/fonts/index.ts`, and the `CommunityStory` and `ProductScenes` components under `src/app/(home)/sections/`.

## Colors

The palette combines woodland greens, moonlit blue, warm paper, and dark ink, with campfire colour used sparingly for emphasis.

### Primary

- **Forest:** the deeper marketing chapter and shared documentation accent.
- **Forest brand:** the incumbent Panda `Primary.forest` colour; retain it in existing illustrations and treatments. It is a distinct source token, not an alias for the deeper marketing forest.

### Secondary

- **Moon:** the blue backdrop for care and assistance stories.
- **Fire:** warm emphasis in ownership diagrams and sequence highlights.

### Neutral

- **Newspaper:** softly green paper for supporting bands, footer, and light actions.
- **Paper:** the warm reading ground and light text on dark chapters.
- **Ink:** principal dark text and the ownership chapter ground.

**The Separate Worlds Rule.** Marketing backgrounds may be expressive; product examples retain neutral, bordered surfaces and compact product hierarchy.

## Typography

**Display Font:** local Joie Grotesk Bold, with sans-serif fallback.
**Body Font:** local Work Sans, with sans-serif fallback.

Joie supplies the website's recognisable voice. Work Sans carries prose, links, navigation, and descriptive labels. The recorded display, headline, and title roles form a fluid hierarchy; body copy has generous leading and a typical limit of 48ch. At small widths, chapter prose becomes 1rem. Captions identify illustrations and screenshots without competing with their content.

Local Hedvig Letters Serif, Intel One Mono, and Gorton retain their established editorial, code, and specification uses. Their presence does not make them interchangeable heading styles. Inter is confined to product UI demonstrations: neutral compact headings, prose, metadata, references, and status labels.

**The Marketing Voice Rule.** Use Joie for marketing headings and Work Sans for their supporting text; use the product typeface only where the page depicts product UI.

## Layout

The main marketing and blog index container is capped at 1184px, with 48px side gutters on wide screens. Gutters become 24px at 900px and 20px at 700px. Chapters use the fluid spacing token and varied two-column arrangements; at 700px, primary chapter layouts and blog lists become single columns. Keep text and its demonstration together in reading order.

Blog articles use a narrower 740px reading frame with 20px minimum side gutters. The documentation shell retains Fumadocs navigation, search, sidebar, table of contents, and API mechanics. Colour and typography additions must preserve those reading and navigation behaviours.

**The Open Chapter Rule.** Use spacing and full-width colour changes for narrative sections. Reserve frames for actual screenshots, diagrams, or product objects.

## Elevation & Depth

The chapter extension is predominantly flat: colour changes, thin borders, inset references, and artwork establish depth. Product examples use visible boundaries; secondary information sits within quieter inset regions. Existing Panda link buttons, cards, and frosted overlays retain their established shadow or blur behaviour where already used; those are not a requirement for new narrative chapters.

The Library and Trail illustrations are static. The ownership diagram uses a slow, continuous dashed flow from the Storyden mark along softly curved branches, without a junction dot or replay controls. Dashes move 24px every 3.6 seconds with linear timing; labels remain still. Playback pauses offscreen or in a hidden tab. Reduced motion and no-JavaScript rendering show complete, static connections. Action colour transitions are short (150ms). Do not make animation a prerequisite for understanding the page.

## Shapes

Marketing actions have gently rounded corners, compact labels have tighter corners, and screenshots and demonstration frames use the larger frame radius. Fine rules separate related passages and list rows. Circular woodland crops and existing landscape artwork remain part of the identity; these shapes do not justify wrapping every paragraph in a card.

## Components

### Actions

Paper actions pair newspaper with ink on dark chapters; forest actions pair forest with white on light grounds. Both use a 48px minimum height, semibold Work Sans, and the action radius. Paper actions turn white on hover; forest actions darken. Text actions retain normal reading colour and gain an underline on hover. Marketing links and buttons expose a two-pixel current-colour focus outline with a five-pixel offset.

### Labels and product examples

The outlined Optional AI label names a real capability condition. Product status labels express actual state, such as scheduled, completed, or needing attention. Illustrative scenes are static examples, labelled as illustrative; actual screenshots are identified separately. Do not make imitation controls appear interactive.

Product scenes use neutral bordered frames and compact Inter. Embedded references and results use inset tone. The ownership connection diagram retains marketing typography: its community source is visually stronger than its branch headings, while supporting protocol names remain subordinate.

### Navigation and reading surfaces

Shared navigation uses a paper ground with a forest live-community action. Active documentation navigation uses forest text and a quiet green tint. The footer groups real links on newspaper. The blog index uses dated editorial rows and landscape artwork, while article pages favour an open reading frame over decorative cards.

## Do's and Don'ts

### Do:

- **Do** preserve the existing hero, brand marks, and woodland assets when extending this identity.
- **Do** distinguish marketing typography from product demonstrations.
- **Do** label illustrative communities and optional assistance clearly.
- **Do** preserve visible focus, readable contrast, mobile reading order, and complete reduced-motion content.

### Don't:

- **Don't** apply this marketing palette or typography as authority over the separate product application.
- **Don't** turn the homepage's chapter order into a universal rule for other surfaces.
- **Don't** imply fabricated customers, metrics, or guaranteed AI results through examples.
- **Don't** canonize isolated layout exceptions or decorative labels as reusable design tokens.
