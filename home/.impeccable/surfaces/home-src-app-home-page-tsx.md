---
version: 1
slug: "home-src-app-home-page-tsx"
primary_target: "home/src/app/(home)/page.tsx"
related_targets:
  [
    "home/src/app/(home)/blog",
    "home/src/app/docs",
    "home/src/app/marketing.css",
  ]
---

# Storyden marketing: a lasting home

Mode: Persuade on the landing page; Read on blog and docs.

## Purpose and audience

Storyden is an open-source forum and knowledgebase for communities whose conversations, discoveries, and knowledge deserve a lasting home. The page speaks to people running clubs, creative communities, gaming groups, research groups, and shared-interest spaces. Its narrative centres on belonging, useful shared knowledge, community care, and ownership.

The primary conversion is visiting the community of the people behind Storyden at `https://makeroom.club`. Starting a self-hosted community is the secondary action. The forum and Library work without a model provider; relevant assistance is explicitly labelled **Optional AI**.

Copy is welcoming, concrete, and outcome-led. Illustrations use fictional members and example activity, without implying customer endorsements or measured usage. Public copy describes the product and community outcomes; it does not describe the design process or explain that an illustration follows Storyden’s layout.

## Identity and sources of truth

Repository-root `PRODUCT.md` defines product context. Root `DESIGN.md` and the implemented `web` application define the product UI represented inside demonstrations. `home/DESIGN.md` defines the marketing visual system.

The marketing identity uses Joie headings, Work Sans prose, woodland artwork, newspaper neutrals, forest green, moonlit blue, and campfire accents. Existing Hedvig, Intel One, and Gorton typography supports specialist passages and technical tables. Product illustrations use compact Inter typography, bordered surfaces, actual product iconography, and a clear metadata hierarchy.

Colour gives the narrative structure: light reading surfaces for discussion, forest green for shared knowledge, moonlit blue for care and assistance, ink for ownership and deployment, and a light closing invitation. Generous product stages and varied text/image arrangements provide pacing without repeated feature-card grids.

## Landing page composition

1. **Hero and product screenshot.** Responsive lake artwork, the Storyden logo, and the Get Started and Live Demo controls introduce the page. Beneath the controls: “An open-source forum and knowledgebase for communities that want their conversations, discoveries, and wisdom to last.” The main application screenshot follows directly and reaches into the opening viewport. Shared navigation has no separate live-demo button.
2. **Manifesto.** The frosted passage below the screenshot contains the copy beginning “In a tangled web of fleeting feeds”, the ownership line, and “Collective memory. Without the noise.” Bold emphasis highlights feeds, what matters, on your terms, data, and corner. It has no separate heading or CTA. Its upper edge overlaps the screenshot by 120px, uses 16px backdrop blur and 72% paper opacity, and fades into opaque paper over 160px.
3. **Conversation.** “Home for your people. And for what they know.” introduces a gaming project discussion beside the copy. “Community you can come back to.” is followed by “A place for conversations to unfold, ideas to take root, and people to stay connected at their own pace.” A documentation link leads to discussion. The Knowledge gardens passage follows, with the sapling artwork and bold emphasis on clubs, teams, and people.
4. **Knowledge.** “Make the good finds part of the furniture.” connects a shared book recommendation to a Library resource and reading list through a static three-step illustration. A real Library directory screenshot and the `organise.png` and `curate.png` diagrams provide product evidence. The copy covers pages, properties, links, and collaborative curation. Optional assistance with titles, tags, summaries, and organisation is presented alongside ordinary editing.
5. **Care.** “A shared space, looked after together.” pairs community-care copy with a submitted report awaiting review. “Running a community is a labour of love.” leads into the community’s character and the tools that help preserve it as it grows. A separate Optional AI passage, “Some of the routine can look after itself.”, pairs configured Robot and Trail assistance with scheduled, completed, and needs-attention states. The three explanatory paragraphs share the same body-text hierarchy; human decisions remain explicit.
6. **Ownership.** “Your community. Your terms.” covers self-hosting, identity, rules, and extension. A centred connection diagram sits directly on the ink background, without an outer frame or caption. Links beneath it read “Your look and feel”, “Your integrations”, and “Everyone’s source code”.
7. **Technical tables and deployment.** The Storyden computer illustration, specification tables, Docker command, and open-source message form the technical close. There is no introductory statistics caption. GitHub stars are fetched with a 24-hour revalidation interval; commits, contributors, lines of code, and API endpoint counts currently render as unavailable. Numeric values must come from verified sources, never estimates or fabricated fallback calculations. Directly after `MilspecSection`, the operations plate “Small by default. Room to grow.” continues the ink and newspaper technical treatment with hairline borders, Gorton headings and diagram labels, and Work Sans prose. Two rows show the web app and API in one Docker image with SQLite and uploads on a persistent local volume, then optional PostgreSQL, S3-compatible file storage, and Redis caching. The footer covers public URLs, a persistent volume, and HTTPS, linking to the Docker Compose guide. The plate presents simple operation with optional growth and inherits `home/DESIGN.md` without introducing new design tokens.
8. **Community invitation.** “Come on in. Make yourself at home.” sits beside a circular crop of `square-campfire.webp`. The copy reads “Meet the people behind Storyden in our community. Then imagine what yours could become.” The primary button is “Explore our community”; the secondary link is “Start your own”.

