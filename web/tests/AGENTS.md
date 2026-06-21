# Playwright Test Guidance

## Shared E2E State

Playwright tests run against one shared backend for the duration of the E2E
harness. This is by design. Do not write tests that assume a pristine database,
fixed resource counts, or that only resources created by the current test exist.

- Create resources with unique names/slugs/handles, usually with a timestamp or
  similarly unique suffix.
- Scope assertions to the specific resource created by the test whenever
  possible, using exact names, IDs, slugs, or test-local text.
- Avoid asserting global list counts unless the test explicitly controls the
  whole dataset. Prefer "contains my resource" or "does not contain my deleted
  resource" over "has exactly N items".
- When selecting from paginated or ordered lists, make sure the test resource is
  discoverable even if other tests have created many resources. Use search,
  filters, deterministic ordering, or API setup that places the resource in the
  relevant page.
- Tests should tolerate unrelated resources from earlier specs, retries, or
  parallel workers. Cleanup is useful, but must not be required for correctness.
- Do not depend on test execution order. A spec should pass when run alone, when
  run after the full suite, and when retried.

## User Behavior

Prefer user-visible interactions over implementation shortcuts. For example, if
a user selects an item through a picker, the test should exercise the picker
rather than navigating directly to the final URL, unless the test is specifically
about route handling.
