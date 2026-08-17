import { render, screen } from "@testing-library/react";
import type { UIDataTypes, UIMessagePart } from "ai";
import { describe, expect, it, vi } from "vitest";

import type { StorydenTools } from "@/api/robots";

import { RobotToolCall } from "./RobotToolCall";

vi.mock("@/api/openapi-client/nodes", () => ({
  useNodeList: () => ({
    data: { nodes: [] },
    error: undefined,
    isLoading: false,
  }),
}));

vi.mock("./RobotChatContext", () => ({
  useRobotChat: () => ({
    messages: [],
    resolveLibraryPageRequest: vi.fn(),
    resolveToolConfirmation: vi.fn(),
  }),
}));

describe("RobotToolCall", () => {
  it("renders an error-shaped tool result without reading the success payload", () => {
    const message =
      'tool "document_get" is Toolset-only; use Toolset system.documents instead';
    const part = {
      type: "tool-tool_load",
      toolCallId: "call-tool-load",
      state: "output-available",
      input: { tools: ["document_get"] },
      output: { error: message },
    } as unknown as UIMessagePart<UIDataTypes, StorydenTools>;

    const { container } = render(<RobotToolCall part={part} />);

    expect(screen.getByRole("alert")).toHaveTextContent(
      "“document_get” can’t be loaded individually; load system.documents instead.",
    );
    expect(container.querySelector("pre")).toHaveTextContent("Toolset-only");
    expect(container.querySelector("pre")).toHaveTextContent(
      "system.documents",
    );
    expect(screen.getByText("Error")).toBeInTheDocument();
    expect(screen.queryByText(/Loaded .* tools/)).not.toBeInTheDocument();
  });

  it("summarises verbose runtime diagnostics while preserving full details", () => {
    const message = [
      "tool 'web_open' not found.",
      "Available tools: document_get, document_search, tool_search",
      "Possible causes: the tool has not been loaded",
    ].join("\n");
    const part = {
      type: "tool-web_open",
      toolCallId: "call-web-open",
      state: "output-available",
      input: { url: "https://example.com" },
      output: { error: message },
    } as unknown as UIMessagePart<UIDataTypes, StorydenTools>;

    const { container } = render(<RobotToolCall part={part} />);

    expect(screen.getByRole("alert")).toHaveTextContent(
      "“web_open” isn’t available in this conversation.",
    );
    expect(screen.getByRole("alert")).not.toHaveTextContent("Available tools");
    expect(container.querySelector("pre")).toHaveTextContent(
      "Available tools: document_get, document_search, tool_search",
    );
  });
});
