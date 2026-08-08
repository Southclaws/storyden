import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { useRobotModelsList } from "@/api/openapi-client/robots";
import type { Robot } from "@/api/openapi-schema";

import { RobotConfigurationForm } from "./RobotConfigurationForm";

vi.mock("@/api/openapi-client/robots", async (importOriginal) => {
  const actual =
    await importOriginal<typeof import("@/api/openapi-client/robots")>();

  return {
    ...actual,
    useRobotModelsList: vi.fn(),
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
  });

  it("prefills the persisted model when editing a robot", () => {
    render(<RobotConfigurationForm robot={robot} />);

    expect(screen.getByPlaceholderText("Select a model")).toHaveValue(
      "openai/gpt-4-0613",
    );
  });
});
