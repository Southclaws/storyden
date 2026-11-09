import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { Permission } from "@/api/openapi-schema";

import { useReplyMenu } from "./useReplyMenu";

const deleteReply = vi.fn();
const revalidate = vi.fn();
const resolveReport = vi.fn();
const onEdit = vi.fn();
const updateConclusionReply = vi.fn();

vi.mock("@/lib/thread/lifecycle", () => ({
  useThreadLifecycle: () => ({ updateConclusionReply }),
}));

let sessionMock: any;

vi.mock("src/auth", () => ({
  useSession: () => sessionMock,
}));

vi.mock("src/utils/client", () => ({
  useShare: () => false,
}));

vi.mock("@/api/client", () => ({
  handle: async (
    fn: () => Promise<void>,
    opts?: { cleanup?: () => Promise<void> },
  ) => {
    await fn();
    await opts?.cleanup?.();
  },
}));

vi.mock("@/lib/report/useReportContext", () => ({
  useReportContext: () => ({
    resolveReport,
  }),
}));

vi.mock("@/lib/thread/mutation", () => ({
  useThreadMutations: () => ({
    deleteReply,
    revalidate,
  }),
}));

vi.mock("@/lib/thread/undo", () => ({
  withUndo: async ({ action }: { action: () => Promise<void> }) => {
    await action();
  },
}));

vi.mock("@/utils/useCopyToClipboard", () => ({
  useCopyToClipboard: () => [undefined, vi.fn()],
}));

describe("useReplyMenu", () => {
  beforeEach(() => {
    deleteReply.mockReset();
    revalidate.mockReset();
    resolveReport.mockReset();
    onEdit.mockReset();
    sessionMock = {
      id: "moderator-1",
      roles: [{ id: "role-1", permissions: [Permission.MANAGE_POSTS] }],
    };
  });

  it("disables edit and delete actions for deleted replies", () => {
    const { result } = renderHook(() =>
      useReplyMenu({
        thread: { slug: "thread-1" } as any,
        reply: {
          id: "reply-1",
          author: { id: "author-1", name: "Author" },
          body: "Reply body",
          deletedAt: "2026-03-23T00:00:00.000Z",
        } as any,
        onEdit,
      }),
    );

    expect(result.current.isEditingEnabled).toBe(false);
    expect(result.current.isDeletingEnabled).toBe(false);
  });

  it("resolves the active report when deleting from the reply menu", async () => {
    const { result } = renderHook(() =>
      useReplyMenu({
        thread: { slug: "thread-1" } as any,
        reply: {
          id: "reply-1",
          author: { id: "author-1", name: "Author" },
          body: "Reply body",
        } as any,
        onEdit,
      }),
    );

    await act(async () => {
      await result.current.handlers.handleDelete();
    });

    expect(deleteReply).toHaveBeenCalledWith("reply-1");
    expect(resolveReport).toHaveBeenCalledTimes(1);
    expect(revalidate).toHaveBeenCalledTimes(1);
  });

  it.each([
    {
      name: "thread author",
      id: "thread-author",
      permissions: [],
      allowed: true,
    },
    {
      name: "post manager",
      id: "manager",
      permissions: [Permission.MANAGE_POSTS],
      allowed: true,
    },
    {
      name: "reply author",
      id: "reply-author",
      permissions: [],
      allowed: false,
    },
    { name: "other member", id: "other", permissions: [], allowed: false },
    { name: "guest", id: undefined, permissions: [], allowed: false },
  ])(
    "allows answer selection by $name: $allowed",
    async ({ id, permissions, allowed }) => {
      sessionMock = id ? { id, roles: [{ permissions }] } : undefined;
      const reply = {
        id: "reply-1",
        author: { id: "reply-author" },
        visibility: "published",
      };
      const { result } = renderHook(() =>
        useReplyMenu({
          thread: { slug: "thread-1", author: { id: "thread-author" } } as any,
          reply: reply as any,
          onEdit,
        }),
      );

      expect(result.current.canUpdateConclusion).toBe(allowed);
      await act(async () => {
        await result.current.handlers.handleConclusionChange();
      });
      expect(updateConclusionReply).toHaveBeenCalledTimes(allowed ? 1 : 0);
      if (allowed) {
        expect(updateConclusionReply).toHaveBeenCalledWith(reply);
      }
    },
  );

  it.each([
    { deletedAt: "2026-11-08T10:00:00Z", visibility: "published" },
    { visibility: "review" },
    { visibility: "draft" },
  ])(
    "does not nominate unavailable replies: $visibility $deletedAt",
    (state) => {
      const { result } = renderHook(() =>
        useReplyMenu({
          thread: { slug: "thread-1" } as any,
          reply: { id: "reply-1", ...state } as any,
          onEdit,
        }),
      );
      expect(result.current.canUpdateConclusion).toBe(false);
    },
  );

  it("removes the selected answer without reopening the thread", async () => {
    const { result } = renderHook(() =>
      useReplyMenu({
        thread: { slug: "thread-1", conclusion_reply_id: "reply-1" } as any,
        reply: { id: "reply-1", visibility: "published" } as any,
        onEdit,
      }),
    );
    expect(result.current.isConclusion).toBe(true);
    await act(async () => {
      await result.current.handlers.handleConclusionChange();
    });
    expect(updateConclusionReply).toHaveBeenCalledWith(null);
  });
});
