import { renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { Permission, type ReadStatus } from "@/api/openapi-schema";

import { useThreadMenu } from "./useThreadMenu";

let sessionMock: any;

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
