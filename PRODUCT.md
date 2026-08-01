# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Storyden serves interest-based internet communities, including clubs, gaming groups, fandoms, and other member-driven spaces. Its adoption audience is the people who operate these communities, while its primary day-to-day users are the members who visit, discuss, and curate content.

Members need a shared place to participate and collaboratively maintain knowledge that remains useful beyond the immediate conversation. Moderators need to keep conversations healthy, manage reports, and maintain quality. Operators and administrators need to establish, host, configure, and extend a community without assembling a collection of disconnected services.

## Product Purpose

Storyden is a platform for building and managing communities and their content. It combines active discussion with durable communal knowledge so a community can run a forum, publish news or a blog, curate links, and build directories, wikis, catalogs, or knowledge bases in one place.

Success means a community can establish its own home, hold useful conversations, and turn what it learns into organized material that remains discoverable and maintainable over time.

## Positioning

Storyden integrates forum discussion with a structured Library for long-lived community knowledge. The Library and discussion system share members, roles, permissions, references, and workflows instead of behaving as separate products. This makes Storyden an alternative for communities that find chat too transient, team workspaces too closed, and legacy forum or wiki software too limited or outdated.

The product is open source, self-hostable, API-driven, extensible through plugins and integrations, and offers AI-assisted features as optional capabilities rather than prerequisites.

## Operating Context

- Community operators deploy and configure an instance, establish roles and permissions, organize discussion categories, moderate activity, and shape the Library.
- Members create threads and replies, react to and reference contributions, discover people and topics, and contribute to shared pages or collections when permitted.
- Communities use the Library to turn conversational knowledge into structured pages, directories, wikis, catalogs, bookmarks, and other maintained resources.
- Developers and operators can automate workflows, build integrations, or provide custom clients through the API, OAuth, MCP tools, and plugin system.

## Capabilities and Constraints

- Discussion includes threads, nested replies, reactions, categories, tags, moderation, reports, and member roles and permissions.
- The Library provides tree-structured pages, rich content, properties, assets, references, revisions, and collaborative curation.
- Community identity and content extend across member profiles, collections, links, search, notifications, and administrative workflows.
- The core deployment is self-hostable with minimal required infrastructure and sane defaults. Optional infrastructure can replace built-in defaults as a deployment grows.
- The system is API-driven and supports custom frontends, OAuth applications, MCP clients, and out-of-process plugins.
- AI-assisted search, organization, and automation are optional; the core community and knowledge workflows must remain useful without them.
- Security, extensibility, and reliable ownership of community data are durable product requirements.
- Storyden is white-label and must adapt to each community through administrator-configured branding, including accent color, images, icons, and custom theme CSS.

## Brand Commitments

- The product name is **Storyden**.
- Storyden is open-source software released under the MIT license.
- Its character is a modern love letter to the community-owned internet and forum culture of the 2000s.
- Product language should be direct, warm, curious, and enthusiastic about people building places around shared interests.
- Storyden must not impose a rigid identity over the communities it hosts; each community's own identity remains primary.
- Existing brand marks and product imagery under `home/public/brand/` and `home/public/` are authoritative assets unless a later explicit rebrand replaces them.

## Evidence on Hand

- `README.md` contains the maintained product overview, capability summary, positioning, deployment example, and current application screenshot.
- `home/content/docs/introduction/` documents the product model, discussion system, Library, members, moderation, API, OAuth, plugins, MCP integration, and deployment workflows.
- `home/public/brand/` contains the existing Storyden icon and logomark variants.
- The working application under `web/` contains implemented member, discussion, Library, collection, link, moderation, administration, search, notification, Robot, and account workflows.
- No testimonials, named customer claims, adoption metrics, or performance benchmarks were confirmed for future product storytelling; they must not be fabricated.

## Product Principles

1. **Let communities own their home.** Preserve self-hosting, open access to the code, extensibility, and control over community data.
2. **Turn conversation into collective memory.** Connect active discussion to structured knowledge so useful contributions do not disappear into a chronological archive.
3. **Make the small deployment complete.** Keep the core product useful with minimal infrastructure while allowing larger communities to adopt specialized services.
4. **Keep intelligence optional and accountable.** AI may reduce organizational work or unlock new workflows, but essential community participation cannot depend on it.
5. **Modernize without discarding forum culture.** Preserve the depth, ownership, and identity of enduring internet communities while improving security, usability, and deployment.

## Accessibility & Inclusion

Storyden targets WCAG 2.2 Level AA across its web product. Core community participation, content creation, moderation, navigation, and account workflows must be operable with a keyboard and assistive technology, support visible focus and sufficient contrast, respect reduced-motion preferences, and remain usable across responsive layouts, user appearance preferences, and lower-powered devices.
