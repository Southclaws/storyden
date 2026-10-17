import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { register } from "./instrumentation";

const addressSettings = [
  "PUBLIC_API_ADDRESS",
  "PUBLIC_WEB_ADDRESS",
  "NEXT_PUBLIC_API_ADDRESS",
  "NEXT_PUBLIC_WEB_ADDRESS",
  "SSR_API_ADDRESS",
] as const;

describe("server address configuration", () => {
  beforeEach(() => {
    vi.resetModules();

    for (const setting of addressSettings) {
      vi.stubEnv(setting, undefined);
    }
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.restoreAllMocks();
  });

  it("preserves localhost defaults when addresses are unset", async () => {
    const { serverEnvironment, getAPIAddress } = await import("./config");

    expect(serverEnvironment()).toEqual({
      API_ADDRESS: "http://localhost:8000",
      WEB_ADDRESS: "http://localhost:3000",
      source: "server",
    });
    expect(getAPIAddress()).toBe("http://localhost:8000");
  });

  it("preserves NEXT_PUBLIC address precedence", async () => {
    vi.stubEnv("PUBLIC_API_ADDRESS", "https://api.example.com");
    vi.stubEnv("PUBLIC_WEB_ADDRESS", "https://example.com");
    vi.stubEnv("NEXT_PUBLIC_API_ADDRESS", "https://api.community.example");
    vi.stubEnv("NEXT_PUBLIC_WEB_ADDRESS", "https://community.example");

    const { serverEnvironment } = await import("./config");

    expect(serverEnvironment()).toEqual({
      API_ADDRESS: "https://api.community.example",
      WEB_ADDRESS: "https://community.example",
      source: "server",
    });
  });

  it.each([
    "http://localhost:8000",
    "http://storyden:8000",
    "http://192.168.1.5:8000",
    "http://[fd00::1]:8000",
    "https://community.example/backend",
  ])("accepts a complete internal or public address: %s", async (address) => {
    vi.stubEnv("PUBLIC_API_ADDRESS", address);
    vi.stubEnv("PUBLIC_WEB_ADDRESS", address);
    vi.stubEnv("SSR_API_ADDRESS", address);

    const { serverEnvironment, getAPIAddress } = await import("./config");

    expect(serverEnvironment().API_ADDRESS).toBe(address);
    expect(serverEnvironment().WEB_ADDRESS).toBe(address);
    expect(getAPIAddress()).toBe(address);
  });

  it("allows an internal HTTP API behind public HTTPS addresses", async () => {
    vi.stubEnv("PUBLIC_API_ADDRESS", "https://api.example.com");
    vi.stubEnv("PUBLIC_WEB_ADDRESS", "https://example.com");
    vi.stubEnv("SSR_API_ADDRESS", "http://backend:8000");

    const { serverEnvironment, getAPIAddress } = await import("./config");

    expect(serverEnvironment().API_ADDRESS).toBe("https://api.example.com");
    expect(getAPIAddress()).toBe("http://backend:8000");
    await expect(register()).resolves.toBeUndefined();
  });

  it("leaves valid but incompatible addresses for Doctor to explain", async () => {
    vi.stubEnv("PUBLIC_API_ADDRESS", "http://api.example.com/backend");
    vi.stubEnv("PUBLIC_WEB_ADDRESS", "https://other.example/community");

    await expect(register()).resolves.toBeUndefined();
  });

  describe.each(addressSettings)("%s", (setting) => {
    it.each([
      "",
      "community.example",
      "/api",
      "ftp://example.com",
      "https:example.com",
      "https://",
      "https://operator:secret@example.com",
    ])("rejects invalid configuration at startup: %s", async (address) => {
      vi.stubEnv(setting, address);

      await expect(import("./config")).rejects.toThrow(
        `Invalid ${setting}: expected an absolute HTTP or HTTPS URL without a username or password.`,
      );
    });
  });

  it("exits when startup configuration validation fails", async () => {
    vi.stubEnv("SSR_API_ADDRESS", "backend:8000");
    const stopped = new Error("Process exited");
    const exit = vi.spyOn(process, "exit").mockImplementation(() => {
      throw stopped;
    });
    const log = vi.spyOn(console, "error").mockImplementation(() => {});

    await expect(register()).rejects.toBe(stopped);

    expect(exit).toHaveBeenCalledWith(1);
    expect(log).toHaveBeenCalledWith(
      expect.objectContaining({
        message:
          "Invalid SSR_API_ADDRESS: expected an absolute HTTP or HTTPS URL without a username or password.",
      }),
    );
  });

  it("validates each serverEnvironment call without disclosing the value", async () => {
    const { serverEnvironment } = await import("./config");
    vi.stubEnv("PUBLIC_API_ADDRESS", "https://operator:secret@example.com");

    expect(serverEnvironment).toThrow(
      /^Invalid PUBLIC_API_ADDRESS: expected an absolute HTTP or HTTPS URL without a username or password\.$/,
    );
  });

  it("applies the same address rules to injected browser configuration", async () => {
    const { ConfigSchema } = await import("./config");

    expect(
      ConfigSchema.safeParse({
        API_ADDRESS: "api.example.com",
        WEB_ADDRESS: "https://example.com",
        source: "script",
      }).success,
    ).toBe(false);
  });
});
