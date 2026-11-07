import { act, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import type { WebMCP } from "webmcp-types";

import { GeneratedWebMCPTool } from "./GeneratedWebMCPTool";
import { WEBMCP_TOOL_DEFINITIONS } from "./tools.generated";

const originalModelContext = Object.getOwnPropertyDescriptor(
  document,
  "modelContext",
);
const registerTool =
  vi.fn<
    (
      tool: WebMCP.ModelContextTool,
      options?: WebMCP.ModelContextRegisterToolOptions,
    ) => Promise<void>
  >();
const definition = WEBMCP_TOOL_DEFINITIONS.index_page_edit_start;

beforeEach(() => {
  registerTool.mockResolvedValue(undefined);
  Object.defineProperty(document, "modelContext", {
    configurable: true,
    value: { registerTool },
  });
});

afterEach(() => {
  if (originalModelContext) {
    Object.defineProperty(document, "modelContext", originalModelContext);
  } else {
    Reflect.deleteProperty(document, "modelContext");
  }
});

describe("GeneratedWebMCPTool", () => {
  it("keeps the latest handler registered and aborts its registration on unmount", async () => {
    const first = vi.fn(() => ({ message: "First render" }));
    const latest = vi.fn(() => ({ message: "Latest render" }));
    const view = render(
      <GeneratedWebMCPTool config={{ ...definition, execute: first }} />,
    );
    await waitFor(() => expect(registerTool).toHaveBeenCalledOnce());
    const [tool, registration] = registerTool.mock.calls[0]!;
    expect(tool).not.toHaveProperty("outputSchema");
    expect(registration?.signal?.aborted).toBe(false);

    view.rerender(
      <GeneratedWebMCPTool config={{ ...definition, execute: latest }} />,
    );

    await act(async () => {
      await expect(
        tool.execute({}, { signal: new AbortController().signal }),
      ).resolves.toEqual({ message: "Latest render" });
    });
    expect(first).not.toHaveBeenCalled();
    expect(latest).toHaveBeenCalledOnce();
    expect(registerTool).toHaveBeenCalledOnce();

    view.unmount();
    expect(registration?.signal?.aborted).toBe(true);
  });

  it("preserves a rejected edit's explanation in the browser result", async () => {
    const execute = vi.fn(() => {
      throw new Error("The assets block is already present.");
    });
    render(<GeneratedWebMCPTool config={{ ...definition, execute }} />);
    await waitFor(() => expect(registerTool).toHaveBeenCalledOnce());
    const [tool] = registerTool.mock.calls[0]!;

    await act(async () => {
      await expect(
        tool.execute({}, { signal: new AbortController().signal }),
      ).resolves.toEqual({
        isError: true,
        content: [
          { type: "text", text: "The assets block is already present." },
        ],
      });
    });
  });

  it("rejects cancellation instead of converting it to an edit error", async () => {
    const controller = new AbortController();
    const cancelled = new DOMException("Edit cancelled", "AbortError");
    const execute = vi.fn(() => {
      controller.abort(cancelled);
      throw cancelled;
    });
    render(<GeneratedWebMCPTool config={{ ...definition, execute }} />);
    await waitFor(() => expect(registerTool).toHaveBeenCalledOnce());
    const [tool] = registerTool.mock.calls[0]!;

    await act(async () => {
      await expect(
        tool.execute({}, { signal: controller.signal }),
      ).rejects.toThrow("Edit cancelled");
    });
  });
});
