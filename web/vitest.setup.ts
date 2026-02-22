import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

if (typeof window !== "undefined" && !(window as any).__storyden__) {
  (window as any).__storyden__ = {
    API_ADDRESS: "http://localhost:8000",
    WEB_ADDRESS: "http://localhost:3000",
    source: "script",
  };
}

afterEach(() => {
  cleanup();
});

if (typeof globalThis.ResizeObserver === "undefined") {
  class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
  }

  globalThis.ResizeObserver = ResizeObserver as typeof globalThis.ResizeObserver;
}

if (typeof document.execCommand !== "function") {
  (document as Document & { execCommand?: typeof document.execCommand }).execCommand =
    () => false;
}