Original artwork files and independently placeable section components remain available in the repository. `Permanence`, `Gardens`, `Organise`, and `Curate` are composed into the landing page; the complete `CollectiveMemorySection` composition and other unused section components remain available separately.

## Product demonstrations

### Discussion

The discussion illustration follows `web/src/screens/thread/ThreadScreen/ThreadScreen.tsx` and its `Byline`, `ThreadStats`, and `Reply`/`CardBox` compositions. It shows the harbour project, replies, reactions, and an inline Library reference. The member/time/category line precedes the title; counts precede separately bordered replies. There is no sidebar, breadcrumb trail, or explanatory caption.

Use the product’s Lucide icons and its Heroicons chat-bubble shape, vendored with its MIT notice. Preserve reading order at every width. Displayed controls are decorative and not focusable; the example does not act as a live thread.

### Library and care

The Library sequence, moderation report, and Trail history are static HTML illustrations with real text. Library references and properties follow the product’s content model; the sequence does not imply an automated recommendation feature. The report uses the product’s status and type pills, reporter and handler metadata, subject preview, and action styling. The Trail uses the product’s route and attention icons with scheduled, completed, and needs-attention states, without promising autonomous decisions.

Screenshots and content diagrams have descriptive alt text. Decorative woodland artwork has empty alt text. Product examples have no process-oriented captions such as “Illustrative run history” or “shown in Storyden’s thread layout”.

The connected-assistance link targets `/docs/introduction/robots`, which introduces Robots and links into a setup-first guide for trusted community staff. The companion Trails overview leads to schedule, event and run-history guides; the Cookbook combines both features in actionable recipes.

### Ownership connections and motion

The Storyden mark connects to two callouts:

- **Make it your own:** Custom sites & tools; API · OAuth · Plugins.
- **Add help, if you want it:** Robots & connected tools; MCP · Optional AI.

Both paths are dashed, with muted green for integrations and a warm accent for optional assistance. A smooth curved split has no junction dot. Animated paths stop short of the arrowheads so dashes never overlap them. “Your members, content, and rules.” uses balanced text wrapping.

Dashes flow continuously by 24px every 3.6 seconds with linear timing. Text stays still, and the paths remain visible without drawing or clearing cycles. Playback pauses offscreen and in hidden tabs. Reduced-motion and no-JavaScript rendering show complete static connections. There are no replay buttons or scroll hijacking.

Desktop uses a compact horizontal composition. At widths up to 900px, the source label sits above the mark and the paths flow down into two upright callout columns. Mobile typography and wrapping keep both columns readable at narrow widths.

Narrative and SVG geometry render on the server. `ConnectionMotion.tsx` is the small client boundary responsible for visibility and motion-preference handling. The illustrations make no model calls and require no backend or public API changes.

## Blog, documentation, and shared shell

The shared navigation and footer use the marketing identity. The blog groups its heading, introduction, and RSS controls in one reading block beside the panoramic `wide-mountain.webp` illustration. The header stacks on narrow screens; the artwork keeps its natural 2:1 aspect ratio and becomes full-width on mobile. A dated editorial index follows a fine divider, with two article columns on desktop and one on mobile. Article bodies, reading frames, and URLs remain intact.

Documentation uses restrained brand colours around Fumadocs. Search, sidebar navigation, table of contents, and API reference behaviour remain intact.

Robots and Trails documentation uses feature overviews followed by setup, usage and technical subpages. Cookbook recipes combine built-in capabilities with tested instructions. Screenshots use the existing browser frame, retain at least 2× pixel density at their intended display size, and avoid lossy recompression. Fumadocs image zoom provides a closer view, including a keyboard-accessible expansion control. Markdown prose uses natural paragraphs without manual line wrapping.

## Validation requirements

Run the home production build, TypeScript, scoped lint/format checks, and relevant asset and route checks for substantive changes. Review representative desktop and mobile renders of the landing page, blog index and article, documentation page, and API reference.

Check keyboard navigation, visible focus, contrast, heading order, diagram alternatives, narrow-screen readability, reduced motion, static rendering, and layout stability. Preserve existing assets and product-faithful demonstrations. Use existing local services; do not start development servers.
