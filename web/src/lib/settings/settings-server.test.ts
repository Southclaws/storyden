import { beforeEach, describe, expect, it, vi } from "vitest";

import { getInfo } from "@/api/openapi-server/misc";

import { DefaultSettings } from "./settings";
import { getSettings, getSettingsResult } from "./settings-server";

vi.mock("server-only", () => ({}));
vi.mock("@/api/openapi-server/misc", () => ({ getInfo: vi.fn() }));

describe("server settings availability", () => {
  beforeEach(() => {
    vi.mocked(getInfo).mockReset();
  });

  it("reports loaded settings as available even when their values match defaults", async () => {
    vi.mocked(getInfo).mockResolvedValue({
      data: { ...DefaultSettings },
      status: 200,
      headers: new Headers(),
    });

    expect(await getSettingsResult()).toEqual({
      ready: true,
      settings: DefaultSettings,
    });
  });

  it("reports loading failures explicitly without substituting settings", async () => {
    const error = new Error("API unavailable");
    vi.mocked(getInfo).mockRejectedValue(error);

    expect(await getSettingsResult()).toEqual({ ready: false, error });
  });

  it("preserves the fallback for callers that only need display settings", async () => {
    vi.mocked(getInfo).mockRejectedValue(new Error("API unavailable"));

    expect(await getSettings()).toEqual(DefaultSettings);
  });
});
