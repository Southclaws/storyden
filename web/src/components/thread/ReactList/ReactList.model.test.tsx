import { describe, expect, test } from "vitest";

import type { React as Reaction } from "@/api/openapi-schema";

import { groupReactions } from "./ReactList.model";

function reaction(id: string, emoji: string, accountID?: string): Reaction {
  return {
    id,
    emoji,
    author: accountID ? { id: accountID } : undefined,
  } as Reaction;
}

describe("groupReactions", () => {
  test("groups reactions in first-seen order and marks the current account", () => {
    const reactions = [
      reaction("1", "🔥", "other"),
      reaction("2", "👍", "current"),
      reaction("3", "🔥", "current"),
    ];

    expect(groupReactions("current", reactions)).toEqual([
      {
        emoji: "🔥",
        count: 2,
        hasReacted: true,
        reactedBy: [],
        reactions: [reactions[0], reactions[2]],
      },
      {
        emoji: "👍",
        count: 1,
        hasReacted: true,
        reactedBy: [],
        reactions: [reactions[1]],
      },
    ]);
  });

  test("does not mark anonymous reactions as belonging to a guest", () => {
    expect(groupReactions(undefined, [reaction("1", "🔥")])).toEqual([
      expect.objectContaining({ hasReacted: false }),
    ]);
  });
});
