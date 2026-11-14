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

const sessionQuery = vi.hoisted(() => ({
  account: undefined as Account | undefined,
}));

let cachedThread: ThreadGetResponse | undefined;
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
  useSession: (initialSession?: Account) =>
    sessionQuery.account ?? initialSession,
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

function threadResponse(
  reactions: Reaction[] = [],
  postID = "reply",
): ThreadGetResponse {
  const thread = {
    id: "thread",
    slug: "thread",
    reacts: [] as Reaction[],
    replies: {
      replies: [{ id: "reply", reacts: [] as Reaction[] } as Reply],
    },
  } as ThreadGetResponse;

  if (postID === "thread") {
    thread.reacts = reactions;
  } else {
    thread.replies.replies[0]!.reacts = reactions;
  }

  return thread;
}

function cachedReactions(postID: string) {
  if (postID === "thread") {
    return cachedThread!.reacts;
  }

  return cachedThread!.replies.replies[0]!.reacts;
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

describe.each(["thread", "reply"])("%s reaction mutations", (postID) => {
  beforeEach(() => {
    sessionQuery.account = { id: "current" } as Account;
    api.postReactAdd.mockReset();
    api.postReactRemove.mockReset();
    mutate.mockClear();
    cachedThread = threadResponse([], postID);
    serverThread = cachedThread;
  });

  test("uses the rendered thread before the client cache has data", async () => {
    const initialThread = threadResponse(
      [
        {
          ...reaction("other-reaction"),
          author: { id: "other" },
        } as Reaction,
      ],
      postID,
    );
    const created = reaction("server-reaction");
    const request = deferred<Reaction>();
    api.postReactAdd.mockReturnValue(request.promise);
    cachedThread = undefined;
    const { result } = renderHook(() =>
      useThreadMutations(initialThread as Thread),
    );

    let operation!: Promise<Reaction | undefined>;
    act(() => {
      operation = result.current.reactionAdd(postID, "🔥");
    });

    await waitFor(() => {
      expect(cachedReactions(postID)).toHaveLength(2);
    });
    expect(cachedReactions(postID)[1]?.author.id).toBe("current");

    await act(async () => {
      request.resolve(created);
      await operation;
    });

    expect(cachedReactions(postID)).toEqual([
      { ...reaction("other-reaction"), author: { id: "other" } },
      created,
    ]);
  });

  test("adds a reaction before the client session fetch completes", async () => {
    sessionQuery.account = undefined;
    const account = { id: "current" } as Account;
    const created = reaction("server-reaction");
    api.postReactAdd.mockResolvedValue(created);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread, undefined, undefined, account),
    );

    await act(async () => {
      await result.current.reactionAdd(postID, "🔥");
    });

    expect(api.postReactAdd).toHaveBeenCalledWith(postID, { emoji: "🔥" });
    expect(cachedReactions(postID)).toEqual([created]);
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
      operation = result.current.reactionAdd(postID, "🔥");
    });

    await waitFor(() => {
      expect(cachedReactions(postID)[0]!.id).toMatch(
        /^optimistic_reply_update_/,
      );
    });
    expect(api.postReactAdd).toHaveBeenCalledOnce();

    await act(async () => {
      request.resolve(created);
      await operation;
    });

    expect(cachedReactions(postID)).toEqual([created]);
  });

  test("rolls back an optimistic add when the delayed request fails", async () => {
    const request = deferred<Reaction>();
    api.postReactAdd.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    let operation!: Promise<Reaction | undefined>;
    act(() => {
      operation = result.current.reactionAdd(postID, "🔥");
    });
    await waitFor(() => {
      expect(cachedReactions(postID)).toHaveLength(1);
    });

    await expect(
      act(async () => {
        request.reject(new Error("network timeout"));
        await operation;
      }),
    ).rejects.toThrow("network timeout");
    expect(cachedReactions(postID)).toEqual([]);
  });

  test("restores server state when a delayed removal fails", async () => {
    const existing = reaction("server-reaction");
    const concurrent = {
      ...reaction("optimistic_reply_update_2"),
      emoji: "👍",
    };
    cachedThread = threadResponse([existing], postID);
    serverThread = cachedThread;
    const request = deferred<void>();
    api.postReactRemove.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    let operation!: Promise<void>;
    act(() => {
      operation = result.current.reactionRemove(postID, existing.id);
    });
    await waitFor(() => {
      expect(cachedReactions(postID)).toEqual([]);
    });

    cachedThread = threadResponse([concurrent], postID);

    await expect(
      act(async () => {
        request.reject(new Error("network timeout"));
        await operation;
      }),
    ).rejects.toThrow("network timeout");
    expect(cachedReactions(postID)).toEqual([existing, concurrent]);
  });

  test("does not send optimistic IDs to the remove endpoint", async () => {
    const optimistic = reaction("optimistic_reply_update_1");
    cachedThread = threadResponse([optimistic], postID);
    const { result } = renderHook(() =>
      useThreadMutations(cachedThread as Thread),
    );

    await act(async () => {
      await result.current.reactionRemove(postID, optimistic.id);
    });

    expect(api.postReactRemove).not.toHaveBeenCalled();
    expect(cachedReactions(postID)).toEqual([optimistic]);
  });
});
