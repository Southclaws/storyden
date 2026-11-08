import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { Account, NotificationSettings } from "@/api/openapi-schema";
import { parseMemberSettings } from "@/lib/settings/member-settings";
import { DefaultSettings, type Settings } from "@/lib/settings/settings";

import { MemberNotificationSettings } from "./MemberNotificationSettings";

const updateAccount = vi.fn();
const mutateAccount = vi.fn();
const updateSettings = vi.fn();
const mutateSettings = vi.fn();

vi.mock("@/api/client", () => ({
  handle: async (action: () => Promise<void>) => action(),
}));

vi.mock("@/api/openapi-client/accounts", () => ({
  accountUpdate: (...args: unknown[]) => updateAccount(...args),
  useAccountGet: () => ({ mutate: mutateAccount }),
}));

vi.mock("@/api/openapi-client/notifications", () => ({
  notificationSettingsUpdate: (...args: unknown[]) => updateSettings(...args),
  useNotificationSettingsGet: () => ({ mutate: mutateSettings }),
}));

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
  meta: { editor: { mode: "markdown" }, custom: "keep" },
  verified_status: "none",
  kind: "human",
  email_addresses: [],
  auth_services: [],
  admin: false,
};

const settings = { ...DefaultSettings, web_push_public_key: "public-key" };
const events = {
  thread_reply: { inbox: true, push: true },
  reply_to_reply: { inbox: true, push: true },
  follow: { inbox: false, push: false },
};
const notificationSettings: NotificationSettings = {
  push: { enabled: true },
  events,
};

function renderSettings(
  preferences: NotificationSettings = notificationSettings,
  site: Settings = settings,
  member: Account = account,
) {
  return render(
    <MemberNotificationSettings
      session={parseMemberSettings(member)}
      settings={site}
      notificationSettings={preferences}
    />,
  );
}

