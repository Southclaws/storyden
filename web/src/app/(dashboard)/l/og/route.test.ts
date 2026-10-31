import { ImageResponse } from "next/og";
import { NextRequest } from "next/server";
import { readFileSync } from "node:fs";
import { createElement } from "react";
import { beforeAll, beforeEach, expect, it, vi } from "vitest";

import { nodeGet } from "@/api/openapi-server/nodes";

import { GET } from "./route";

vi.mock("server-only", () => ({}));
vi.mock("@/config", () => ({ getAPIAddress: () => "http://127.0.0.1:8000" }));
vi.mock("@/api/openapi-server/nodes", () => ({ nodeGet: vi.fn() }));
vi.mock("@/lib/settings/settings-server", () => ({
  getSettings: async () => ({ accent_colour: "#fff" }),
}));

beforeAll(async () => {
  await new ImageResponse(createElement("div", null, "warm up")).arrayBuffer();
});

beforeEach(() => {
  vi.mocked(nodeGet).mockResolvedValue({
    data: {
      visibility: "published",
      name: "A public page",
      description: "A description",
    },
  } as any);
});

it("requires a slug", async () => {
  expect((await GET(new NextRequest("http://example.com/l/og"))).status).toBe(
    400,
  );
  expect(nodeGet).not.toHaveBeenCalled();
});

it("only renders publicly visible pages through an anonymous read", async () => {
  vi.mocked(nodeGet).mockResolvedValue({
    data: { visibility: "draft" },
  } as any);
  const response = await GET(
    new NextRequest("http://example.com/l/og?slug=private"),
  );
  expect(response.status).toBe(404);
  expect(response.headers.get("Cache-Control")).toBe("no-store");
  expect(nodeGet).toHaveBeenCalledWith(
    "private",
    undefined,
    expect.objectContaining({ anonymous: true, cache: "no-store" }),
  );
});

it.each([false, true])(
  "renders a 1200 by 630 PNG with image=%s",
  async (withImage) => {
    if (withImage) {
      vi.mocked(nodeGet).mockResolvedValue({
        data: {
          visibility: "published",
          name: "A public page",
          description: "A description",
          primary_image: {
            filename: "image.png",
            path: "http://untrusted.example/image.png",
          },
        },
      } as any);
    }
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(
        new Response(
          readFileSync("../app/services/branding/icon/embed/default_icon.png"),
          { headers: { "Content-Type": "image/png" } },
        ),
      );
    const response = await GET(
      new NextRequest("http://example.com/l/og?slug=public"),
    );
    expect(response.status).toBe(200);
    expect(response.headers.get("Content-Type")).toBe("image/png");
    const png = Buffer.from(await response.arrayBuffer());
    expect(png.subarray(0, 8)).toEqual(
      Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]),
    );
    expect(png.readUInt32BE(16)).toBe(1200);
    expect(png.readUInt32BE(20)).toBe(630);
    expect(png.length).toBeLessThan(5_000_000);
    if (withImage) {
      expect(fetchMock.mock.calls[0]![0]).toBe(
        "http://127.0.0.1:8000/api/assets/image.png?variant=l",
      );
    } else {
      expect(fetchMock).not.toHaveBeenCalled();
    }
  },
);
