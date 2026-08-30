import { test } from "uvu";
import * as assert from "uvu/assert";

import {
  EMPTY_THEME_MANIFEST,
  parseThemeManifest,
  resolveThemeAssetHref,
} from "./manifest";

const valid = {
  api_version: "v1",
  enabled: true,
  revision: "revision-123",
  stylesheets: [
    {
      id: "css00000000000000000",
      filename: "theme.css",
      path: "/api/info/theme/assets/theme.css",
      mime_type: "text/css",
      size: 123,
      integrity: "sha256-dGVzdC1kaWdlc3Q=",
    },
  ],
  scripts: [],
};

test("accepts ordered same-origin immutable theme assets", () => {
  const parsed = parseThemeManifest(valid, "https://community.example/api");
  assert.is(
    parsed.stylesheets[0]?.href,
    "https://community.example/api/info/theme/assets/theme.css",
  );
  assert.is(parsed.revision, "revision-123");
});

test("resolves browser asset requests against the configured API origin", () => {
  assert.is(
    resolveThemeAssetHref(
      "/api/info/theme/assets/theme.css",
      "http://localhost:8000",
    ),
    "http://localhost:8000/api/info/theme/assets/theme.css",
  );
});

test("fails open for malformed manifests", () => {
  assert.equal(
    parseThemeManifest(
      { ...valid, api_version: "v2" },
      "https://community.example/api",
    ),
    EMPTY_THEME_MANIFEST,
  );
});

test("rejects cross-origin and non-theme asset paths", () => {
  const remote = structuredClone(valid);
  remote.stylesheets[0]!.path = "https://evil.example/theme.css";
  assert.equal(
    parseThemeManifest(remote, "https://community.example/api"),
    EMPTY_THEME_MANIFEST,
  );

  const wrongPath = structuredClone(valid);
  wrongPath.stylesheets[0]!.path = "/api/assets/theme.css";
  assert.equal(
    parseThemeManifest(wrongPath, "https://community.example/api"),
    EMPTY_THEME_MANIFEST,
  );
});

test("a disabled manifest does not expose configured assets", () => {
  assert.equal(
    parseThemeManifest(
      { ...valid, enabled: false },
      "https://community.example/api",
    ),
    EMPTY_THEME_MANIFEST,
  );
});

test.run();
