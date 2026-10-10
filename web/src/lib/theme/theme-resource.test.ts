import assert from "node:assert/strict";
import { test } from "vitest";

import { createThemeResourceResponse } from "./theme-resource";

test("returns a cacheable theme resource with a strong content ETag", async () => {
  const response = createThemeResourceResponse(
    new Request("https://community.example/theme.css"),
    "hello",
    "text/css",
  );

  assert.strictEqual(response.status, 200);
  assert.strictEqual(response.headers.get("Cache-Control"), "public, no-cache");
  assert.strictEqual(
    response.headers.get("Content-Type"),
    "text/css; charset=utf-8",
  );
  assert.strictEqual(
    response.headers.get("ETag"),
    '"sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="',
  );
  assert.strictEqual(await response.text(), "hello");
});

test("returns no body when the theme resource ETag matches", async () => {
  const request = new Request("https://community.example/theme.js", {
    headers: {
      "If-None-Match":
        '"sha256-previous", W/"sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="',
    },
  });

  const response = createThemeResourceResponse(
    request,
    "hello",
    "application/javascript",
  );

  assert.strictEqual(response.status, 304);
  assert.strictEqual(
    response.headers.get("ETag"),
    '"sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="',
  );
  assert.strictEqual(await response.text(), "");
});
