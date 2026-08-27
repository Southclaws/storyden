import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { useRobotToolsList } from "@/api/openapi-client/robots";
import type { RobotToolset } from "@/api/openapi-schema";

import { RobotToolsetConfigurationForm } from "./RobotToolsetConfigurationForm";

vi.mock("@/api/openapi-client/robots", async (importOriginal) => {
  const actual =
    await importOriginal<typeof import("@/api/openapi-client/robots")>();
  return { ...actual, useRobotToolsList: vi.fn() };
});

describe("RobotToolsetConfigurationForm", () => {
  beforeEach(() => {
    vi.mocked(useRobotToolsList).mockReturnValue({
      data: {
        tools: [
          {
            id: "mcp:calendar:echo",
            callable_name: "calendar_echo",
            name: "Echo event details",
            description: "Echo calendar event details.",
            source: "mcp",
            available: true,
            requires_confirmation: false,
            requires_workspace: false,
            toolset_only: false,
          },
        ],
      },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/tools"],
    });
  });

  it("presents an MCP Toolset as read-only and managed from Robot settings", () => {
    const toolset: RobotToolset = {
      id: "mcp:calendar",
      name: "Calendar",
      description: "Connected calendar capabilities.",
      instruction: "",
      tools: ["mcp:calendar:echo"],
      source: "mcp",
      source_ref: "server-id",
      editable: false,
      usage_count: 0,
      requires_workspace: false,
    };

    render(<RobotToolsetConfigurationForm toolset={toolset} />);

    expect(screen.getByText("mcp")).toBeVisible();
    expect(
      screen.getByText(
        "This Toolset mirrors a connected MCP server and is managed from the Robots settings.",
      ),
    ).toBeVisible();
    expect(screen.getByRole("textbox", { name: "Name" })).toBeDisabled();
    expect(screen.getByText("Echo event details")).toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Save" }),
    ).not.toBeInTheDocument();
  });
});
