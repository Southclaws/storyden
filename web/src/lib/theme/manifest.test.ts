import assert from "node:assert/strict";
import { test } from "vitest";

import {
  EMPTY_THEME_MANIFEST,
  parseThemeManifest,
  resolveThemeAssetHref,
} from "./manifest";

const valid = {
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
  assert.strictEqual(
    parsed.stylesheets[0]?.href,
    "https://community.example/api/info/theme/assets/theme.css",
  );
});

test("resolves browser asset requests against the configured API origin", () => {
  assert.strictEqual(
    resolveThemeAssetHref(
      "/api/info/theme/assets/theme.css",
      "http://localhost:8000",
    ),
    "http://localhost:8000/api/info/theme/assets/theme.css",
  );
});

test("fails open for malformed manifests", () => {
  assert.deepStrictEqual(
    parseThemeManifest(
      { ...valid, stylesheets: "not-an-array" },
      "https://community.example/api",
    ),
    EMPTY_THEME_MANIFEST,
  );
});

test("rejects cross-origin and non-theme asset paths", () => {
  const remote = structuredClone(valid);
  remote.stylesheets[0]!.path = "https://evil.example/theme.css";
  assert.deepStrictEqual(
    parseThemeManifest(remote, "https://community.example/api"),
    EMPTY_THEME_MANIFEST,
  );

  const wrongPath = structuredClone(valid);
  wrongPath.stylesheets[0]!.path = "/api/assets/theme.css";
  assert.deepStrictEqual(
    parseThemeManifest(wrongPath, "https://community.example/api"),
    EMPTY_THEME_MANIFEST,
  );
});
