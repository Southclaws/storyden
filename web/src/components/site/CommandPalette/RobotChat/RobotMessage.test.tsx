import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { Asset } from "@/api/openapi-schema";

import { RobotMessage } from "./RobotMessage";

vi.mock("@/components/asset/AssetThumbnailList", () => ({
  AssetThumbnailList: ({
    assets,
    align,
  }: {
    assets: Asset[];
    align?: "start" | "end";
  }) => (
    <div role="group" aria-label="Message images" data-align={align}>
      {assets.map((asset) => (
        <span key={asset.id}>{asset.filename}</span>
      ))}
    </div>
  ),
}));

describe("RobotMessage", () => {
  it("renders every attached image in the message thumbnail list", () => {
    render(
      <RobotMessage
        id="message-with-images"
        role="user"
        parts={[{ type: "text", text: "What is in these images?" }]}
        assets={[
          {
            id: "d4fa1vkrvimc73eq4jpg",
            filename: "first.png",
            path: "/api/assets/d4fa1vkrvimc73eq4jpg-first.png",
            mime_type: "image/png",
            width: 800,
            height: 600,
          },
          {
            id: "d4fa20crvimc73eq4jq0",
            filename: "second.webp",
            path: "/api/assets/d4fa20crvimc73eq4jq0-second.webp",
            mime_type: "image/webp",
            width: 1200,
            height: 900,
          },
        ]}
      />,
    );

    expect(screen.getByText("What is in these images?")).toBeInTheDocument();
    expect(
      screen.getByRole("group", { name: "Message images" }),
    ).toHaveTextContent("first.png");
    expect(
      screen.getByRole("group", { name: "Message images" }),
    ).toHaveTextContent("second.webp");
    expect(
      screen.getByRole("group", { name: "Message images" }),
    ).toHaveAttribute("data-align", "end");
  });
});
