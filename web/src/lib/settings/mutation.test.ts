import assert from "node:assert/strict";
import { test } from "vitest";

import type { AdminSettingsProps, Info } from "@/api/openapi-schema";
import {
  AuthMode,
  InstanceCapability,
  RegistrationMode,
} from "@/api/openapi-schema";

import { mergeAdminSettingsIntoInfo } from "./mutation";

function baseInfo(overrides: Partial<Info> = {}): Info {
  return {
    title: "Storyden",
    description: "A forum for the modern age.",
    content: "",
    accent_colour: "#123456",
    onboarding_status: "complete",
    authentication_mode: AuthMode.handle,
    registration_mode: RegistrationMode.public,
    capabilities: [InstanceCapability.robots],
    metadata: {},
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
    metadata: {},
    api_address: "http://localhost:8000",
    web_address: "http://localhost:3000",
    ...overrides,
  };
}

test("admin settings updates preserve public-only info fields", () => {
  const current = baseInfo({
    metadata: { navigation: { items: [] } },
  });
  const updated = mergeAdminSettingsIntoInfo(
    current,
    baseAdminSettings({
      title: "Updated Storyden",
      metadata: { navigation: { items: [{ type: "robots" }] } },
    }),
  );

  assert.deepStrictEqual(updated.capabilities, [InstanceCapability.robots]);
  assert.strictEqual(updated.onboarding_status, "complete");
  assert.strictEqual(updated.title, "Updated Storyden");
  assert.deepStrictEqual(updated.metadata, {
    navigation: { items: [{ type: "robots" }] },
  });
});

test("admin settings updates use capabilities when the response includes them", () => {
  const updated = mergeAdminSettingsIntoInfo(
    baseInfo(),
    baseAdminSettings({ capabilities: [] }),
  );

  assert.deepStrictEqual(updated.capabilities, []);
});
