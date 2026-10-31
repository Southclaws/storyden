import "server-only";

import { getAPIAddress } from "@/config";

type ImagePath = "/info/banner" | "/info/icon/32x32" | `/assets/${string}`;

const IMAGE_TYPES = new Set([
  "image/png",
  "image/jpeg",
  "image/webp",
  "image/gif",
]);
const RESPONSE_HEADERS = [
  "Content-Type",
  "Cache-Control",
  "ETag",
  "Last-Modified",
];
const MAX_RENDER_IMAGE_BYTES = 5_000_000;

export async function streamMetadataImage(path: ImagePath): Promise<Response> {
  const signal = AbortSignal.timeout(10_000);

  try {
    const response = await fetch(
      `${getAPIAddress().replace(/\/$/, "")}/api${path}`,
      {
        cache: "no-store",
        redirect: "error",
        signal,
        headers: {
          Accept: "image/png,image/jpeg,image/webp,image/gif",
          "Accept-Encoding": "identity",
        },
      }
    );
    const contentType = response.headers
      .get("Content-Type")
      ?.split(";")[0]
      ?.trim();

    if (!response.ok || !response.body || !IMAGE_TYPES.has(contentType ?? "")) {
      await response.body?.cancel();
      return imageError(response.status >= 400 ? response.status : 502);
    }

    const headers = new Headers();
    for (const name of RESPONSE_HEADERS) {
      const value = response.headers.get(name);
      if (value) {
        headers.set(name, value);
      }
    }
    if (!headers.has("Cache-Control")) {
      headers.set("Cache-Control", "no-cache");
    }

    // Fetch decodes compressed bodies, so their original length may be wrong.
    if (!response.headers.has("Content-Encoding")) {
      const length = response.headers.get("Content-Length");
      if (length) headers.set("Content-Length", length);
    }

    return new Response(response.body, { headers });
  } catch (error) {
    console.error("Metadata image request failed", path, error);
    return imageError(signal.aborted ? 504 : 502);
  }
}

export async function loadMetadataImage(
  path: `/assets/${string}`
): Promise<ArrayBuffer | Response> {
  const response = await streamMetadataImage(path);
  if (!response.ok) return response;

  const reader = response.body!.getReader();
  const chunks: Uint8Array[] = [];
  let length = 0;

  try {
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      length += value.byteLength;
      if (length > MAX_RENDER_IMAGE_BYTES) {
        await reader.cancel();
        return imageError(502);
      }
      chunks.push(value);
    }
  } catch (error) {
    console.error("Metadata image body failed", path, error);
    return imageError(502);
  } finally {
    reader.releaseLock();
  }

  const bytes = new Uint8Array(length);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.byteLength;
  }
  return bytes.buffer;
}

function imageError(status: number): Response {
  return new Response("Image unavailable", {
    status,
    headers: { "Cache-Control": "no-store" },
  });
}
