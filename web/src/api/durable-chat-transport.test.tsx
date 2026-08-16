import type { UIMessageChunk } from "ai";
import { describe, expect, it, vi } from "vitest";

import {
  createDurableChatTransport,
  observeRobotSession,
} from "./durable-chat-transport";
import type { RobotSessionStreamEvent } from "./openapi-schema/robotSessionStreamEvent";

describe("createDurableChatTransport", () => {
  it("submits a command without coupling it to the response stream", async () => {
    const onCommandAccepted = vi.fn();
    const fetchClient = vi.fn<typeof fetch>(async () => {
      return new Response(
        JSON.stringify({
          streamUrl: "http://api.test/api/robots/sessions/session-1/stream",
          sessionId: "session-1",
          messageId: "message-1",
        }),
        { status: 202, headers: { "Content-Type": "application/json" } },
      );
    });
    const transport = createDurableChatTransport({
      api: "http://api.test/api/robots/sessions",
      fetchClient,
      onCommandAccepted,
    });

    const response = await transport.sendMessages({
      trigger: "submit-message",
      chatId: "session-1",
      messageId: undefined,
      messages: [{ id: "message-1", role: "user", parts: [] }],
      abortSignal: undefined,
    });

    const chunks: UIMessageChunk[] = [];
    for await (const chunk of response) chunks.push(chunk);
    expect(chunks).toEqual([]);
    expect(fetchClient).toHaveBeenCalledTimes(1);
    expect(onCommandAccepted).toHaveBeenCalledWith({
      sessionId: "session-1",
      messageId: "message-1",
      clientMessageId: "message-1",
      clientMessageRole: "user",
    });
  });
});

describe("observeRobotSession", () => {
  it("starts the session feed from the snapshot offset", async () => {
    const fetchClient = vi.fn<typeof fetch>(async (input) => {
      expect(input.toString()).toContain(
        "/api/robots/sessions/session-1/stream?offset=0000000000000000_0000000000000042",
      );
      return new Response(
        JSON.stringify([
          {
            sequence: 43,
            turn_id: "turn-1",
            event_kind: "turn_completed",
            parts: [],
          },
        ]),
        {
          headers: {
            "Content-Type": "application/json",
            "Stream-Next-Offset": "0000000000000000_0000000000000043",
            "Stream-Up-To-Date": "true",
            "Stream-Closed": "true",
          },
        },
      );
    });

    const events = await observeRobotSession({
      url: "http://api.test/api/robots/sessions/session-1/stream",
      offset: "0000000000000000_0000000000000042",
      fetchClient,
    });

    const received: RobotSessionStreamEvent[] = [];
    for await (const event of events) received.push(event);
    expect(received).toEqual([
      {
        sequence: 43,
        turn_id: "turn-1",
        event_kind: "turn_completed",
        parts: [],
      },
    ]);
  });
});
