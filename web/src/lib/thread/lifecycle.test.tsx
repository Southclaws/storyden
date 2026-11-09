import { act, renderHook } from "@testing-library/react";
import type { PropsWithChildren } from "react";
import useSWR, { type Cache, SWRConfig, unstable_serialize } from "swr";
import { describe, expect, it, vi } from "vitest";

import { getReplyGetKey } from "@/api/openapi-client/replies";
import {
  getThreadGetKey,
  getThreadListKey,
  threadUpdate,
} from "@/api/openapi-client/threads";
import type { Reply, Thread, ThreadListOKResponse } from "@/api/openapi-schema";

import { useThreadLifecycle } from "./lifecycle";

vi.mock("@/api/openapi-client/threads", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@/api/openapi-client/threads")>()),
  threadUpdate: vi.fn(),
}));

const thread = {
  id: "thread-id",
  slug: "thread-slug",
  lifecycle_state: "resolved",
  conclusion_reply_id: "reply-id",
  updatedAt: "2026-11-08T10:00:00Z",
  replies: { replies: [{ id: "first-page-reply" }] },
} as Thread;

function withCache(cache: Cache) {
  return function Wrapper({ children }: PropsWithChildren) {
    return <SWRConfig value={{ provider: () => cache }}>{children}</SWRConfig>;
  };
}

describe("thread lifecycle", () => {
  it("nominates and removes an answer while synchronising the detail and reply caches", async () => {
    const cache = new Map();
    const initial = {
      ...thread,
      lifecycle_state: "open",
      conclusion_reply_id: undefined,
    } satisfies Thread;
    const answer = {
      id: "new-answer",
      root_id: thread.id,
      body: "<p>An answer</p>",
    } as Reply;
    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: initial,
        });
        return { data, ...useThreadLifecycle(data) };
      },
      { wrapper: withCache(cache) },
    );
    vi.mocked(threadUpdate).mockResolvedValueOnce({
      ...thread,
      conclusion_reply_id: answer.id,
    });

    await act(async () => {
      await result.current.updateConclusionReply(answer);
    });

    expect(threadUpdate).toHaveBeenCalledWith(thread.id, {
      conclusion_reply_id: answer.id,
    });
    expect(result.current.data.lifecycle_state).toBe("resolved");
    expect(result.current.data.conclusion_reply_id).toBe(answer.id);
    expect(
      cache.get(unstable_serialize(getReplyGetKey(thread.slug, answer.id)))
        .data,
    ).toEqual(answer);

    vi.mocked(threadUpdate).mockResolvedValueOnce({
      ...thread,
      conclusion_reply_id: undefined,
    });
    await act(async () => {
      await result.current.updateConclusionReply(null);
    });
    expect(threadUpdate).toHaveBeenLastCalledWith(thread.id, {
      conclusion_reply_id: null,
    });
    expect(result.current.data.conclusion_reply_id).toBeUndefined();
    expect(result.current.data.lifecycle_state).toBe("resolved");
  });

  it("reopens across feeds and every detail page without replacing replies", async () => {
    const firstPageKey = getThreadGetKey(thread.slug, { page: "1" });
    const secondPageKey = getThreadGetKey(thread.id, { page: "2" });
    const secondPage = {
      ...thread,
      replies: { replies: [{ id: "second-page-reply" }] },
    };
    const unrelated = { ...thread, id: "other-id", slug: "other-slug" };
    const otherThreadKey = getThreadGetKey(unrelated.slug);
    const feedKeys = [
      getThreadListKey({ page: "1" }),
      getThreadListKey({ page: "2", categories: ["support"] }),
    ];
    const feed: ThreadListOKResponse = {
      threads: [thread, unrelated],
      current_page: 1,
      page_size: 50,
      results: 2,
      total_pages: 1,
    };
    const cache = new Map();

    cache.set(unstable_serialize(firstPageKey), {
      data: thread,
      _k: firstPageKey,
    });
    cache.set(unstable_serialize(secondPageKey), {
      data: secondPage,
      _k: secondPageKey,
    });
    cache.set(unstable_serialize(otherThreadKey), {
      data: unrelated,
      _k: otherThreadKey,
    });
    for (const key of feedKeys) {
      cache.set(unstable_serialize(key), { data: feed, _k: key });
    }

    const { result } = renderHook(() => useThreadLifecycle(thread), {
      wrapper: withCache(cache),
    });
    const updated = {
      ...thread,
      lifecycle_state: "open",
      conclusion_reply_id: undefined,
      updatedAt: "2026-11-08T11:00:00Z",
    } satisfies Thread;
    vi.mocked(threadUpdate).mockResolvedValueOnce(updated);

    await act(async () => {
      await result.current.updateLifecycleState("open");
    });

    expect(threadUpdate).toHaveBeenCalledWith(thread.id, {
      lifecycle_state: "open",
    });
    for (const key of [firstPageKey, secondPageKey]) {
      const cached = cache.get(unstable_serialize(key)).data;
      expect(cached.lifecycle_state).toBe("open");
      expect(cached.conclusion_reply_id).toBeUndefined();
      expect(cached.updatedAt).toBe(updated.updatedAt);
    }
    expect(cache.get(unstable_serialize(secondPageKey)).data.replies).toEqual(
      secondPage.replies,
    );
    expect(cache.get(unstable_serialize(otherThreadKey)).data).toEqual(
      unrelated,
    );
    for (const key of feedKeys) {
      const cached = cache.get(unstable_serialize(key)).data;
      expect(cached.threads[0].lifecycle_state).toBe("open");
      expect(cached.threads[0].conclusion_reply_id).toBeUndefined();
      expect(cached.threads[1]).toEqual(unrelated);
    }
  });

  it("resolves a detail screen before fallback data enters the shared cache", async () => {
    const initial = {
      ...thread,
      lifecycle_state: "open",
      conclusion_reply_id: undefined,
    } satisfies Thread;
    const cache = new Map();
    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: initial,
        });
        return { data, ...useThreadLifecycle(data) };
      },
      { wrapper: withCache(cache) },
    );
    vi.mocked(threadUpdate).mockResolvedValueOnce({
      ...initial,
      lifecycle_state: "resolved",
    });

    await act(async () => {
      await result.current.updateLifecycleState("resolved");
    });

    expect(result.current.data.lifecycle_state).toBe("resolved");
    expect(result.current.data.replies).toEqual(initial.replies);
  });

  it("keeps the existing state and conclusion when saving fails", async () => {
    const cache = new Map();
    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: thread,
        });
        return { data, ...useThreadLifecycle(data) };
      },
      { wrapper: withCache(cache) },
    );
    vi.mocked(threadUpdate).mockRejectedValueOnce(new Error("Unable to save"));

    await act(async () => {
      await expect(result.current.updateLifecycleState("open")).rejects.toThrow(
        "Unable to save",
      );
    });

    expect(result.current.data).toEqual(thread);
  });
});
