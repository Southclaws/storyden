import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { Button } from ".";

describe("Button", () => {
  it("renders children when not loading even if loadingText is provided", () => {
    render(
      <Button loading={false} loadingText="Retrying...">
        Retry Now
      </Button>,
    );

    expect(screen.getByRole("button")).toHaveTextContent("Retry Now");
    expect(screen.queryByText("Retrying...")).not.toBeInTheDocument();
  });

  it("renders loadingText when loading is true", () => {
    render(
      <Button loading loadingText="Retrying...">
        Retry Now
      </Button>,
    );

    expect(screen.getByRole("button")).toHaveTextContent("Retrying...");
    expect(screen.queryByText("Retry Now")).not.toBeInTheDocument();
  });

  it("renders a spinner when loading is true without loadingText", () => {
    const { container } = render(<Button loading>Retry Now</Button>);

    expect(screen.getByRole("button", { name: "Retry Now" })).toBeDisabled();
    expect(screen.getByRole("button")).toHaveAttribute("aria-busy", "true");
    expect(screen.getByText("Retry Now")).toHaveClass(
      "button__loading-content",
    );
    expect(
      container.querySelector(".button__loading-indicator"),
    ).not.toBeNull();
  });

  it("applies semantic intent independently from visual priority", () => {
    render(
      <Button variant="outline" intent="destructive">
        Delete account
      </Button>,
    );

    expect(screen.getByRole("button")).toHaveClass(
      "button--variant_outline",
      "button--intent_destructive",
    );
  });
});
