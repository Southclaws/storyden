import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  useRobotModelsList,
  useRobotToolsList,
  useRobotToolsetsList,
} from "@/api/openapi-client/robots";
import type { Robot } from "@/api/openapi-schema";

import { RobotConfigurationForm } from "./RobotConfigurationForm";

vi.mock("@/api/openapi-client/robots", async (importOriginal) => {
  const actual =
    await importOriginal<typeof import("@/api/openapi-client/robots")>();

  return {
    ...actual,
    useRobotModelsList: vi.fn(),
    useRobotToolsList: vi.fn(),
    useRobotToolsetsList: vi.fn(),
  };
});

const robot: Robot = {
  author: {
    handle: "odin",
    id: "account-id",
    joined: "2026-01-01T00:00:00.000Z",
    name: "Odin",
    roles: [],
  },
  createdAt: "2026-01-01T00:00:00.000Z",
  description: "Moderates content",
  id: "robot-id",
  model: "openai/gpt-4-0613",
  name: "Moderation Bot",
  playbook: "Moderate content",
  tools: [],
  toolsets: [],
  updatedAt: "2026-01-01T00:00:00.000Z",
};

describe("RobotConfigurationForm", () => {
  beforeEach(() => {
    const modelsResponse: ReturnType<typeof useRobotModelsList> = {
      data: { models: [] },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/models"],
    };

    vi.mocked(useRobotModelsList).mockReturnValue(modelsResponse);
    const toolsResponse: ReturnType<typeof useRobotToolsList> = {
      data: { tools: [] },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/tools"],
    };

    vi.mocked(useRobotToolsList).mockReturnValue(toolsResponse);
    const toolsetsResponse: ReturnType<typeof useRobotToolsetsList> = {
      data: { toolsets: [] },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/toolsets"],
    };

    vi.mocked(useRobotToolsetsList).mockReturnValue(toolsetsResponse);
  });

  it("prefills the persisted model when editing a robot", () => {
    render(<RobotConfigurationForm robot={robot} />);

    expect(screen.getByPlaceholderText("Select a model")).toHaveValue(
      "openai/gpt-4-0613",
    );
  });

  it("removes an individual tool when a Toolset containing it is selected", async () => {
    const user = userEvent.setup();
    vi.mocked(useRobotToolsList).mockReturnValue({
      data: {
        tools: [
          {
            id: "thread_get",
            callable_name: "thread_get",
            name: "Get thread",
            description: "Get a discussion thread",
          source: "native",
          available: true,
          requires_confirmation: false,
          requires_workspace: false,
          },
        ],
      },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/tools"],
    });
    vi.mocked(useRobotToolsetsList).mockReturnValue({
      data: {
        toolsets: [
          {
            id: "system.discussions",
            name: "Discussion management",
            description: "Discussion tools",
            instruction: "Use discussion data as the source of truth.",
            tools: ["thread_get"],
          source: "system",
          editable: false,
          usage_count: 0,
          requires_workspace: false,
          },
        ],
      },
      error: undefined,
      isLoading: false,
      isValidating: false,
      mutate: vi.fn(),
      swrKey: ["/robots/toolsets"],
    });

    render(
      <RobotConfigurationForm robot={{ ...robot, tools: ["thread_get"] }} />,
    );

    await user.click(screen.getByRole("button", { name: "Select Toolsets" }));
    await user.click(
      screen.getByRole("menuitem", { name: "Discussion management" }),
    );
    await user.click(
      screen.getByRole("button", { name: "Select individual tools" }),
    );

    expect(
      screen.queryByRole("button", { name: "Remove Get thread" }),
    ).not.toBeInTheDocument();
  });
});
