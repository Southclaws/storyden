import assert from "node:assert/strict";
import { test } from "vitest";

import type {
  AdminSettingsProps,
  Info,
  MessageOfTheDay,
} from "@/api/openapi-schema";
import { AuthMode, RegistrationMode } from "@/api/openapi-schema";

import { DefaultFeedConfig } from "./feed";
import { DefaultNavigationConfig } from "./navigation";
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
    registration_mode: RegistrationMode.public,
    capabilities: [],
    metadata: DefaultFrontendConfig,
    api_address: "http://localhost:8000",
    web_address: "http://localhost:3000",
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
    registration_mode: RegistrationMode.public,
    capabilities: [],
    metadata: DefaultFrontendConfig,
    api_address: "http://localhost:8000",
    web_address: "http://localhost:3000",
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

  assert.deepStrictEqual(parsed.metadata.feed, {
    blocks: [
      { type: "categories", layout: "grid" },
      { type: "quick-share", showCategorySelect: false },
      { type: "threads", source: "uncategorised" },
    ],
  });
  assert.deepStrictEqual(parsed.metadata.editor, { mode: "richtext" });
  assert.deepStrictEqual(parsed.metadata.navigation, DefaultNavigationConfig);
});

test("parseSettings uses the fresh feed for absent or empty metadata", () => {
  const absent = parseSettings(baseInfo({ metadata: undefined }));
  const empty = parseSettings(baseInfo({ metadata: {} }));

  assert.deepStrictEqual(absent.metadata.feed, DefaultFeedConfig);
  assert.deepStrictEqual(empty.metadata.feed, DefaultFeedConfig);
});

test("parseSettings defaults a missing feed without discarding other metadata", () => {
  const parsed = parseSettings(
    baseInfo({
      metadata: {
        editor: { mode: "markdown" },
        signatures: { enabled: false, maxHeight: 320 },
      },
    }),
  );

  assert.deepStrictEqual(parsed.metadata, {
    feed: DefaultFeedConfig,
    navigation: DefaultNavigationConfig,
    editor: { mode: "markdown" },
    signatures: { enabled: false, maxHeight: 320 },
  });
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
  assert.deepStrictEqual(parsed.metadata, DefaultFrontendConfig);
});

test("parseAdminSettings fills missing editor with defaults", () => {
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

  assert.deepStrictEqual(parsed.metadata.feed, {
    blocks: [
      { type: "quick-share", showCategorySelect: true },
      { type: "threads", source: "all" },
    ],
  });
  assert.deepStrictEqual(parsed.metadata.editor, { mode: "richtext" });
  assert.deepStrictEqual(parsed.metadata.navigation, DefaultNavigationConfig);
});

test("parseSettings keeps valid motd metadata type", () => {
  const parsed = parseSettings(
    baseInfo({
      motd: {
        content: "<p>hi</p>",
        metadata: { type: "alert" },
      },
    }),
  );

  assert.deepStrictEqual(parsed.motd?.metadata?.type, "alert");
});

test("parseSettings drops invalid motd metadata type", () => {
  const parsed = parseSettings(
    baseInfo({
      motd: {
        content: "<p>hi</p>",
        metadata: { type: "invalid" } as unknown as Info["metadata"],
      } as MessageOfTheDay,
    }),
  );

  assert.deepStrictEqual(parsed.motd?.metadata, undefined);
});
