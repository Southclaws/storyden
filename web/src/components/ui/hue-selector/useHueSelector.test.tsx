import { act, renderHook } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { useHueSelector } from "./useHueSelector";

describe("useHueSelector", () => {
  test("reads an existing colour and commits its hue on release", () => {
    const onChange = vi.fn();
    const { result, rerender } = renderHook(
      ({ value }) => useHueSelector({ value, onChange, onUpdate: vi.fn() }),
      { initialProps: { value: "#00ff00" } },
    );
    expect(result.current.value).toBe("hsl(120, 15%, 80%)");
    rerender({ value: "rgb(0, 0, 255)" });
    expect(result.current.value).toBe("hsl(240, 15%, 80%)");
    act(() => result.current.onPointerUp());
    expect(onChange).toHaveBeenCalledWith("hsl(240, 15%, 80%)");
  });
});
