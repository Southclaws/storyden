import { stream } from "@durable-streams/client";
import type { ChatTransport, UIMessage } from "ai";

import type { RobotSessionStreamEvent } from "./openapi-schema/robotSessionStreamEvent";

export type CommandAccepted = {
  sessionId: string;
  messageId: string;
  clientMessageId?: string;
  clientMessageRole?: string;
};

type DurableChatTransportOptions = {
  api: string;
  headers?: HeadersInit;
  fetchClient?: typeof fetch;
  onCommandAccepted?: (command: CommandAccepted) => void;
};

type ObserveSessionOptions = {
  url: string;
  offset: string;
  signal?: AbortSignal;
  fetchClient?: typeof fetch;
};

function mergeHeaders(headers?: HeadersInit) {
  if (!headers) return {};
  return Object.fromEntries(new Headers(headers).entries());
}

function emptyMessageStream() {
  return new ReadableStream<never>({
    start(controller) {
      controller.close();
    },
  });
}

function parseCommandAccepted(body: unknown): CommandAccepted | undefined {
  if (
    !body ||
    typeof body !== "object" ||
    !("sessionId" in body) ||
    !("messageId" in body) ||
    typeof body.sessionId !== "string" ||
    typeof body.messageId !== "string"
  ) {
    return;
  }
  return { sessionId: body.sessionId, messageId: body.messageId };
}

export async function observeRobotSession({
  url,
  offset,
  signal,
  fetchClient = fetch,
}: ObserveSessionOptions) {
  const response = await stream<RobotSessionStreamEvent>({
    url,
    offset,
    live: "sse",
    json: true,
    signal,
    fetch: fetchClient,
  });
  return response.jsonStream();
}

export function createDurableChatTransport<UI_MESSAGE extends UIMessage>({
  api,
  headers,
  fetchClient = fetch,
  onCommandAccepted,
}: DurableChatTransportOptions): ChatTransport<UI_MESSAGE> {
  return {
    async sendMessages({
      trigger,
      chatId,
      messageId,
      messages,
      abortSignal,
      body,
      headers: requestHeaders,
    }) {
      const response = await fetchClient(api, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...mergeHeaders(headers),
          ...mergeHeaders(requestHeaders),
        },
        body: JSON.stringify({
          ...(body ?? {}),
          id: chatId,
          messages,
          trigger,
          messageId,
        }),
        signal: abortSignal,
      });
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `HTTP error ${response.status}`);
      }

      const accepted = parseCommandAccepted(await response.json());
      if (!accepted) {
        throw new Error(
          "Robot message identity missing from command response.",
        );
      }
      const clientMessage = messages[messages.length - 1];
      onCommandAccepted?.({
        ...accepted,
        clientMessageId: clientMessage?.id,
        clientMessageRole: clientMessage?.role,
      });
      return emptyMessageStream();
    },

    async reconnectToStream() {
      return null;
    },
  };
}
