import { act, renderHook, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, test, vi } from "vitest";

import type {
  Account,
  React as Reaction,
  Reply,
  Thread,
  ThreadGetResponse,
} from "@/api/openapi-schema";

import { useThreadMutations } from "./mutation";

const api = vi.hoisted(() => ({
  postReactAdd: vi.fn(),
  postReactRemove: vi.fn(),
}));

let cachedThread: ThreadGetResponse;
let serverThread: ThreadGetResponse;

const mutate = vi.fn(async (_key, updater, _options) => {
  if (typeof updater === "function") {
    cachedThread = updater(cachedThread);
  } else {
    cachedThread = serverThread;
  }
  return cachedThread;
});

vi.mock("swr", () => ({
  useSWRConfig: () => ({ mutate }),
}));

vi.mock("@/auth", () => ({
  useSession: () => ({ id: "current" }) as Account,
}));

vi.mock("@/api/openapi-client/posts", () => ({
  postDelete: vi.fn(),
  postReactAdd: api.postReactAdd,
  postReactRemove: api.postReactRemove,
  postUpdate: vi.fn(),
}));

vi.mock("@/api/openapi-client/replies", () => ({
  replyCreate: vi.fn(),
}));

vi.mock("@/api/openapi-client/threads", () => ({
  getThreadGetKey: () => ["/threads/thread", { page: "1" }],
  threadUpdate: vi.fn(),
}));

function reaction(id: string): Reaction {
  return {
    id,
    emoji: "🔥",
    author: { id: "current" },
  } as Reaction;
}

function threadResponse(reactions: Reaction[] = []): ThreadGetResponse {
  return {
    id: "thread",
    slug: "thread",
    replies: {
      replies: [{ id: "reply", reacts: reactions } as Reply],
    },
  } as ThreadGetResponse;
}

function cachedReactions() {
  return cachedThread.replies.replies[0]!.reacts;
}

function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((resolvePromise, rejectPromise) => {
    resolve = resolvePromise;
    reject = rejectPromise;
  });

  return { promise, reject, resolve };
}

describe("reaction mutations", () => {
  beforeEach(() => {
    api.postReactAdd.mockReset();
    api.postReactRemove.mockReset();
    mutate.mockClear();
    cachedThread = threadResponse();
    serverThread = cachedThread;
  });

  test("keeps an optimistic add during latency then replaces it with the server reaction", async () => {
    const request = deferred<Reaction>();
    const created = reaction("server-reaction");
    api.postReactAdd.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    let operation!: Promise<Reaction | undefined>;
    act(() => {
      operation = result.current.reactionAdd("reply", "🔥");
    });

    await waitFor(() => {
      expect(cachedReactions()[0]!.id).toMatch(/^optimistic_reply_update_/);
    });
    expect(api.postReactAdd).toHaveBeenCalledOnce();

    await act(async () => {
      request.resolve(created);
      await operation;
    });

    expect(cachedReactions()).toEqual([created]);
  });

  test("rolls back an optimistic add when the delayed request fails", async () => {
    const request = deferred<Reaction>();
    api.postReactAdd.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    let operation!: Promise<Reaction | undefined>;
    act(() => {
      operation = result.current.reactionAdd("reply", "🔥");
    });
    await waitFor(() => {
      expect(cachedReactions()).toHaveLength(1);
    });

    await expect(
      act(async () => {
        request.reject(new Error("network timeout"));
        await operation;
      }),
    ).rejects.toThrow("network timeout");
    expect(cachedReactions()).toEqual([]);
  });

  test("restores server state when a delayed removal fails", async () => {
    const existing = reaction("server-reaction");
    const concurrent = {
      ...reaction("optimistic_reply_update_2"),
      emoji: "👍",
    };
    cachedThread = threadResponse([existing]);
    serverThread = cachedThread;
    const request = deferred<void>();
    api.postReactRemove.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    let operation!: Promise<void>;
    act(() => {
      operation = result.current.reactionRemove("reply", existing.id);
    });
    await waitFor(() => {
      expect(cachedReactions()).toEqual([]);
    });

    cachedThread = threadResponse([concurrent]);

    await expect(
      act(async () => {
        request.reject(new Error("network timeout"));
        await operation;
      }),
    ).rejects.toThrow("network timeout");
    expect(cachedReactions()).toEqual([existing, concurrent]);
  });

  test("does not send optimistic IDs to the remove endpoint", async () => {
    const optimistic = reaction("optimistic_reply_update_1");
    cachedThread = threadResponse([optimistic]);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    await act(async () => {
      await result.current.reactionRemove("reply", optimistic.id);
    });

    expect(api.postReactRemove).not.toHaveBeenCalled();
    expect(cachedReactions()).toEqual([optimistic]);
  });
});
