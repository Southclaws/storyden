# AGENTS.md

## Storybook

Storybook is the canonical place to work on `web` UI components in isolation.
When changing components under `src/components/ui/**`, add or update the
co-located `*.stories.tsx` file in the component folder.

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

Component folders generally follow this shape:

- `index.ts` re-exports the public API so consumer import paths stay stable.
- `<ComponentName>.recipe.ts` contains the Panda recipe when the component has
  one.
- `<ComponentName>.internal.tsx` contains the lower-level implementation.
- `<ComponentName>.tsx` may contain a closed/simple variant when useful.
- `<ComponentName>.form.tsx` contains React Hook Form wrappers for form variants.
- `<ComponentName>.stories.tsx` documents the component in Storybook.
- `<ComponentName>.test.tsx` covers complex behavior when needed.

Do not reintroduce a shared `ui/form` component directory. Form wrappers for an
existing component live inside that component's folder; standalone form
primitives such as labels and error text live in their own flat folders.
