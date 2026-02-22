import { test } from "uvu";
import * as assert from "uvu/assert";

import type { AdminSettingsProps, Info } from "@/api/openapi-schema";
import { AuthMode } from "@/api/openapi-schema";

import {
  DefaultFrontendConfig,
  parseAdminSettings,
  parseSettings,
} from "./settings";

function baseInfo(overrides: Partial<Info> = {}): Info {
  return {
    title: "Storyden",
    description: "A forum for the modern age.",
    content: "",
    accent_colour: "#123456",
    onboarding_status: "complete",
    authentication_mode: AuthMode.handle,
    capabilities: [],
    metadata: DefaultFrontendConfig,
    ...overrides,
  };
}

function baseAdminSettings(
  overrides: Partial<AdminSettingsProps> = {},
): AdminSettingsProps {
  return {
    title: "Storyden",
    description: "A forum for the modern age.",
    content: "",
    accent_colour: "#123456",
    authentication_mode: AuthMode.handle,
    capabilities: [],
    metadata: DefaultFrontendConfig,
    ...overrides,
  };
}

test("parseSettings keeps valid metadata and applies nested defaults", () => {
  const parsed = parseSettings(
    baseInfo({
      metadata: {
        feed: {
          layout: { type: "grid" },
          source: { type: "categories" },
        },
      },
    }),
  );

  assert.equal(parsed.metadata.feed, {
    layout: { type: "grid" },
    source: { type: "categories", threadListMode: "uncategorised", quickShare: "enabled" },
  });
  assert.equal(parsed.metadata.editor, { mode: "richtext" });
  assert.equal(parsed.metadata.sidebar, { defaultState: "closed" });
});

test("parseSettings falls back to defaults for invalid metadata", () => {
  const originalWarn = console.warn;
  let warned = false;
  let parsed!: ReturnType<typeof parseSettings>;

  console.warn = () => {
    warned = true;
  };

  try {
    parsed = parseSettings(
      baseInfo({
        metadata: {
          feed: {
            layout: { type: "invalid-layout" },
            source: { type: "threads" },
          },
        } as unknown as Info["metadata"],
      }),
    );
  } finally {
    console.warn = originalWarn;
  }

  assert.ok(warned);
  assert.equal(parsed.metadata, DefaultFrontendConfig);
});

test("parseAdminSettings fills missing editor/sidebar with defaults", () => {
  const parsed = parseAdminSettings(
    baseAdminSettings({
      metadata: {
        feed: {
          layout: { type: "list" },
          source: { type: "threads" },
        },
      },
    }),
  );

  assert.equal(parsed.metadata.feed, {
    layout: { type: "list" },
    source: { type: "threads", quickShare: "enabled" },
  });
  assert.equal(parsed.metadata.editor, { mode: "richtext" });
  assert.equal(parsed.metadata.sidebar, { defaultState: "closed" });
});

test.run();
