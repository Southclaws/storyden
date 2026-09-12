# Frontend WebMCP tools

This directory is the source of truth for tools implemented by the Storyden
frontend. These tools describe capabilities of the current browser page, not
the Storyden backend API.

Each tool uses the same wrapper shape and `x-storyden-tool` metadata convention
as `api/robots/tools`. Add the tool to `schema.yaml`, then regenerate the
TypeScript contracts and runtime catalogue:

```bash
task generate:webmcp
```

Generated output lives in `web/src/lib/webmcp/tools.generated.ts`. Do not edit
that file by hand.

## Registration lifecycle

Library layout tools are registered only while a Library page is in direct edit
mode.

The `useWebMCP` hooks own registration and unregister automatically when their
React component unmounts. The polyfill supplies the same-document API when the
browser has no native WebMCP implementation and leaves a native implementation
untouched.
