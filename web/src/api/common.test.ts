import { describe, expect, it } from "vitest";

import { buildResult } from "./common";

describe("API response decoding", () => {
  it("returns plain text as a string, including version identifiers", async () => {
    const response = new Response("edge-abc1234", {
      headers: { "Content-Type": "text/plain; charset=utf-8" },
    });

    expect(await buildResult<string>(response)).toBe("edge-abc1234");
  });

  it("returns structured JSON responses", async () => {
    const response = Response.json({ title: "My community" });

    expect(await buildResult(response)).toEqual({ title: "My community" });
  });

  it("preserves binary responses as blobs", async () => {
    const bytes = new Uint8Array([0, 255, 17, 128]);
    const response = new Response(bytes, {
      headers: { "Content-Type": "application/octet-stream" },
    });

    const result = await buildResult<Blob>(response);

    expect(result).toBeInstanceOf(Blob);
    expect(new Uint8Array(await result.arrayBuffer())).toEqual(bytes);
  });
});
