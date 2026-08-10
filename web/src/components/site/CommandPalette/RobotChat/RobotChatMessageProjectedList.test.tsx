import { describe, expect, it } from "vitest";

import { StorydenUIMessage } from "@/api/robots-types";

import {
  projectRobotMessages,
  projectToolOutputs,
} from "./RobotChatMessageProjection";

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

describe("projectRobotMessages", () => {
  it("groups a delegated Robot branch under its parent call", () => {
    const messages = [
      {
        id: "coordinator-call",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_delegate_1",
            toolName: "robot_d9n7eodo2dtgv230rpng",
            state: "input-available",
            input: { request: "Inspect the submitted library page." },
          },
        ],
      },
      {
        id: "specialist-result",
        role: "assistant",
        robot: {
          id: "d9n7eodo2dtgv230rpng",
          name: "Library Curator",
        },
        branch: "robot_d9n7eodo2dtgv230rpng@call_delegate_1",
        parts: [
          {
            type: "text",
            text: "The page has substantive library coverage.",
            state: "done",
          },
        ],
      },
      {
        id: "coordinator-result",
        role: "user",
        parts: [
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_delegate_1",
            toolName: "robot_d9n7eodo2dtgv230rpng",
            state: "output-available",
            input: {},
            output: { result: "The page has substantive library coverage." },
          },
        ],
      },
      {
        id: "coordinator-summary",
        role: "assistant",
        parts: [
          {
            type: "text",
            text: "The curator approved the page.",
            state: "done",
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const projected = projectRobotMessages(messages);

    expect(projected).toHaveLength(2);
    expect(projected[0]?.parts).toEqual([
      {
        type: "data-delegation",
        id: "call_delegate_1",
        data: {
          callId: "call_delegate_1",
          robot: {
            id: "d9n7eodo2dtgv230rpng",
            name: "Library Curator",
          },
          request: "Inspect the submitted library page.",
          status: "completed",
          messages: [messages[1]],
          error: undefined,
        },
      },
    ]);
    expect(projected[1]?.parts).toEqual([
      {
        type: "text",
        text: "The curator approved the page.",
        state: "done",
      },
    ]);
    expect(JSON.stringify(projected)).not.toContain(
      "tool-robot_d9n7eodo2dtgv230rpng",
    );
  });

  it("merges a delegation resumed across chat requests", () => {
    const messages = [
      {
        id: "first-response",
        role: "assistant",
        parts: [
          {
            type: "data-delegation",
            id: "call_delegate_2",
            data: {
              callId: "call_delegate_2",
              robot: {
                id: "d9n7eodo2dtgv230rpng",
                name: "Library Curator",
              },
              request: "Review the submitted page.",
              status: "running",
              messages: [
                {
                  id: "specialist-search",
                  role: "assistant",
                  parts: [
                    {
                      type: "text",
                      text: "I’ll inspect the page.",
                      state: "done",
                    },
                  ],
                },
              ],
            },
          },
        ],
      },
      {
        id: "resumed-response",
        role: "assistant",
        parts: [
          {
            type: "data-delegation",
            id: "call_delegate_2",
            data: {
              callId: "call_delegate_2",
              robot: {
                id: "d9n7eodo2dtgv230rpng",
                name: "Delegated Robot",
              },
              request: "",
              status: "completed",
              messages: [
                {
                  id: "specialist-result",
                  role: "assistant",
                  parts: [
                    {
                      type: "text",
                      text: "The page meets the curation criteria.",
                      state: "done",
                    },
                  ],
                },
              ],
            },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const projected = projectRobotMessages(messages);

    expect(projected).toHaveLength(1);
    expect(projected[0]?.parts[0]).toMatchObject({
      type: "data-delegation",
      id: "call_delegate_2",
      data: {
        callId: "call_delegate_2",
        robot: {
          id: "d9n7eodo2dtgv230rpng",
          name: "Library Curator",
        },
        request: "Review the submitted page.",
        status: "completed",
        messages: [{ id: "specialist-search" }, { id: "specialist-result" }],
      },
    });
  });
});
