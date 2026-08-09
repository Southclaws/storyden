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
