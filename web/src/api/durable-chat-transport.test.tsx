import { describe, expect, it, vi } from "vitest";
import type { UIMessageChunk } from "ai";

import { createDurableChatTransport } from "./durable-chat-transport";

describe("createDurableChatTransport", () => {
  it("uses the configured fetch client for the durable stream read", async () => {
    const fetchClient = vi.fn<typeof fetch>(async (input, init) => {
      const url = input.toString();
      if (init?.method === "POST") {
        return new Response(
          JSON.stringify({ streamUrl: "http://api.test/api/robots/sessions/session-1/turns/turn-1" }),
          {
            status: 201,
            headers: {
              "Content-Type": "application/json",
              Location: "http://api.test/api/robots/sessions/session-1/turns/turn-1",
            },
          },
        );
      }

      expect(url).toContain("/api/robots/sessions/session-1/turns/turn-1?offset=-1");
      return new Response(
        JSON.stringify([
          { type: "start", messageId: "response-1" },
          { type: "finish" },
        ]),
        {
          headers: {
            "Content-Type": "application/json",
            "Stream-Next-Offset": "0000000000000000_0000000000000002",
            "Stream-Up-To-Date": "true",
            "Stream-Closed": "true",
          },
        },
      );
    });

    const transport = createDurableChatTransport({
      api: "http://api.test/api/robots/sessions",
      fetchClient,
    });
    const response = await transport.sendMessages({
      trigger: "submit-message",
      chatId: "session-1",
      messageId: undefined,
      messages: [],
      abortSignal: undefined,
    });

    const chunks: UIMessageChunk[] = [];
    for await (const chunk of response) {
      chunks.push(chunk);
    }

    expect(chunks).toEqual([
      { type: "start", messageId: "response-1" },
      { type: "finish" },
    ]);
    expect(fetchClient).toHaveBeenCalledTimes(2);
  });
});
