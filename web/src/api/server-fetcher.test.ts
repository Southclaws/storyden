import { cookies, headers } from "next/headers";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetcher } from "./server";

vi.mock("@/config", () => ({ getAPIAddress: () => "https://api.example" }));
vi.mock("next/headers", () => ({
  cookies: vi.fn(async () => ({
    getAll: () => [{ name: "session", value: "admin-session" }],
  })),
  headers: vi.fn(async () => new Headers({ "x-forwarded-for": "203.0.113.1" })),
}));

beforeEach(() => vi.restoreAllMocks());

describe("server API requests", () => {
  it("continues forwarding the session and client address by default", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(Response.json({ ok: true }));

    await fetcher("/info", { method: "GET" });

    const request = fetchMock.mock.calls[0]![0] as Request;
    expect(request.headers.get("Cookie")).toBe("session=admin-session");
    expect(request.headers.get("x-forwarded-for")).toBe("203.0.113.1");
    expect(request.headers.get("x-storyden-ssr")).toBe("1");
  });

  it("makes anonymous requests without credentials or request-context access", async () => {
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValue(Response.json({ ok: true }));
    const controller = new AbortController();

    const result = await fetcher("/nodes", {
      method: "GET",
      anonymous: true,
      cache: "no-store",
      signal: controller.signal,
      headers: {
        Cookie: "session=admin-session",
        Authorization: "Bearer admin-token",
      },
    });

    const request = fetchMock.mock.calls[0]![0] as Request;
    expect(request.url).toBe("https://api.example/api/nodes");
    expect(request.cache).toBe("no-store");
    expect(request.headers.has("Cookie")).toBe(false);
    expect(request.headers.has("Authorization")).toBe(false);
    expect(request.headers.has("x-forwarded-for")).toBe(false);
    expect(request.headers.has("x-storyden-ssr")).toBe(false);
    expect(cookies).not.toHaveBeenCalled();
    expect(headers).not.toHaveBeenCalled();
    expect(result).toMatchObject({ data: { ok: true }, status: 200 });

    controller.abort();
    expect(request.signal.aborted).toBe(true);
  });
});
