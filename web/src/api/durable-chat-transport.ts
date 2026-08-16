import { stream } from "@durable-streams/client";
import type {
  ChatTransport,
  UIMessage,
  UIMessageChunk,
} from "ai";

type DurableChatTransportOptions = {
  api: string;
  reconnectApi?: string;
  headers?: HeadersInit;
  fetchClient?: typeof fetch;
};

function mergeHeaders(headers?: HeadersInit) {
  if (!headers) return {};
  return Object.fromEntries(new Headers(headers).entries());
}

function parseBodyStreamURL(body: unknown) {
  if (
    body &&
    typeof body === "object" &&
    "streamUrl" in body &&
    typeof body.streamUrl === "string" &&
    body.streamUrl.length > 0
  ) {
    return body.streamUrl;
  }
  return undefined;
}

async function parseJSON(response: Response) {
  if (!response.headers.get("content-type")?.includes("application/json")) {
    return;
  }
  try {
    return (await response.json()) as unknown;
  } catch {
    return;
  }
}

function resolveStreamURL(streamURL: string, responseURL: string, requestURL: string) {
  if (/^[a-zA-Z][a-zA-Z\d+\-.]*:/.test(streamURL)) {
    return streamURL;
  }

  for (const baseURL of [responseURL, requestURL, window.location.href]) {
    if (!baseURL) continue;
    try {
      return new URL(streamURL, baseURL).toString();
    } catch {
      continue;
    }
  }
  throw new Error(`Failed to resolve durable stream URL "${streamURL}".`);
}

function toReadableStream<T>(iterable: AsyncIterable<T>) {
  const iterator = iterable[Symbol.asyncIterator]();
  return new ReadableStream<T>({
    async pull(controller) {
      const result = await iterator.next();
      if (result.done) {
        controller.close();
        return;
      }
      controller.enqueue(result.value);
    },
    async cancel() {
      await iterator.return?.();
    },
  });
}

async function readUIMessageChunks(
  streamURL: string,
  abortSignal: AbortSignal | undefined,
  fetchClient: typeof fetch,
) {
  const response = await stream<UIMessageChunk>({
    url: streamURL,
    live: "sse",
    json: true,
    signal: abortSignal,
    fetch: fetchClient,
  });
  return toReadableStream(response.jsonStream());
}

async function streamURLFromResponse(response: Response) {
  return (
    response.headers.get("Location") ??
    parseBodyStreamURL(await parseJSON(response))
  );
}

export function createDurableChatTransport<UI_MESSAGE extends UIMessage>({
  api,
  reconnectApi,
  headers,
  fetchClient = fetch,
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

      const streamURL = await streamURLFromResponse(response);
      if (!streamURL) {
        throw new Error("Durable stream URL missing from chat response.");
      }
      return readUIMessageChunks(
        resolveStreamURL(streamURL, response.url, api),
        abortSignal,
        fetchClient,
      );
    },

    async reconnectToStream({ chatId, headers: requestHeaders }) {
      const endpoint =
        reconnectApi ?? `${api.replace(/\/$/, "")}/${chatId}/stream`;
      const response = await fetchClient(endpoint, {
        method: "GET",
        headers: {
          ...mergeHeaders(headers),
          ...mergeHeaders(requestHeaders),
        },
      });
      if (response.status === httpStatusNoContent) {
        return null;
      }
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `HTTP error ${response.status}`);
      }

      const streamURL = await streamURLFromResponse(response);
      if (!streamURL) {
        throw new Error("Durable stream URL missing from reconnect response.");
      }
      return readUIMessageChunks(
        resolveStreamURL(streamURL, response.url, endpoint),
        undefined,
        fetchClient,
      );
    },
  };
}

const httpStatusNoContent = 204;
