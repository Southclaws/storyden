# AGENTS.md

## Storybook

Storybook is the canonical place to work on `web` UI components in isolation.
When changing components under `src/components/ui/**`, add or update the
co-located `*.stories.tsx` file in the component folder.

Storybook is also the directional design-system authority, not merely a prop
catalogue. Before composing a new screen or choosing visual variants, read
`Foundations/Principles/Storyden` and the relevant component/composition
stories. Treat those stories as the approved decision rules for hierarchy,
layering, density, actions, navigation, feedback, and screen composition.

- Use stories to exercise all Panda recipe variants exposed by the component,
  especially `size`, `variant`, `kind`, `shape`, and slot recipes.
- Prefer public component imports from the folder `index.ts`; avoid reaching
  into `*.internal.tsx` from stories unless the component only exposes slot
  primitives that way.
- Keep stories realistic and lightweight. Avoid importing whole feature screens
  or API-backed flows for UI primitive stories.
- For icons and other unusual exports, constrain the Storybook preview to match
  normal product usage instead of changing the component export just for the
  catalogue.

Use explicit slash-delimited story titles so the intentionally flat UI source
directory still produces a predictable Storybook hierarchy:

- `Foundations/Tokens/<Token group>`, `Foundations/Icons/<Catalogue>`, and
  `Foundations/Layout/<System>` document the raw design-system foundations.
- `Foundations/Principles/<Guide>` documents why and when those foundations
  are used in Storyden. Principles precede token or component choice.
- `Components/<Role>/<Component>` contains isolated public UI components. Use
  the established role folders: `Actions`, `Forms`, `Navigation`, `Feedback`,
  `Overlays`, `Data Display`, `Layout`, and `Typography`.
- `Compositions/<Domain>/<Pattern>` contains realistic arrangements of several
  components without importing a complete product screen.
- `Screens/<Screen>` contains complete representative screens and shell
  examples.

Do not introduce another top-level category or a miscellaneous component folder
without first agreeing the taxonomy. The canonical top-level ordering lives in
`.storybook/preview.tsx`.

Do not start the long-running Storybook dev server from an agent unless the user
explicitly asks. If the user already has Storybook running, use it for visual
inspection. Otherwise validate with the static build.

## Frontend Commands

Use Node 24 for frontend work:

```sh
eval "$(fnm env --shell zsh)"
fnm use 24
```

Useful validation commands:

```sh
pnpm --dir web exec tsc --noEmit
pnpm --dir web exec prettier --check "src/components/ui/**/*.stories.tsx"
pnpm --dir web build-storybook
pnpm --dir web exec vitest run src/components/ui
```

Use the targeted commands that match the change. For broad UI component changes,
prefer at least TypeScript and Storybook build before handing back.

## Component Structure

`src/components/ui` is intentionally flat: one component folder per public UI
component or primitive group.

## Component Reuse

Always prefer an existing component, recipe, or established composition before
creating a new component. Search `src/components/ui`, its Storybook stories, and
representative product call sites before writing route-local UI from raw
elements and bespoke CSS.

Before editing a product screen:

1. Read `Foundations/Principles/Component Selection` in Storybook, or inspect
   `src/stories/foundations/ComponentSelection.stories.tsx` when Storybook is
   not running.
2. Search components by semantic intent, not by the markup you expect to write.
   For example, search for status, feedback, member identity, relative time,
   multiline input, page header, or surface before searching for Badge, Box, or
   `styled.*`.
3. Read the closest component stories and at least two representative product
   call sites. Stories define the supported role and boundary; call sites show
   how the component participates in real layouts.
4. State which existing primitives and compositions the screen will reuse. If
   a semantic role is missing, identify that gap before implementing the
   feature-local version.

`styled.*`, `Box`, `Stack`, and generated Panda recipes are implementation
tools, not evidence that the design system lacks a component. Do not rebuild a
semantic component from these primitives when Storybook already documents that
role.

Do not create a new visual component, including a feature-local or route-local
component, without first confirming the decision with the user. Present the
closest existing components, explain why they cannot satisfy the requirement,
and agree where the new component belongs. A component that owns new layout,
spacing, sizing, states, or responsive presentation is a new visual component
even when it is small, uses BEM classes, or is only consumed once.

Feature-specific compositions are allowed without introducing a new primitive
when they only arrange existing canonical components and do not establish their
own visual language. Prefer extending an existing recipe or extracting a shared
composition when multiple features need the same pattern.

Do not extract a component only to shorten JSX or name a small fragment. Keep
one-off labels, counts, and simple arrangements inline. A feature-local
component earns its boundary by owning a domain concept, substantial behavior,
or a coherent section of the screen. Repeated semantic styling or behavior
across three or more call sites is a signal to stop and extract a shared
component before adding another local implementation.

When adding or extending a shared component, its Storybook documentation must
state:

- the semantic role it owns;
- when to use it;
- the nearest component or primitive people might confuse it with;
- when not to use it;
- realistic variants and edge cases.

Component folders generally follow this shape:

- `index.ts` re-exports the public API so consumer import paths stay stable.
- `<ComponentName>.recipe.ts` contains the Panda recipe when the component has
  one.
- `<ComponentName>.internal.tsx` contains the lower-level implementation.
- `<ComponentName>.tsx` may contain a closed/simple variant when useful.
- `<ComponentName>.field.tsx` contains type-safe React Hook Form adapters for
  field variants.
- `<ComponentName>.stories.tsx` documents the component in Storybook.
- `<ComponentName>.test.tsx` covers complex behavior when needed.

Do not reintroduce a shared `ui/form` component directory. Form wrappers for an
existing component live inside that component's folder; standalone form
primitives such as labels and error text live in their own flat folders.
