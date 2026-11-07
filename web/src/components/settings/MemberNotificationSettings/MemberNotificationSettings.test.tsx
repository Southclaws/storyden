import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { Account } from "@/api/openapi-schema";
import { parseMemberSettings } from "@/lib/settings/member-settings";
import { DefaultSettings } from "@/lib/settings/settings";

import { MemberNotificationSettings } from "./MemberNotificationSettings";

const update = vi.fn();
const mutate = vi.fn();

vi.mock("@/api/client", () => ({
  handle: async (action: () => Promise<void>) => action(),
}));

vi.mock("@/api/openapi-client/accounts", () => ({
  accountUpdate: (...args: unknown[]) => update(...args),
  useAccountGet: () => ({ mutate }),
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

const notificationSettings = { push: { enabled: true }, events };

describe("MemberNotificationSettings", () => {
  beforeEach(() => {
    update.mockResolvedValue(account);
    mutate.mockResolvedValue(account);
  });

  it("saves both channels and preserves unrelated settings", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: notificationSettings,
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    await user.click(screen.getByRole("button", { name: "Replies: Inbox" }));
    await user.click(
      screen.getByRole("checkbox", {
        name: "Mark notifications as read when opened",
      }),
    );

    expect(
      screen.getByRole("button", { name: "Replies: Push" }),
    ).toHaveAttribute("aria-pressed", "false");
    expect(update).not.toHaveBeenCalled();

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith({
        meta: {
          editor: { mode: "markdown" },
          custom: "keep",
          notifications: { markReadOnClick: false },
        },
        notification_settings: {
          push: { enabled: true },
          events: {
            follow: { inbox: false, push: false },
            thread_reply: { inbox: false, push: false },
            reply_to_reply: { inbox: false, push: false },
          },
        },
      });
      expect(mutate).toHaveBeenCalledWith(account, { revalidate: false });
    });
  });

  it("enables Inbox with Push and preserves Inbox when Push is turned off", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: { enabled: true },
        events: {
          ...events,
          thread_reply: { inbox: false, push: false },
          reply_to_reply: { inbox: false, push: false },
        },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

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
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: {
            push: { enabled: true },
            events: {
              ...events,
              thread_reply: { inbox: true, push: false },
              reply_to_reply: { inbox: true, push: false },
            },
          },
        }),
      );
    });
  });

  it("preserves mixed per-event choices until the group is changed", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: { enabled: true },
        events: {
          ...events,
          thread_reply: { inbox: false, push: false },
        },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    const inbox = screen.getByRole("button", { name: "Replies: Inbox" });
    const push = screen.getByRole("button", { name: "Replies: Push" });
    expect(inbox).toHaveAttribute("aria-pressed", "mixed");
    expect(push).toHaveAttribute("aria-pressed", "mixed");

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: session.notification_settings,
        }),
      );
    });

    await user.click(push);
    expect(inbox).toHaveAttribute("aria-pressed", "true");
    expect(push).toHaveAttribute("aria-pressed", "true");
    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenLastCalledWith(
        expect.objectContaining({
          notification_settings: notificationSettings,
        }),
      );
    });
  });

  it("pauses push without clearing choices or disabling Inbox controls", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: notificationSettings,
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

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
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: {
            events,
            push: { enabled: false },
          },
        }),
      );
    });
  });

  it("shows Library and Trail controls when those features are available", () => {
    const session = parseMemberSettings({
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
    });

    render(
      <MemberNotificationSettings
        session={session}
        settings={{ ...settings, capabilities: ["robots"] }}
      />,
    );

    expect(
      screen.getByRole("button", { name: "Library contributions: Inbox" }),
    ).toBeEnabled();
    expect(
      screen.getByRole("button", { name: "Trail alerts: Inbox" }),
    ).toBeEnabled();
  });

  it("does not expose unavailable features or legacy event types", () => {
    render(
      <MemberNotificationSettings
        session={parseMemberSettings(account)}
        settings={settings}
      />,
    );

    const names = screen
      .getAllByRole("button", { pressed: false })
      .map((control) => control.getAttribute("aria-label"));

    expect(names).toEqual([
      "Replies: Inbox",
      "Replies: Push",
      "Reports and moderation: Inbox",
      "Reports and moderation: Push",
    ]);
  });

  it("keeps Inbox filtering available when push is disabled on the site", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: notificationSettings,
    });

    render(
      <MemberNotificationSettings
        session={session}
        settings={DefaultSettings}
      />,
    );

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
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: {
            push: { enabled: true },
            events: {
              ...events,
              thread_reply: { inbox: false, push: false },
              reply_to_reply: { inbox: false, push: false },
            },
          },
        }),
      );
    });
  });
});
