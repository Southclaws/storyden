import { describe, expect, it } from "vitest";

import { StorydenUIMessage } from "@/api/robots-types";

import {
  assistantCompletedToolOutputIDs,
  assistantHasTextAfterToolOutput,
  assistantToolOutputsAreComplete,
  hasUnhydratedToolOutput,
  reconcileMessages,
  shouldReplaceMessages,
} from "./RobotChatContext";

describe("assistantHasTextAfterToolOutput", () => {
  it("allows assistant text before a pending confirmation tool", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        { type: "text", text: "Deleting this now." },
        {
          type: "tool-robot_delete",
          toolCallId: "call_1",
          toolName: "robot_delete",
          state: "output-available",
          input: { id: "robot_1" },
          output: { _storyden_confirmation: { approved: true } },
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantHasTextAfterToolOutput(message)).toBe(false);
  });

  it("allows assistant text before multiple completed confirmation tools", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        { type: "text", text: "Deleting these now." },
        {
          type: "tool-robot_delete",
          toolCallId: "call_1",
          toolName: "robot_delete",
          state: "output-available",
          input: { id: "robot_1" },
          output: { _storyden_confirmation: { approved: true } },
        },
        {
          type: "tool-robot_delete",
          toolCallId: "call_2",
          toolName: "robot_delete",
          state: "output-available",
          input: { id: "robot_2" },
          output: { _storyden_confirmation: { approved: true } },
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantHasTextAfterToolOutput(message)).toBe(false);
  });

  it("blocks auto-send after the assistant has responded to a tool output", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        {
          type: "tool-robot_delete",
          toolCallId: "call_1",
          toolName: "robot_delete",
          state: "output-available",
          input: { id: "robot_1" },
          output: { _storyden_confirmation: { approved: true } },
        },
        { type: "text", text: "Deleted." },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantHasTextAfterToolOutput(message)).toBe(true);
  });
});

describe("assistantToolOutputsAreComplete", () => {
  it("does not auto-submit while one confirmation in the assistant turn is still pending", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        {
          type: "tool-robot_delete",
          toolCallId: "call_1",
          state: "output-available",
          input: { id: "robot_1" },
          output: { _storyden_confirmation: { approved: true } },
        },
        {
          type: "tool-robot_delete",
          toolCallId: "call_2",
          state: "input-available",
          input: { id: "robot_2" },
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantToolOutputsAreComplete(message)).toBe(false);
  });

  it("auto-submits once all confirmations in the assistant turn are resolved", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        {
          type: "tool-robot_delete",
          toolCallId: "call_1",
          state: "output-available",
          input: { id: "robot_1" },
          output: { _storyden_confirmation: { approved: true } },
        },
        {
          type: "tool-robot_delete",
          toolCallId: "call_2",
          state: "output-available",
          input: { id: "robot_2" },
          output: { _storyden_confirmation: { approved: true } },
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantToolOutputsAreComplete(message)).toBe(true);
  });
});

describe("assistantCompletedToolOutputIDs", () => {
  it("returns completed tool IDs from the current assistant step", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        {
          type: "tool-robot_switch",
          toolCallId: "call_switch",
          toolName: "robot_switch",
          state: "output-available",
          input: { robot_id: "plugin_builder" },
          output: { success: true, robot_id: "plugin_builder" },
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantCompletedToolOutputIDs(message)).toEqual(["call_switch"]);
  });

  it("ignores provider-executed tools", () => {
    const message = {
      id: "message-1",
      role: "assistant",
      parts: [
        {
          type: "tool-plugin_file_read",
          toolCallId: "call_read",
          state: "output-available",
          providerExecuted: true,
        },
      ],
    } as unknown as StorydenUIMessage;

    expect(assistantCompletedToolOutputIDs(message)).toEqual([]);
  });
});

describe("hasUnhydratedToolOutput", () => {
  it("detects local confirmation output that incoming history has not persisted yet", () => {
    const localMessages = [
      {
        id: "message-1",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_delete",
            toolCallId: "call_1",
            toolName: "robot_delete",
            state: "output-available",
            input: { id: "robot_1" },
            output: { _storyden_confirmation: { approved: true } },
          },
          {
            type: "tool-robot_delete",
            toolCallId: "call_2",
            toolName: "robot_delete",
            state: "input-available",
            input: { id: "robot_2" },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const incomingMessages = [
      {
        id: "message-1",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_delete",
            toolCallId: "call_1",
            toolName: "robot_delete",
            state: "input-available",
            input: { id: "robot_1" },
          },
          {
            type: "tool-robot_delete",
            toolCallId: "call_2",
            toolName: "robot_delete",
            state: "input-available",
            input: { id: "robot_2" },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    expect(hasUnhydratedToolOutput(localMessages, incomingMessages)).toBe(true);
  });

  it("allows incoming history once matching tool output is present", () => {
    const messages = [
      {
        id: "message-1",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_delete",
            toolCallId: "call_1",
            toolName: "robot_delete",
            state: "output-available",
            input: { id: "robot_1" },
            output: { _storyden_confirmation: { approved: true } },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    expect(hasUnhydratedToolOutput(messages, messages)).toBe(false);
  });
});

describe("shouldReplaceMessages", () => {
  it("replaces stale local windows when incoming contains newer message ids", () => {
    const localMessages = [
      { id: "old-1", role: "user", parts: [{ type: "text", text: "old" }] },
      {
        id: "old-2",
        role: "assistant",
        parts: [{ type: "text", text: "older" }],
      },
    ] as unknown as StorydenUIMessage[];

    const incomingMessages = [
      {
        id: "old-2",
        role: "assistant",
        parts: [{ type: "text", text: "older" }],
      },
      {
        id: "new-1",
        role: "assistant",
        parts: [{ type: "text", text: "new" }],
      },
    ] as unknown as StorydenUIMessage[];

    expect(shouldReplaceMessages(localMessages, incomingMessages)).toBe(true);
  });

  it("keeps local messages when incoming is only an already-present subset", () => {
    const localMessages = [
      { id: "old-1", role: "user", parts: [{ type: "text", text: "old" }] },
      {
        id: "old-2",
        role: "assistant",
        parts: [{ type: "text", text: "older" }],
      },
    ] as unknown as StorydenUIMessage[];

    const incomingMessages = [
      { id: "old-1", role: "user", parts: [{ type: "text", text: "old" }] },
    ] as unknown as StorydenUIMessage[];

    expect(shouldReplaceMessages(localMessages, incomingMessages)).toBe(false);
  });
});

describe("reconcileMessages", () => {
  it("preserves older local messages before the incoming server window", () => {
    const localMessages = [
      { id: "older-1", role: "user", parts: [{ type: "text", text: "older" }] },
      {
        id: "overlap",
        role: "assistant",
        parts: [{ type: "text", text: "old copy" }],
      },
      {
        id: "stale-tail",
        role: "assistant",
        parts: [{ type: "text", text: "stale" }],
      },
    ] as unknown as StorydenUIMessage[];

    const incomingMessages = [
      {
        id: "overlap",
        role: "assistant",
        parts: [{ type: "text", text: "new copy" }],
      },
      {
        id: "new-tail",
        role: "assistant",
        parts: [{ type: "text", text: "new" }],
      },
    ] as unknown as StorydenUIMessage[];

    expect(
      reconcileMessages(localMessages, incomingMessages).map((m) => m.id),
    ).toEqual(["older-1", "overlap", "new-tail"]);
  });
});
