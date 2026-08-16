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

  it("keeps a top-level user message carrying a stale delegation scope", () => {
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
        isolation_scope: "call_delegate_1",
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
        id: "follow-up",
        role: "user",
        isolation_scope: "call_delegate_1",
        parts: [
          {
            type: "text",
            text: "Nice!",
            state: "done",
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const projected = projectRobotMessages(messages);

    expect(projected.map((message) => message.id)).toEqual([
      "coordinator-call",
      "follow-up",
    ]);
    expect(projected[0]?.parts[0]).toMatchObject({
      type: "data-delegation",
      data: {
        messages: [{ id: "specialist-result" }],
      },
    });
    expect(projected[1]).toMatchObject({
      role: "user",
      parts: [{ type: "text", text: "Nice!" }],
    });
  });

  it("does not infer a delegation from an isolation scope alone", () => {
    const messages = [
      {
        id: "follow-up",
        role: "user",
        isolation_scope: "call_completed_delegation",
        parts: [
          {
            type: "text",
            text: "Okay",
            state: "done",
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    expect(projectRobotMessages(messages)).toEqual(messages);
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

  it("keeps an asynchronous delegation running until its completed result arrives", () => {
    const baseMessages = [
      {
        id: "coordinator-call",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_async_delegate",
            state: "input-available",
            input: { request: "Inspect the thread." },
          },
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_async_delegate",
            state: "output-available",
            output: { status: "pending" },
          },
        ],
      },
      {
        id: "specialist-progress",
        role: "assistant",
        branch: "robot_d9n7eodo2dtgv230rpng",
        isolation_scope: "call_async_delegate",
        robot: {
          id: "d9n7eodo2dtgv230rpng",
          name: "Thread Researcher",
        },
        parts: [{ type: "text", text: "Inspecting now.", state: "done" }],
      },
    ] as unknown as StorydenUIMessage[];

    const pending = projectRobotMessages(baseMessages);
    expect(pending[0]?.parts[0]).toMatchObject({
      type: "data-delegation",
      data: { status: "running" },
    });

    const completed = projectRobotMessages([
      ...baseMessages,
      {
        id: "coordinator-result",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_async_delegate",
            state: "output-available",
            output: {
              status: "completed",
              summary: "The thread was inspected.",
            },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[]);
    expect(completed[0]?.parts[0]).toMatchObject({
      type: "data-delegation",
      data: { status: "completed" },
    });
  });

  it("uses the specialist finish result after hydration without rendering the runtime tool", () => {
    const messages = [
      {
        id: "coordinator-call",
        role: "assistant",
        parts: [
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_async_delegate",
            state: "input-available",
            input: { request: "Inspect the thread." },
          },
          {
            type: "tool-robot_d9n7eodo2dtgv230rpng",
            toolCallId: "call_async_delegate",
            state: "output-available",
            output: { status: "pending" },
          },
        ],
      },
      {
        id: "specialist-result",
        role: "assistant",
        branch: "robot_d9n7eodo2dtgv230rpng",
        isolation_scope: "call_async_delegate",
        robot: {
          id: "d9n7eodo2dtgv230rpng",
          name: "Thread Researcher",
        },
        parts: [
          { type: "text", text: "The thread was inspected.", state: "done" },
          {
            type: "tool-robot_run_finish",
            toolCallId: "call_specialist_finish",
            state: "input-available",
            input: {
              status: "completed",
              summary: "The thread was inspected.",
            },
          },
        ],
      },
      {
        id: "specialist-finish-output",
        role: "assistant",
        branch: "robot_d9n7eodo2dtgv230rpng",
        isolation_scope: "call_async_delegate",
        robot: {
          id: "d9n7eodo2dtgv230rpng",
          name: "Thread Researcher",
        },
        parts: [
          {
            type: "tool-robot_run_finish",
            toolCallId: "call_specialist_finish",
            state: "output-available",
            output: {
              status: "completed",
              summary: "The thread was inspected.",
            },
          },
        ],
      },
    ] as unknown as StorydenUIMessage[];

    const projected = projectRobotMessages(messages);

    expect(projected[0]?.parts[0]).toMatchObject({
      type: "data-delegation",
      data: {
        status: "completed",
        messages: [
          {
            id: "specialist-result",
            parts: [{ type: "text", text: "The thread was inspected." }],
          },
        ],
      },
    });
    expect(JSON.stringify(projected)).not.toContain("robot_run_finish");
  });
});