describe("MemberNotificationSettings", () => {
  beforeEach(() => {
    updateAccount.mockResolvedValue(account);
    mutateAccount.mockResolvedValue(account);
    updateSettings.mockImplementation(
      async (value: NotificationSettings) => value,
    );
    mutateSettings.mockResolvedValue(undefined);
  });

  it("saves delivery settings separately and preserves unrelated account metadata", async () => {
    const user = userEvent.setup();
    renderSettings();
    await user.click(screen.getByRole("button", { name: "Replies: Inbox" }));
    await user.click(
      screen.getByRole("checkbox", {
        name: "Mark notifications as read when opened",
      }),
    );
    expect(
      screen.getByRole("button", { name: "Replies: Push" }),
    ).toHaveAttribute("aria-pressed", "false");
    expect(updateSettings).not.toHaveBeenCalled();
    await user.click(screen.getByRole("button", { name: "Save" }));

    const expected = {
      push: { enabled: true },
      events: {
        follow: { inbox: false, push: false },
        thread_reply: { inbox: false, push: false },
        reply_to_reply: { inbox: false, push: false },
      },
    };
    await waitFor(() => {
      expect(updateSettings).toHaveBeenCalledWith(expected);
      expect(mutateSettings).toHaveBeenCalledWith(expected, {
        revalidate: false,
      });
      expect(updateAccount).toHaveBeenCalledWith({
        meta: {
          editor: { mode: "markdown" },
          custom: "keep",
          notifications: { markReadOnClick: false },
        },
      });
      expect(mutateAccount).toHaveBeenCalledWith(account, {
        revalidate: false,
      });
    });
  });

  it("enables Inbox with Push and preserves Inbox when Push is turned off", async () => {
    const user = userEvent.setup();
    renderSettings({
      push: { enabled: true },
      events: {
        ...events,
        thread_reply: { inbox: false, push: false },
        reply_to_reply: { inbox: false, push: false },
      },
    });
    const inbox = screen.getByRole("button", { name: "Replies: Inbox" });
    const push = screen.getByRole("button", { name: "Replies: Push" });
    await user.click(push);
    expect(inbox).toHaveAttribute("aria-pressed", "true");
    expect(push).toHaveAttribute("aria-pressed", "true");
    await user.click(push);
    expect(inbox).toHaveAttribute("aria-pressed", "true");
    expect(push).toHaveAttribute("aria-pressed", "false");
    await user.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() => {
      expect(updateSettings).toHaveBeenCalledWith({
        push: { enabled: true },
        events: {
          ...events,
          thread_reply: { inbox: true, push: false },
          reply_to_reply: { inbox: true, push: false },
        },
      });
    });
    expect(updateAccount).not.toHaveBeenCalled();
  });

  it("preserves mixed per-event choices until the group is changed", async () => {
    const user = userEvent.setup();
    const mixed = {
      push: { enabled: true },
      events: { ...events, thread_reply: { inbox: false, push: false } },
    };
    renderSettings(mixed);
    const inbox = screen.getByRole("button", { name: "Replies: Inbox" });
    const push = screen.getByRole("button", { name: "Replies: Push" });
    expect(inbox).toHaveAttribute("aria-pressed", "mixed");
    expect(push).toHaveAttribute("aria-pressed", "mixed");
    await user.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() => expect(updateSettings).toHaveBeenCalledWith(mixed));
    await user.click(push);
    expect(inbox).toHaveAttribute("aria-pressed", "true");
    expect(push).toHaveAttribute("aria-pressed", "true");
    await user.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() =>
      expect(updateSettings).toHaveBeenLastCalledWith(notificationSettings),
    );
  });

  it("pauses push without clearing choices or disabling Inbox controls", async () => {
    const user = userEvent.setup();
    renderSettings();
    await user.click(
      screen.getByRole("checkbox", { name: "Receive push notifications" }),
    );
    expect(
      screen.getByRole("button", { name: "Replies: Push" }),
    ).toBeDisabled();
    expect(
      screen.getByRole("button", { name: "Replies: Inbox" }),
    ).toBeEnabled();
    await user.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() => {
      expect(updateSettings).toHaveBeenCalledWith({
        events,
        push: { enabled: false },
      });
    });
    expect(updateAccount).not.toHaveBeenCalled();
  });

  it("shows Library and Trail controls when those features are available", () => {
    renderSettings(
      notificationSettings,
      { ...settings, capabilities: ["robots"] },
      {
        ...account,
        roles: [
          {
            id: "role-1",
            createdAt: account.createdAt,
            updatedAt: account.updatedAt,
            name: "Member",
            colour: "blue",
            permissions: ["SUBMIT_LIBRARY_NODE_CHANGES", "MANAGE_TRAILS"],
            badge: false,
            default: false,
          },
        ],
      },
    );
    expect(
      screen.getByRole("button", { name: "Library contributions: Inbox" }),
    ).toBeEnabled();
    expect(
      screen.getByRole("button", { name: "Trail alerts: Inbox" }),
    ).toBeEnabled();
  });

  it("does not expose unavailable features or legacy event types", () => {
    renderSettings({ events: {}, push: { enabled: true } });
    const names = screen
      .getAllByRole("button", { pressed: false })
      .map((control) => control.getAttribute("aria-label"));
    expect(names).toEqual([
      "Replies: Inbox",
      "Replies: Push",
      "Reports and moderation: Inbox",
      "Reports and moderation: Push",
      "All new threads: Inbox",
      "All new threads: Push",
      "All replies: Inbox",
      "All replies: Push",
      "All reactions: Inbox",
      "All reactions: Push",
    ]);
  });

  it("keeps Inbox filtering available when push is disabled on the site", async () => {
    const user = userEvent.setup();
    renderSettings(notificationSettings, DefaultSettings);
    expect(
      screen.queryByRole("checkbox", { name: "Receive push notifications" }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Replies: Push" }),
    ).toBeDisabled();
    await user.click(screen.getByRole("button", { name: "Replies: Inbox" }));
    expect(
      screen.getByRole("button", { name: "Replies: Push" }),
    ).toHaveAttribute("aria-pressed", "false");
    await user.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() => {
      expect(updateSettings).toHaveBeenCalledWith({
        push: { enabled: true },
        events: {
          ...events,
          thread_reply: { inbox: false, push: false },
          reply_to_reply: { inbox: false, push: false },
        },
      });
    });
  });
});
