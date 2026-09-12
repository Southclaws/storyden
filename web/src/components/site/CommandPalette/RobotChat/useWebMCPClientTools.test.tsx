import { act, renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { StorydenUIMessage } from "@/api/robots-types";

import {
  usePendingWebMCPToolCalls,
  useWebMCPClientTools,
} from "./useWebMCPClientTools";
import { getWebMCPClientID } from "./webMCPClient";

const originalModelContext = Object.getOwnPropertyDescriptor(
  document,
  "modelContext",
);

afterEach(() => {
  if (originalModelContext) {
    Object.defineProperty(document, "modelContext", originalModelContext);
  } else {
    Reflect.deleteProperty(document, "modelContext");
  }
  vi.restoreAllMocks();
});

describe("useWebMCPClientTools", () => {
  it("executes once without waiting for tool output persistence", async () => {
    const tool = {
      name: "library_page_layout_get",
      description: "Read the current library layout.",
      inputSchema: { type: "object" },
    };
    const executeTool = vi.fn().mockResolvedValue('{"blocks":[]}');
    Object.defineProperty(document, "modelContext", {
      configurable: true,
      value: {
        getTools: vi.fn().mockResolvedValue([tool]),
        executeTool,
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
      },
    });
    const submitOutput = vi.fn(() => new Promise<void>(() => undefined));
    const { result } = renderHook(() =>
      useWebMCPClientTools(async () => ({ page_type: "library" })),
    );
    act(() => {
      result.current.setToolOutputSubmitter(submitOutput);
    });
    const toolCall = {
      toolName: "library_page_layout_get",
      toolCallId: "call-1",
      input: {},
      dynamic: true,
      providerMetadata: {
        storyden: {
          source: "webmcp",
          client_id: getWebMCPClientID(),
          scope: { page_type: "library" },
        },
      },
    };

    await act(async () => {
      await expect(result.current.handleWebMCPToolCall(toolCall)).resolves.toBe(
        true,
      );
    });

    await act(async () => {
      await result.current.handleWebMCPToolCall(toolCall);
    });

    expect(executeTool).toHaveBeenCalledOnce();
    expect(executeTool).toHaveBeenCalledWith(tool, "{}", {
      signal: undefined,
    });
    expect(submitOutput).toHaveBeenCalledWith({
      toolName: "library_page_layout_get",
      toolCallId: "call-1",
      output: { blocks: [] },
    });
  });

  it("retries a pending same-browser call when its page tool remounts", async () => {
    const tool = {
      name: "library_page_block_add",
      description: "Add a block.",
      inputSchema: { type: "object" },
    };
    let registeredTools: (typeof tool)[] = [];
    let toolChangeListener: EventListener | undefined;
    const executeTool = vi.fn().mockResolvedValue('{"added":true}');
    Object.defineProperty(document, "modelContext", {
      configurable: true,
      value: {
        getTools: vi.fn().mockImplementation(() => registeredTools),
        executeTool,
        addEventListener: vi.fn(
          (_type: string, listener: EventListener) =>
            (toolChangeListener = listener),
        ),
        removeEventListener: vi.fn(),
      },
    });
    const submitOutput = vi.fn();
    const { result } = renderHook(() =>
      useWebMCPClientTools(async () => ({ page_type: "library" })),
    );
    act(() => {
      result.current.setToolOutputSubmitter(submitOutput);
    });

    await act(async () => {
      await result.current.handleWebMCPToolCall({
        toolName: "library_page_block_add",
        toolCallId: "call-1",
        input: { block: "gallery" },
        dynamic: true,
        providerMetadata: {
          storyden: {
            source: "webmcp",
            client_id: getWebMCPClientID(),
            scope: { page_type: "library" },
          },
        },
      });
    });

    expect(submitOutput).not.toHaveBeenCalled();

    registeredTools = [tool];
    act(() => toolChangeListener?.(new Event("toolchange")));

    await waitFor(() =>
      expect(submitOutput).toHaveBeenCalledWith({
        toolName: "library_page_block_add",
        toolCallId: "call-1",
        output: { added: true },
      }),
    );
  });
});

describe("usePendingWebMCPToolCalls", () => {
  it("recovers a pending WebMCP call from hydrated session history", async () => {
    const handleToolCall = vi.fn().mockResolvedValue(true);
    renderHook(() =>
      usePendingWebMCPToolCalls(
        [
          {
            id: "message-1",
            role: "assistant",
            parts: [
              {
                type: "dynamic-tool",
                toolName: "library_page_layout_get",
                toolCallId: "call-1",
                state: "input-available",
                input: {},
                callProviderMetadata: {
                  storyden: { source: "webmcp", client_id: "browser-1" },
                },
              },
            ],
          },
        ],
        handleToolCall,
      ),
    );

    await waitFor(() =>
      expect(handleToolCall).toHaveBeenCalledWith({
        toolName: "library_page_layout_get",
        toolCallId: "call-1",
        input: {},
        dynamic: true,
        providerMetadata: {
          storyden: { source: "webmcp", client_id: "browser-1" },
        },
      }),
    );
  });

  it("does not repeat a call whose output is already in hydrated history", () => {
    const handleToolCall = vi.fn().mockResolvedValue(true);
    renderHook(() =>
      usePendingWebMCPToolCalls(
        [
          {
            id: "message-1",
            role: "assistant",
            parts: [
              {
                type: "dynamic-tool",
                toolName: "library_page_layout_get",
                toolCallId: "call-1",
                state: "input-available",
                input: {},
              },
            ],
          },
          {
            id: "message-2",
            role: "user",
            parts: [
              {
                type: "tool-library_page_layout_get",
                toolCallId: "call-1",
                state: "output-available",
                input: {},
                output: { blocks: [] },
              },
            ],
          },
        ] as unknown as StorydenUIMessage[],
        handleToolCall,
      ),
    );

    expect(handleToolCall).not.toHaveBeenCalled();
  });

  it("ignores persisted messages with null parts", () => {
    const handleToolCall = vi.fn().mockResolvedValue(true);
    renderHook(() =>
      usePendingWebMCPToolCalls(
        [
          { id: "message-1", role: "assistant", parts: null },
        ] as unknown as StorydenUIMessage[],
        handleToolCall,
      ),
    );

    expect(handleToolCall).not.toHaveBeenCalled();
  });
});
