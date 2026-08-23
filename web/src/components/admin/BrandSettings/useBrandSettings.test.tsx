import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import { useBrandSettings } from "./useBrandSettings";

const revalidate = vi.fn();
const updateSettings = vi.fn();

vi.mock("@/api/client", () => ({
  handle: async (
    fn: () => Promise<void>,
    options?: { cleanup?: () => Promise<void> },
  ) => {
    await fn();
    await options?.cleanup?.();
  },
}));

vi.mock("@/api/openapi-client/misc", () => ({
  iconUpload: vi.fn(),
}));

vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ revalidate, updateSettings }),
}));

describe("useBrandSettings", () => {
  beforeEach(() => {
    revalidate.mockReset();
    updateSettings.mockReset().mockResolvedValue(undefined);
    vi.stubGlobal(
      "fetch",
      vi.fn(() => new Promise(() => undefined)),
    );
  });

  it("stages MOTD removal until the form is saved", async () => {
    const settings = {
      ...DefaultSettings,
      motd: {
        content: "Community maintenance tonight.",
        metadata: { type: "information" as const },
      },
    };
    const { result } = renderHook(() => useBrandSettings({ settings }));

    expect(result.current.canClearMotd).toBe(true);
    expect(result.current.motdClearPending).toBe(false);

    act(() => result.current.onClearMotd());

    expect(result.current.canClearMotd).toBe(false);
    expect(result.current.motdClearPending).toBe(true);
    expect(updateSettings).not.toHaveBeenCalled();

    await act(async () => {
      await result.current.onSubmit();
    });

    expect(updateSettings).toHaveBeenCalledWith(
      expect.objectContaining({ motd: {} }),
    );
    expect(result.current.motdClearPending).toBe(false);
  });
});
