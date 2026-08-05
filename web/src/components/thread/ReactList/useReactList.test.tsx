import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, test, vi } from "vitest";

import type {
  Account,
  React as Reaction,
  Reply,
  Thread,
} from "@/api/openapi-schema";

import { useReactionList } from "./useReactList";

const mutations = vi.hoisted(() => ({
  reactionAdd: vi.fn(),
  reactionRemove: vi.fn(),
}));

vi.mock("@/api/client", () => ({
  handle: async (operation: () => Promise<unknown>) => {
    try {
      return await operation();
    } catch {
      return undefined;
    }
  },
}));

vi.mock("@/auth", () => ({
  useSession: (initialSession?: Account) => initialSession,
}));

vi.mock("@/lib/thread/mutation", () => ({
  useThreadMutations: () => mutations,
}));

const session = { id: "current" } as Account;
const thread = { id: "thread", slug: "thread" } as Thread;

function reaction(id: string, emoji: string, accountID: string): Reaction {
  return {
    id,
    emoji,
    author: { id: accountID },
  } as Reaction;
}

function reply(reactions: Reaction[] = []): Reply {
  return {
    id: "reply",
    reacts: reactions,
  } as Reply;
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

describe("useReactionList", () => {
  beforeEach(() => {
    mutations.reactionAdd.mockReset();
    mutations.reactionRemove.mockReset();
  });

  test("deduplicates picker selections while a slow add request is pending", async () => {
    const request = deferred<void>();
    mutations.reactionAdd.mockReturnValue(request.promise);
    const { result } = renderHook(() =>
      useReactionList({ initialSession: session, thread, reply: reply() }),
    );

    let first!: Promise<void>;
    let second!: Promise<void>;
    act(() => {
      first = result.current.handlers.handleReactPicker("🔥");
      second = result.current.handlers.handleReactPicker("🔥");
    });

    expect(mutations.reactionAdd).toHaveBeenCalledOnce();
    expect(result.current.data.pendingEmojis.has("🔥")).toBe(true);

    await act(async () => {
      request.resolve();
      await Promise.all([first, second]);
    });

    expect(result.current.data.pendingEmojis.has("🔥")).toBe(false);
  });

  test("deduplicates toggles while a slow remove request is pending", async () => {
    const request = deferred<void>();
    mutations.reactionRemove.mockReturnValue(request.promise);
    const currentReaction = reaction("reaction", "🔥", session.id);
    const { result } = renderHook(() =>
      useReactionList({
        initialSession: session,
        thread,
        reply: reply([currentReaction]),
      }),
    );

    let first!: Promise<void>;
    let second!: Promise<void>;
    act(() => {
      first = result.current.handlers.handleReactExisting("🔥");
      second = result.current.handlers.handleReactExisting("🔥");
    });

    expect(mutations.reactionRemove).toHaveBeenCalledOnce();
    expect(mutations.reactionRemove).toHaveBeenCalledWith("reply", "reaction");

    await act(async () => {
      request.resolve();
      await Promise.all([first, second]);
    });
  });

  test("never sends an optimistic reaction ID to removal", async () => {
    const optimistic = reaction("optimistic_reply_update_1", "🔥", session.id);
    const { result } = renderHook(() =>
      useReactionList({
        initialSession: session,
        thread,
        reply: reply([optimistic]),
      }),
    );

    await act(async () => {
      await result.current.handlers.handleReactExisting("🔥");
    });

    expect(mutations.reactionRemove).not.toHaveBeenCalled();
  });

  test("clears pending state after failure so the action can be retried", async () => {
    mutations.reactionAdd
      .mockRejectedValueOnce(new Error("network timeout"))
      .mockResolvedValueOnce(undefined);
    const { result } = renderHook(() =>
      useReactionList({ initialSession: session, thread, reply: reply() }),
    );

    await act(async () => {
      await result.current.handlers.handleReactPicker("🔥");
    });
    expect(result.current.data.pendingEmojis.has("🔥")).toBe(false);

    await act(async () => {
      await result.current.handlers.handleReactPicker("🔥");
    });
    expect(mutations.reactionAdd).toHaveBeenCalledTimes(2);
  });

  test("does not add a reaction already owned by the current account", async () => {
    const { result } = renderHook(() =>
      useReactionList({
        initialSession: session,
        thread,
        reply: reply([reaction("reaction", "🔥", session.id)]),
      }),
    );

    await act(async () => {
      await result.current.handlers.handleReactPicker("🔥");
    });

    expect(mutations.reactionAdd).not.toHaveBeenCalled();
  });
});
