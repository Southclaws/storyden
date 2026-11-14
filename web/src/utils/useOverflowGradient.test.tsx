import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { useOverflowGradient } from "./useOverflowGradient";

let resize: ResizeObserverCallback;

beforeEach(() => {
  vi.stubGlobal(
    "ResizeObserver",
    class {
      constructor(callback: ResizeObserverCallback) {
        resize = callback;
      }
      observe() {}
      disconnect() {}
    },
  );
});

afterEach(() => vi.unstubAllGlobals());

function Viewport({
  enabled = true,
  direction = "ltr",
}: {
  enabled?: boolean;
  direction?: "ltr" | "rtl";
}) {
  const ref = useOverflowGradient<HTMLDivElement>(enabled);

  return (
    <div ref={ref} role="region" aria-label="Scrollable views" dir={direction}>
      Content
    </div>
  );
}

describe("overflow gradients", () => {
  for (const axis of ["block", "inline"] as const) {
    it(`shows fades only where more ${axis} content exists, and updates after resizing`, async () => {
      render(<Viewport />);
      const viewport = screen.getByRole("region");
      const scrollSize = axis === "inline" ? "scrollWidth" : "scrollHeight";
      const clientSize = axis === "inline" ? "clientWidth" : "clientHeight";
      const offset = axis === "inline" ? "scrollLeft" : "scrollTop";
      const start = `data-overflow-${axis}-start`;
      const end = `data-overflow-${axis}-end`;

      Object.defineProperties(viewport, {
        [scrollSize]: { configurable: true, value: 600 },
        [clientSize]: { configurable: true, value: 300 },
      });
      fireEvent.scroll(viewport);
      await waitFor(() => {
        expect(viewport).toHaveAttribute(start, "false");
        expect(viewport).toHaveAttribute(end, "true");
      });

      viewport[offset] = 150;
      fireEvent.scroll(viewport);
      await waitFor(() => {
        expect(viewport).toHaveAttribute(start, "true");
        expect(viewport).toHaveAttribute(end, "true");
      });

      viewport[offset] = 300;
      fireEvent.scroll(viewport);
      await waitFor(() => {
        expect(viewport).toHaveAttribute(start, "true");
        expect(viewport).toHaveAttribute(end, "false");
      });

      Object.defineProperty(viewport, clientSize, { value: 600 });
      viewport[offset] = 0;
      resize([], {} as ResizeObserver);
      await waitFor(() => {
        expect(viewport).toHaveAttribute(start, "false");
        expect(viewport).toHaveAttribute(end, "false");
      });
    });
  }

  it("tracks logical edges when scrolling right to left", async () => {
    render(<Viewport direction="rtl" />);
    const viewport = screen.getByRole("region");
    Object.defineProperties(viewport, {
      scrollWidth: { value: 600 },
      clientWidth: { value: 300 },
    });
    viewport.scrollLeft = -300;
    fireEvent.scroll(viewport);

    await waitFor(() => {
      expect(viewport).toHaveAttribute("data-overflow-inline-start", "true");
      expect(viewport).toHaveAttribute("data-overflow-inline-end", "false");
    });
  });

  it("clears fades when disabled", async () => {
    const { rerender } = render(<Viewport />);
    const viewport = screen.getByRole("region");
    await waitFor(() => {
      expect(viewport).toHaveAttribute("data-overflow-inline-end", "false");
    });

    rerender(<Viewport enabled={false} />);
    expect(viewport).not.toHaveAttribute("data-overflow-inline-start");
    expect(viewport).not.toHaveAttribute("data-overflow-inline-end");
  });
});
