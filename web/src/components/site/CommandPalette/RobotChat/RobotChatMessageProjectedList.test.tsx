import { describe, expect, it } from "vitest";

import { StorydenUIMessage } from "@/api/robots-types";

import { projectToolOutputs } from "./RobotChatMessageProjectedList";

describe("projectToolOutputs", () => {
  it("collapses hydrated tool input and output while preserving input arguments", () => {
    const messages = [
      {
        id: "tool-call-message",
        role: "assistant",
        parts: [
          {
            type: "tool-content_search",
            toolCallId: "call_HJMwtXHkKUPOMxATv9n0VkG4",
            toolName: "content_search",
            state: "input-available",
            input: {
              kind: ["thread", "reply"],
              max_results: 10,
              query: "robots agents",
            },
          },
        ],
      },
      {
        id: "tool-output-message",
        role: "user",
        parts: [
          {
            type: "tool-content_search",
            toolCallId: "call_HJMwtXHkKUPOMxATv9n0VkG4",
            toolName: "content_search",
            state: "output-available",
            input: {
              items: [],
              results: 0,
            },
            output: {
              items: [],
              results: 0,
            },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const projected = projectToolOutputs(messages);
    const [part] = projected[0]?.parts ?? [];

    expect(part).toMatchObject({
      type: "tool-content_search",
      toolCallId: "call_HJMwtXHkKUPOMxATv9n0VkG4",
      state: "output-available",
      input: {
        kind: ["thread", "reply"],
        max_results: 10,
        query: "robots agents",
      },
      output: {
        items: [],
        results: 0,
      },
    });
    expect(projected[1]?.parts).toEqual([]);
  });
});
