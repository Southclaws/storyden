import { test } from "uvu";
import * as assert from "uvu/assert";

import { createThemeResourceResponse } from "./theme-resource";

test("returns a cacheable theme resource with a strong content ETag", async () => {
  const response = createThemeResourceResponse(
    new Request("https://community.example/theme.css"),
    "hello",
    "text/css",
  );

  assert.is(response.status, 200);
  assert.is(response.headers.get("Cache-Control"), "public, no-cache");
  assert.is(response.headers.get("Content-Type"), "text/css; charset=utf-8");
  assert.is(
    response.headers.get("ETag"),
    '"sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="',
  );
  assert.is(await response.text(), "hello");
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

  assert.is(response.status, 304);
  assert.is(
    response.headers.get("ETag"),
    '"sha256-LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ="',
  );
  assert.is(await response.text(), "");
});

test.run();
