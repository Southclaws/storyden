import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import Icon from "@/app/icon";
import Image from "@/app/opengraph-image";
import { getAPIAddress } from "@/config";

import { loadMetadataImage, streamMetadataImage } from "./image-response";

vi.mock("server-only", () => ({}));
vi.mock("@/config", () => ({ getAPIAddress: vi.fn() }));

beforeEach(() => {
  vi.mocked(getAPIAddress).mockReturnValue("http://127.0.0.1:8000/backend/");
});
afterEach(() => vi.restoreAllMocks());

describe("metadata image streaming", () => {
  it("serves chunks before the backend finishes and preserves bytes and cache policy", async () => {
    let source!: ReadableStreamDefaultController<Uint8Array>;
    const cancel = vi.fn();
    const body = new ReadableStream<Uint8Array>({
      start(controller) {
        source = controller;
      },
      cancel,
    });
    const upstream = new Response(body, {
      headers: {
        "Content-Type": "image/png",
        "Content-Length": "5",
        "Cache-Control": "public, no-cache",
        ETag: '"banner-v2"',
        "Set-Cookie": "session=private",
        "Content-Disposition": "attachment",
      },
    });
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(upstream);
    const response = await Image();
    expect(response.body).toBe(upstream.body);
    expect(fetchMock).toHaveBeenCalledWith(
      "http://127.0.0.1:8000/backend/api/info/banner",
      expect.objectContaining({
        cache: "no-store",
        redirect: "error",
        signal: expect.any(AbortSignal),
      }),
    );
    expect(response.headers.get("Cache-Control")).toBe("public, no-cache");
    expect(response.headers.get("Content-Length")).toBe("5");
    expect(response.headers.get("ETag")).toBe('"banner-v2"');
    expect(response.headers.has("Set-Cookie")).toBe(false);
    expect(response.headers.has("Content-Disposition")).toBe(false);

    const reader = response.body!.getReader();
    source.enqueue(new Uint8Array([0, 128, 255]));
    expect((await reader.read()).value).toEqual(new Uint8Array([0, 128, 255]));
    source.enqueue(new Uint8Array([1, 2]));
    source.close();
    expect((await reader.read()).value).toEqual(new Uint8Array([1, 2]));
    expect((await reader.read()).done).toBe(true);
    expect(cancel).not.toHaveBeenCalled();
  });

  it("uses the configured backend at request time for the icon", async () => {
    vi.mocked(getAPIAddress).mockReturnValue("http://backend:8000");
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response("png", {
        headers: { "Content-Type": "image/png" },
      }),
    );
    expect(await (await Icon()).text()).toBe("png");
    expect(fetchMock.mock.calls[0]![0]).toBe(
      "http://backend:8000/api/info/icon/32x32",
    );
    expect(fetchMock.mock.calls[0]![1]?.headers).not.toHaveProperty("Cookie");
  });

  it("cancels the upstream body when the consumer disconnects", async () => {
    const cancel = vi.fn();
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(new ReadableStream({ cancel }), {
        headers: { "Content-Type": "image/png" },
      }),
    );
    const response = await streamMetadataImage("/info/banner");
    await response.body!.cancel("disconnected");
    expect(cancel).toHaveBeenCalledWith("disconnected");
  });

  it("does not forward compressed lengths after fetch decodes the body", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response("decoded", {
        headers: {
          "Content-Type": "image/png",
          "Content-Encoding": "gzip",
          "Content-Length": "3",
        },
      }),
    );
    const response = await streamMetadataImage("/info/banner");
    expect(response.headers.has("Content-Encoding")).toBe(false);
    expect(response.headers.has("Content-Length")).toBe(false);
    expect(await response.text()).toBe("decoded");
  });

  it.each([404, 429, 503])(
    "preserves an upstream %s without caching an empty success",
    async (status) => {
      const cancel = vi.fn();
      vi.spyOn(globalThis, "fetch").mockResolvedValue(
        new Response(new ReadableStream({ cancel }), {
          status,
          headers: {
            "Content-Type": "text/html",
            "Cache-Control": "public, max-age=600",
          },
        }),
      );
      const response = await streamMetadataImage("/info/banner");
      expect(response.status).toBe(status);
      expect(response.headers.get("Cache-Control")).toBe("no-store");
      expect(await response.text()).toBe("Image unavailable");
      expect(cancel).toHaveBeenCalled();
    },
  );

  it("rejects a successful HTML error page", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response("login", { headers: { "Content-Type": "text/html" } }),
    );
    expect((await streamMetadataImage("/info/banner")).status).toBe(502);
  });

  it("reports connection failures and timeouts as errors", async () => {
    vi.spyOn(console, "error").mockImplementation(() => {});
    vi.spyOn(globalThis, "fetch").mockRejectedValue(
      new TypeError("connection refused"),
    );
    expect((await streamMetadataImage("/info/banner")).status).toBe(502);
    vi.spyOn(AbortSignal, "timeout").mockReturnValue(
      AbortSignal.abort(new DOMException("timeout", "TimeoutError")),
    );
    expect((await streamMetadataImage("/info/banner")).status).toBe(504);
    expect(AbortSignal.timeout).toHaveBeenCalledWith(10_000);
  });
});

describe("images used in composed previews", () => {
  it("returns binary bytes without base64 conversion", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(new Uint8Array([0, 128, 255]), {
        headers: { "Content-Type": "image/png" },
      }),
    );
    const image = await loadMetadataImage("/assets/image.png?variant=l");
    expect(image).toBeInstanceOf(ArrayBuffer);
    expect(new Uint8Array(image as ArrayBuffer)).toEqual(
      new Uint8Array([0, 128, 255]),
    );
  });

  it("accepts a source above 5 MB within the API's default upload budget", async () => {
    const bytes = new Uint8Array(5_688_566);
    bytes[bytes.length - 1] = 255;
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(bytes, { headers: { "Content-Type": "image/png" } }),
    );
    const image = await loadMetadataImage("/assets/legacy-cover.png?variant=l");
    expect(image).toBeInstanceOf(ArrayBuffer);
    expect((image as ArrayBuffer).byteLength).toBe(bytes.length);
    expect(new Uint8Array(image as ArrayBuffer).at(-1)).toBe(255);
  });

  it("limits chunked input even without a Content-Length", async () => {
    const cancel = vi.fn();
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(
        new ReadableStream({
          start(controller) {
            controller.enqueue(new Uint8Array(8_388_609));
          },
          cancel,
        }),
        { headers: { "Content-Type": "image/png" } },
      ),
    );
    const result = await loadMetadataImage("/assets/image.png?variant=l");
    expect(result).toBeInstanceOf(Response);
    expect((result as Response).status).toBe(502);
    expect(cancel).toHaveBeenCalled();
  });
});
