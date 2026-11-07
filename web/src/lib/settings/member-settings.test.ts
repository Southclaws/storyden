import { describe, expect, it, vi } from "vitest";

import type { Account } from "@/api/openapi-schema";

import { parseMemberSettings } from "./member-settings";
import { DefaultFrontendConfig } from "./settings";

const account: Account = {
  id: "member-1",
  createdAt: "2026-11-07T12:00:00Z",
  updatedAt: "2026-11-07T12:00:00Z",
  joined: "2026-11-07T12:00:00Z",
  handle: "member",
  name: "Member",
  roles: [],
  bio: "",
  links: [],
  meta: {},
  verified_status: "none",
  kind: "human",
  email_addresses: [],
  auth_services: [],
  admin: false,
};

describe("parseMemberSettings", () => {
  it.each([{}, { editor: { mode: "markdown" } }, { notifications: {} }])(
    "marks opened notifications as read by default for existing metadata %j",
    (meta) => {
      const member = parseMemberSettings({ ...account, meta });

      expect(member.meta.notifications.markReadOnClick).toBe(true);
    },
  );

  it("preserves an explicit false and other account metadata", () => {
    const member = parseMemberSettings({
      ...account,
      meta: {
        editor: { mode: "markdown" },
        notifications: { markReadOnClick: false },
        custom: { colour: "blue" },
      },
    });

    expect(member.meta).toEqual({
      editor: { mode: "markdown" },
      notifications: { markReadOnClick: false },
      custom: { colour: "blue" },
    });
  });

  it("inherits the global editor without changing notification defaults", () => {
    const member = parseMemberSettings(account, {
      ...DefaultFrontendConfig,
      editor: { mode: "markdown" },
    });

    expect(member.meta).toEqual({
      editor: { mode: "markdown" },
      notifications: { markReadOnClick: true },
    });
  });

  it("falls back to the default for malformed notification preferences", () => {
    vi.spyOn(console, "warn").mockImplementation(() => {});

    const member = parseMemberSettings({
      ...account,
      meta: { notifications: { markReadOnClick: "false" } },
    });

    expect(member.meta.notifications.markReadOnClick).toBe(true);
  });
});
