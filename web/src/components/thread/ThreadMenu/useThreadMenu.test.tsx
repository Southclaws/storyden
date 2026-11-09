import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { Permission, type ReadStatus } from "@/api/openapi-schema";

import { useThreadMenu } from "./useThreadMenu";

let sessionMock: any;
const updateLifecycleState = vi.fn();
const updateLocked = vi.fn();

vi.mock("@/lib/thread/lifecycle", () => ({
  useThreadLifecycle: () => ({ updateLifecycleState, updateLocked }),
}));

vi.mock("next/navigation", () => ({
  usePathname: () => "/t/thread-1",
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock("nuqs", () => ({
  parseAsBoolean: {},
  useQueryState: () => [undefined, vi.fn()],
}));

vi.mock("@/auth", () => ({
  useSession: () => sessionMock,
}));

vi.mock("@/components/site/useConfirmation", () => ({
  useConfirmation: () => ({
    isConfirming: false,
    handleConfirmAction: vi.fn(),
    handleCancelAction: vi.fn(),
  }),
}));

vi.mock("@/lib/feed/mutation", () => ({
  useFeedMutations: () => ({
    deleteThread: vi.fn(),
    updateThread: vi.fn(),
    revalidate: vi.fn(),
  }),
}));

vi.mock("@/lib/report/useReportContext", () => ({
  useReportContext: () => ({
    resolveReport: vi.fn(),
  }),
}));

vi.mock("@/utils/client", () => ({
  useShare: () => false,
}));

vi.mock("@/utils/useCopyToClipboard", () => ({
  useCopyToClipboard: () => [undefined, vi.fn()],
}));

describe("useThreadMenu", () => {
  beforeEach(() => {
    updateLocked.mockReset();
    sessionMock = {
      id: "moderator-1",
      roles: [{ id: "role-1", permissions: [Permission.MANAGE_POSTS] }],
    };
  });

  it("disables editing for deleted threads", () => {
    const { result } = renderHook(() =>
      useThreadMenu({
        thread: {
          id: "thread-1",
          slug: "thread-1",
          author: { id: "author-1", name: "Author" },
          deletedAt: "2026-03-23T00:00:00.000Z",
        } as any,
        editingEnabled: true,
      }),
    );

    expect(result.current.isEditingEnabled).toBe(false);
    expect(result.current.canUpdateLifecycle).toBe(false);
  });

  it.each([
    { name: "author", id: "author-1", permissions: [], allowed: true },
    {
      name: "member with Manage Posts",
      id: "manager-1",
      permissions: [Permission.MANAGE_POSTS],
      allowed: true,
    },
    { name: "another member", id: "other-1", permissions: [], allowed: false },
    { name: "guest", id: undefined, permissions: [], allowed: false },
  ])("lifecycle access for $name", ({ id, permissions, allowed }) => {
    sessionMock = id ? { id, roles: [{ permissions }] } : undefined;
    const { result } = renderHook(() =>
      useThreadMenu({
        thread: {
          id: "thread-1",
          slug: "thread-1",
          author: { id: "author-1" },
          lifecycle_state: "open",
        } as any,
      }),
    );

    expect(result.current.canUpdateLifecycle).toBe(allowed);
    expect(result.current.isThreadResolved).toBe(false);

    if (!allowed) {
      act(() => {
        void result.current.handlers.handleLifecycleChange("resolved");
      });
      expect(updateLifecycleState).not.toHaveBeenCalled();
    }
  });

  it.each([
    { name: "author", permissions: [], allowed: false },
    { name: "manager", permissions: [Permission.MANAGE_POSTS], allowed: true },
    {
      name: "administrator",
      permissions: [Permission.ADMINISTRATOR],
      allowed: true,
    },
  ])("lock access for $name", async ({ permissions, allowed }) => {
    sessionMock = { id: "author-1", roles: [{ permissions }] };
    const { result } = renderHook(() =>
      useThreadMenu({
        thread: {
          id: "thread-1",
          slug: "thread-1",
          author: { id: "author-1" },
        } as any,
      }),
    );

    expect(result.current.canLockThread).toBe(allowed);
    await act(async () => {
      await result.current.handlers.handleLockChange(true);
    });
    expect(updateLocked).toHaveBeenCalledTimes(allowed ? 1 : 0);
  });

  it("disables lifecycle changes while saving and allows reopening afterward", async () => {
    let finish!: () => void;
    updateLifecycleState.mockImplementationOnce(
      () =>
        new Promise<void>((resolve) => {
          finish = resolve;
        }),
    );
    const { result, rerender } = renderHook(
      ({ state }: { state: "open" | "resolved" }) =>
        useThreadMenu({
          thread: {
            id: "thread-1",
            slug: "thread-1",
            author: { id: "author-1" },
            lifecycle_state: state,
          } as any,
        }),
      { initialProps: { state: "open" } },
    );

    let pending!: Promise<void>;
    act(() => {
      pending = result.current.handlers.handleLifecycleChange("resolved");
    });
    expect(result.current.isUpdatingLifecycle).toBe(true);

    await act(async () => {
      await result.current.handlers.handleLifecycleChange("resolved");
    });
    expect(updateLifecycleState).toHaveBeenCalledTimes(1);

    await act(async () => {
      finish();
      await pending;
    });
    expect(result.current.isUpdatingLifecycle).toBe(false);

    rerender({ state: "resolved" });
    expect(result.current.isThreadResolved).toBe(true);

    await act(async () => {
      await result.current.handlers.handleLifecycleChange("open");
    });
    expect(updateLifecycleState).toHaveBeenLastCalledWith("open");
  });

  it.each<{
    state: string;
    readStatus?: ReadStatus;
    isUnread: boolean;
  }>([
    { state: "never visited", isUnread: true },
    {
      state: "caught up",
      readStatus: { last_read_at: "2026-11-08T10:00:00Z", replies_since: 0 },
      isUnread: false,
    },
    {
      state: "new replies",
      readStatus: { last_read_at: "2026-11-08T10:00:00Z", replies_since: 5 },
      isUnread: true,
    },
    {
      state: "explicitly marked unread",
      readStatus: {
        last_read_at: "2026-11-08T10:00:00Z",
        replies_since: 0,
        marked_unread_at: "2026-11-08T11:00:00Z",
      },
      isUnread: true,
    },
  ])("recognises $state threads", ({ readStatus, isUnread }) => {
    const { result } = renderHook(() =>
      useThreadMenu({
        thread: {
          id: "thread-1",
          slug: "thread-1",
          author: { id: "author-1", name: "Author" },
          read_status: readStatus,
        } as any,
      }),
    );

    expect(result.current.canUpdateReadStatus).toBe(true);
    expect(result.current.isUnread).toBe(isUnread);
  });
});
