import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { NotificationListResult } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import type { Member } from "@/lib/settings/member-settings";

import { useNotifications } from "./useNotifications";

const notificationUpdate = vi.fn();
const updateMany = vi.fn();
const mutate = vi.fn();

vi.mock("@/api/client", () => ({
  handle: async (action: () => Promise<void>) => action(),
}));

vi.mock("@/auth", () => ({
  useSession: vi.fn(),
}));

vi.mock("@/api/openapi-client/notifications", () => ({
  notificationUpdate: (...args: unknown[]) => notificationUpdate(...args),
  useNotificationList: () => ({ data: notifications, mutate }),
  useNotificationUpdateMany: () => ({ trigger: updateMany }),
}));

const notifications: NotificationListResult = {
  current_page: 1,
  page_size: 20,
  results: 1,
  total_pages: 1,
  notifications: [
    {
      id: "notification-1",
      created_at: "2026-11-07T12:00:00Z",
      event: "trail_run_attention",
      status: "unread",
    },
  ],
};

const session: Member = {
  id: "member-1",
  createdAt: "2026-11-07T12:00:00Z",
  updatedAt: "2026-11-07T12:00:00Z",
  joined: "2026-11-07T12:00:00Z",
  handle: "member",
  name: "Member",
  roles: [],
  bio: "",
  links: [],
  meta: {
    editor: { mode: "richtext" },
    notifications: { markReadOnClick: true },
  },
  verified_status: "none",
  kind: "human",
  email_addresses: [],
  auth_services: [],
  admin: false,
};

describe("useNotifications", () => {
  beforeEach(() => {
    vi.mocked(useSession).mockReturnValue(session);
    notificationUpdate.mockResolvedValue(undefined);
    updateMany.mockResolvedValue(undefined);
    mutate.mockResolvedValue(undefined);
  });

  it("marks opened notifications as read when enabled", async () => {
    const { result } = renderHook(() => useNotifications({ status: "unread" }));

    await act(async () => {
      await result.current.handlers?.handleOpen("notification-1");
    });

    expect(notificationUpdate).toHaveBeenCalledWith("notification-1", {
      status: "read",
    });
  });

  it("leaves opened notifications unread when disabled", async () => {
    vi.mocked(useSession).mockReturnValue({
      ...session,
      meta: {
        ...session.meta,
        notifications: { markReadOnClick: false },
      },
    });

    const { result } = renderHook(() => useNotifications({ status: "unread" }));

    await act(async () => {
      await result.current.handlers?.handleOpen("notification-1");
    });

    expect(notificationUpdate).not.toHaveBeenCalled();
    expect(mutate).not.toHaveBeenCalled();
  });

  it("still supports manual read, unread, and mark-all actions when disabled", async () => {
    vi.mocked(useSession).mockReturnValue({
      ...session,
      meta: {
        ...session.meta,
        notifications: { markReadOnClick: false },
      },
    });

    const { result } = renderHook(() => useNotifications({ status: "unread" }));

    await act(async () => {
      await result.current.handlers?.handleMarkAs("notification-1", "read");
      await result.current.handlers?.handleMarkAs("notification-1", "unread");
      await result.current.handlers?.handleMarkAllAsRead();
    });

    expect(notificationUpdate).toHaveBeenNthCalledWith(1, "notification-1", {
      status: "read",
    });
    expect(notificationUpdate).toHaveBeenNthCalledWith(2, "notification-1", {
      status: "unread",
    });
    expect(updateMany).toHaveBeenCalledWith({
      notifications: [{ id: "notification-1", status: "read" }],
    });
  });

  it("waits for account preferences before marking anything read", async () => {
    vi.mocked(useSession).mockReturnValue(undefined);

    const { result } = renderHook(() => useNotifications({ status: "unread" }));

    await act(async () => {
      await result.current.handlers?.handleOpen("notification-1");
    });

    expect(notificationUpdate).not.toHaveBeenCalled();
  });
});
