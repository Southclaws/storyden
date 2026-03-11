# Theming Contract (v1)

Storyden v1 custom theming exposes a small, stable frontend contract:

- `data-sd-theme-api="v1"` on the root `<html>` element.
- `sd-*` class names on shared app/layout/navigation wrappers.
- `sd-screen` and `sd-screen--*` route/surface wrapper classes.

## Compatibility policy

- These markers are part of the public theming API surface.
- Renames/removals are considered breaking changes.
- If a class/marker must change, Storyden must ship alias classes/markers for
  at least one non-major release window and note the deprecation in changelog.
- Internal framework classes (for example, generated CSS-module names and Panda
  utility classes) are not part of this contract and may change at any time.
