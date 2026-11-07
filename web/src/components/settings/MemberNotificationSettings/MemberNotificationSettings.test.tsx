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

describe("MemberNotificationSettings", () => {
  beforeEach(() => {
    update.mockResolvedValue(account);
    mutate.mockResolvedValue(account);
  });

  it("saves grouped controls as event overrides and preserves other preferences", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: { enabled: true, events: { follow: false } },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    await user.click(screen.getByRole("checkbox", { name: "Replies" }));
    await user.click(
      screen.getByRole("checkbox", {
        name: "Mark notifications as read when opened",
      }),
    );

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
          push: {
            enabled: true,
            events: {
              follow: false,
              thread_reply: false,
              reply_to_reply: false,
            },
          },
        },
      });
      expect(mutate).toHaveBeenCalledWith(account, { revalidate: false });
    });
  });

  it("preserves mixed event overrides when saving an unrelated control", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: {
          enabled: true,
          events: {
            thread_reply: false,
            reply_to_reply: true,
            post_like: false,
          },
        },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    await waitFor(() => {
      expect(
        screen.getByRole("checkbox", { name: "Replies" }),
      ).toBePartiallyChecked();
    });

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: session.notification_settings,
        }),
      );
    });
  });

  it("pauses push without clearing event preferences", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: { enabled: true, events: { thread_reply: false } },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    await user.click(
      screen.getByRole("checkbox", { name: "Receive push notifications" }),
    );

    expect(screen.getByRole("checkbox", { name: "Replies" })).toBeDisabled();

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: {
            push: {
              enabled: false,
              events: { thread_reply: false },
            },
          },
        }),
      );
    });
  });

  it("changes all events in a mixed group without changing other events", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      notification_settings: {
        push: {
          enabled: true,
          events: { thread_reply: false, reply_to_reply: true, follow: false },
        },
      },
    });

    render(
      <MemberNotificationSettings session={session} settings={settings} />,
    );

    const replies = screen.getByRole("checkbox", { name: "Replies" });

    await waitFor(() => {
      expect(replies).toBePartiallyChecked();
    });

    await user.click(replies);
    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith(
        expect.objectContaining({
          notification_settings: {
            push: {
              enabled: true,
              events: {
                thread_reply: true,
                reply_to_reply: true,
                follow: false,
              },
            },
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
      screen.getByRole("checkbox", { name: "Library contributions" }),
    ).toBeChecked();
    expect(
      screen.getByRole("checkbox", { name: "Trail alerts" }),
    ).toBeChecked();
  });

  it("does not expose unavailable features or legacy notification types", () => {
    render(
      <MemberNotificationSettings
        session={parseMemberSettings(account)}
        settings={settings}
      />,
    );

    const labels = screen
      .getAllByRole("checkbox")
      .map(
        (control) =>
          control.getAttribute("aria-label") ??
          control.closest("label")?.textContent,
      );

    expect(labels).toEqual([
      "Mark notifications as read when opened",
      "Receive push notifications",
      "Replies",
      "Reports and moderation",
    ]);
  });

  it("keeps inbox settings available when push is disabled on the site", () => {
    render(
      <MemberNotificationSettings
        session={parseMemberSettings(account)}
        settings={DefaultSettings}
      />,
    );

    expect(
      screen.getByRole("checkbox", {
        name: "Mark notifications as read when opened",
      }),
    ).toBeChecked();
    expect(
      screen.queryByRole("checkbox", { name: "Receive push notifications" }),
    ).not.toBeInTheDocument();
  });
});
