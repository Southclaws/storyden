import { test } from "uvu";
import * as assert from "uvu/assert";

import {
  filterAllowedThemeAssets,
  parseThemeManifest,
} from "./manifest";

test("parseThemeManifest returns defaults for invalid payload", () => {
  const parsed = parseThemeManifest({
    css: "bad",
    scripts: null,
  });

  assert.equal(parsed, {
    css: [],
    scripts: [],
  });
});

test("filterAllowedThemeAssets keeps only web/api origins and resolves relative paths", () => {
  const filtered = filterAllowedThemeAssets(
    {
      css: [
        "/api/info/theme/assets/local-css",
        "/theme.css",
        "https://app.storyden.test/theme.css",
        "https://cdn.evil.test/theme.css",
      ],
      scripts: [
        "/api/info/theme/assets/local-script",
        "https://api.storyden.test/api/info/theme/assets/a",
        "https://evil.example.test/x.js",
      ],
    },
    {
      webAddress: "https://app.storyden.test",
      apiAddress: "https://api.storyden.test",
    },
  );

  assert.equal(filtered.css, [
    "https://api.storyden.test/api/info/theme/assets/local-css",
    "https://app.storyden.test/theme.css",
    "https://app.storyden.test/theme.css",
  ]);
  assert.equal(filtered.scripts, [
    "https://api.storyden.test/api/info/theme/assets/local-script",
    "https://api.storyden.test/api/info/theme/assets/a",
  ]);
});

test.run();
