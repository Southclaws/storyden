import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { ReactList, type ReactListItem } from "./ReactList.internal";

function reaction(): ReactListItem {
  return {
    emoji: "🔥",
    count: 2,
    hasReacted: false,
    reactedBy: ["odin", "southclaws"],
  };
}

describe("ReactList", () => {
  test("renders guest reactions as genuinely disabled controls", () => {
    render(
      <ReactList
        isLoggedIn={false}
        reactions={[reaction()]}
        onReactExisting={vi.fn()}
        onReactPicker={vi.fn()}
      />,
    );

    expect(screen.getByTitle("🔥: odin, southclaws")).toBeDisabled();
    expect(
      screen.queryByRole("button", { name: "Add reaction" }),
    ).not.toBeInTheDocument();
  });

  test("marks a reaction control busy while its request is pending", () => {
    render(
      <ReactList
        isLoggedIn
        pendingEmojis={new Set(["🔥"])}
        reactions={[reaction()]}
        onReactExisting={vi.fn()}
        onReactPicker={vi.fn()}
      />,
    );

    const button = screen.getByTitle("🔥: odin, southclaws");
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute("aria-busy", "true");
  });

  test("updates the visible count immediately before the request resolves", () => {
    const onReactExisting = vi.fn();
    render(
      <ReactList
        isLoggedIn
        reactions={[reaction()]}
        onReactExisting={onReactExisting}
        onReactPicker={vi.fn()}
      />,
    );

    const button = screen.getByTitle("🔥: odin, southclaws");
    fireEvent.click(button);

    expect(button).toHaveTextContent("🔥3");
    expect(screen.getByText("3")).toHaveAttribute("data-direction", "1");
    expect(onReactExisting).toHaveBeenCalledOnce();
    expect(onReactExisting).toHaveBeenCalledWith("🔥");
  });

  test("reconciles the animated count with updated data", () => {
    const props = {
      isLoggedIn: true,
      onReactExisting: vi.fn(),
      onReactPicker: vi.fn(),
    };
    const { rerender } = render(
      <ReactList {...props} reactions={[reaction()]} />,
    );

    rerender(
      <ReactList {...props} reactions={[{ ...reaction(), count: 4 }]} />,
    );

    expect(screen.getByTitle("🔥: odin, southclaws")).toHaveTextContent("🔥4");
  });

  test("animates removing a reaction in the opposite direction", () => {
    const onReactExisting = vi.fn();
    render(
      <ReactList
        isLoggedIn
        reactions={[{ ...reaction(), hasReacted: true }]}
        onReactExisting={onReactExisting}
        onReactPicker={vi.fn()}
      />,
    );
    fireEvent.click(screen.getByTitle("🔥: odin, southclaws"));
    expect(screen.getByText("1")).toHaveAttribute("data-direction", "-1");
    expect(onReactExisting).toHaveBeenCalledWith("🔥");
  });

  test("removes the button when the final reaction is removed", () => {
    render(
      <ReactList
        isLoggedIn
        reactions={[{ ...reaction(), count: 1, hasReacted: true }]}
        onReactExisting={vi.fn()}
        onReactPicker={vi.fn()}
      />,
    );
    fireEvent.click(screen.getByTitle("🔥: odin, southclaws"));
    expect(screen.queryByTitle("🔥: odin, southclaws")).not.toBeInTheDocument();
  });
});
