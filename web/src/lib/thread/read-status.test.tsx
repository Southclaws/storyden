import { act, renderHook, waitFor } from "@testing-library/react";
import type { PropsWithChildren } from "react";
import useSWR, { type Cache, SWRConfig, unstable_serialize } from "swr";
import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  getThreadGetKey,
  getThreadListKey,
  threadReadStatusUpdate,
} from "@/api/openapi-client/threads";
import {
  DatagraphItemKind,
  type ReadStatus,
  type Thread,
  type ThreadListOKResponse,
} from "@/api/openapi-schema";
import { useBeacon } from "@/lib/beacon/useBeacon";

import { useThreadReadStatus } from "./read-status";

vi.mock("@/auth", () => ({
  useSession: () => undefined,
}));

vi.mock("@/api/openapi-client/threads", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@/api/openapi-client/threads")>()),
  threadReadStatusUpdate: vi.fn(),
}));

const baseline = "2026-11-08T10:00:00Z";
const markedAt = "2026-11-08T11:00:00Z";

const thread = {
  id: "thread-id",
  slug: "thread-slug",
  read_status: { last_read_at: baseline, replies_since: 5 },
  replies: { replies: [] },
} as unknown as Thread;

function withCache(cache: Cache) {
  return function Wrapper({ children }: PropsWithChildren) {
    return <SWRConfig value={{ provider: () => cache }}>{children}</SWRConfig>;
  };
}

describe("thread read status", () => {
  beforeEach(() => {
    Object.defineProperty(navigator, "sendBeacon", {
      configurable: true,
      value: vi.fn(() => true),
    });
  });

  it("updates every cached page and feed without another visit beacon", async () => {
    const firstPageKey = getThreadGetKey(thread.slug, { page: "1" });
    const secondPageKey = getThreadGetKey(thread.slug, { page: "2" });
    const feedKey = getThreadListKey({ page: "1" });
    const categoryKey = getThreadListKey({
      page: "2",
      categories: ["support"],
    });
    const unrelated = { ...thread, id: "other-id", slug: "other-slug" };
    const feed: ThreadListOKResponse = {
      threads: [thread, unrelated],
      current_page: 1,
      page_size: 50,
      results: 2,
      total_pages: 1,
    };
    const cache = new Map();

    for (const key of [firstPageKey, secondPageKey]) {
      cache.set(unstable_serialize(key), { data: thread, _k: key });
    }

    for (const key of [feedKey, categoryKey]) {
      cache.set(unstable_serialize(key), { data: feed, _k: key });
    }

    const { result } = renderHook(
      () => {
        const { data } = useSWR<Thread>(firstPageKey, null);
        const mutations = useThreadReadStatus(data ?? thread);
        useBeacon(DatagraphItemKind.thread, data?.id, mutations.recordRead);

        return { data, ...mutations };
      },
      { wrapper: withCache(cache) },
    );

    await waitFor(() => expect(navigator.sendBeacon).toHaveBeenCalledTimes(1));
    await waitFor(() =>
      expect(result.current.data?.read_status?.replies_since).toBe(0),
    );
    expect(result.current.data?.read_status?.marked_unread_at).toBeUndefined();

    const unread: ReadStatus = {
      last_read_at: baseline,
      replies_since: 5,
      marked_unread_at: markedAt,
    };
    vi.mocked(threadReadStatusUpdate).mockResolvedValueOnce(unread);

    await act(async () => {
      await result.current.updateReadStatus("unread");
    });

    expect(threadReadStatusUpdate).toHaveBeenCalledWith(thread.slug, {
      status: "unread",
    });
    expect(result.current.data?.read_status).toEqual(unread);

    for (const key of [firstPageKey, secondPageKey]) {
      expect(cache.get(unstable_serialize(key)).data.read_status).toEqual(
        unread,
      );
    }

    for (const key of [feedKey, categoryKey]) {
      const cached = cache.get(unstable_serialize(key)).data;
      expect(cached.threads[0].read_status).toEqual(unread);
      expect(cached.threads[1]).toEqual(unrelated);
    }

    const read: ReadStatus = { last_read_at: markedAt, replies_since: 0 };
    vi.mocked(threadReadStatusUpdate).mockResolvedValueOnce(read);

    await act(async () => {
      await result.current.updateReadStatus("read");
    });

    expect(result.current.data?.read_status).toEqual(read);
    expect(navigator.sendBeacon).toHaveBeenCalledTimes(1);

    vi.mocked(threadReadStatusUpdate).mockResolvedValueOnce(unread);

    await act(async () => {
      await result.current.updateReadStatus("unread");
    });

    await act(async () => {
      await result.current.recordRead();
    });

    expect(result.current.data?.read_status?.marked_unread_at).toBeUndefined();
    expect(result.current.data?.read_status?.replies_since).toBe(0);
    expect(navigator.sendBeacon).toHaveBeenCalledTimes(1);
  });

  it("updates a detail screen before SWR has populated its shared cache", async () => {
    const cache = new Map();

    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: thread,
        });
        return { data, ...useThreadReadStatus(data) };
      },
      { wrapper: withCache(cache) },
    );

    const read = { last_read_at: markedAt, replies_since: 0 };
    vi.mocked(threadReadStatusUpdate).mockResolvedValueOnce(read);

    await act(async () => {
      await result.current.updateReadStatus("read");
    });

    expect(result.current.data.read_status).toEqual(read);
  });

  it("preserves cached state when the request fails", async () => {
    const cache = new Map();

    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: thread,
        });
        return { data, ...useThreadReadStatus(data) };
      },
      { wrapper: withCache(cache) },
    );

    vi.mocked(threadReadStatusUpdate).mockRejectedValueOnce(
      new Error("Unavailable"),
    );

    await act(async () => {
      await expect(result.current.updateReadStatus("unread")).rejects.toThrow(
        "Unavailable",
      );
    });

    expect(result.current.data.read_status).toEqual(thread.read_status);
  });

  it("keeps the current read state when the browser rejects a visit beacon", async () => {
    vi.mocked(navigator.sendBeacon).mockReturnValue(false);

    const { result } = renderHook(
      () => {
        const { data } = useSWR(getThreadGetKey(thread.slug), null, {
          fallbackData: thread,
        });
        const mutations = useThreadReadStatus(data);
        useBeacon(DatagraphItemKind.thread, data.id, mutations.recordRead);

        return data;
      },
      { wrapper: withCache(new Map()) },
    );

    await waitFor(() => expect(navigator.sendBeacon).toHaveBeenCalledTimes(1));
    expect(result.current.read_status).toEqual(thread.read_status);
  });
});
