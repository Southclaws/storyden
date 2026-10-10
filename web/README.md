# Storyden frontend reference implementation

The official Storyden frontend implementation. Production ready and implements most of the Storyden functionality.

## Stack

- Next.js
- Panda CSS

Run frontend commands from `web/` with the Node version specified in `.nvmrc`.

## Custom frontend

If you are implementing your own frontend, use this as a reference and feel free to borrow code and components.

## Checking

While editing this part of the codebase, make sure to utilise the linting, checking and formatting tools:

- `pnpm test` runs all unit and React tests with Vitest.
- `pnpm test:unit` runs `src/**/*.test.ts` in Node.
- `pnpm test:ui` runs `src/**/*.test.tsx` in jsdom with React Testing Library setup.
- `pnpm typecheck` runs TypeScript separately.
- `pnpm lint` checks handwritten code with Oxlint.
- `pnpm lint:fix` applies safe automatic fixes. Review the resulting diff.
- `pnpm format` sorts import members with Oxlint, then formats files and groups imports with Oxfmt.
- `pnpm format:check` checks both without writing changes.

`.oxlintrc.json` enables TypeScript, React, Next.js and accessibility checks. Definite mistakes such as conditional hooks, duplicate keys and unsafe optional chaining are errors. Other correctness findings are warnings. Warnings do not fail the command; errors do.

Unused-variable checks and experimental React Compiler rules are disabled. The semantic-tag preference is also disabled because it does not understand all our UI wrappers. Storybook render/decorator callbacks are exempt from the hook naming rule; application components still get that check. Generated API clients, schemas, Panda output and build artifacts are excluded.

The repository-root config limits editor linting to `web/`. Rules and exclusions live in the frontend config, shared by the CLI and language server.
