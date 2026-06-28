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

## Flake Resistance

- Wait for durable state transitions, not just visible text, before starting the
  next user action. For example, after creating the first robot chat message,
  wait until the route has changed from `/robots/chats/new` to the persisted
  `/robots/chats/:id` route before sending another message.
- Reacquire locators inside retryable interactions when the app may remount a
  control during route confirmation, cache refresh, or optimistic updates. A
  button that was enabled a moment ago may detach and reappear disabled if the
  composer clears or remounts.
- Prefer helper functions that retry the complete user action when a remount is
  expected: fill the input, verify the input value, verify the submit control is
  enabled, then click. Avoid retrying only the final click after earlier state
  may have been lost.
- Scope assertions to the latest/current UI element when previous messages,
  retries, or projected tool output can legitimately render similar content.
  Use `.last()` or a message/container-specific locator instead of a global
  strict locator when repeated summaries are valid.
- Treat long Playwright timeouts as a signal to inspect app state assumptions.
  Slow CI usually explains small delays, not a 60 second wait for a disabled
  control or a strict locator matching multiple valid elements.
