import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { Card } from ".";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ prefetch: vi.fn() }),
}));

describe("Card accessibility", () => {
  it("offers one named destination while keeping the excerpt readable", () => {
    render(
      <Card
        id="harbour"
        title="Harbour notes"
        text="Boats are arriving."
        url="/t/harbour"
        titleStatus="Pinned discussion"
      />,
    );

    expect(
      screen.getByRole("article", {
        name: "Harbour notes Pinned discussion",
      }),
    ).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2 })).toHaveTextContent(
      "Harbour notes",
    );
    expect(screen.getAllByRole("link")).toHaveLength(1);
    expect(screen.getByRole("link", { name: "Harbour notes" })).toHaveAttribute(
      "href",
      "/t/harbour",
    );
    expect(screen.getByText("Boats are arriving.").closest("a")).toBeNull();
  });

  it("announces a meaningful foreground image once", () => {
    const { container } = render(
      <Card
        id="photo"
        title="Harbour photo"
        url="/t/photo"
        image="/photo.png"
        imageAlt="A red boat beside the wooden pier"
        headingAs="h3"
      />,
    );
    expect(screen.getByRole("heading", { level: 3 })).toBeInTheDocument();
    expect(screen.getAllByRole("img")).toHaveLength(1);
    expect(screen.getByRole("img")).toHaveAccessibleName(
      "A red boat beside the wooden pier",
    );
    expect(container.querySelector('img[aria-hidden="true"]')).toHaveAttribute(
      "alt",
      "",
    );
  });

  it("keeps decorative teasers out of image navigation", () => {
    render(
      <Card id="teaser" title="Harbour" url="/t/harbour" image="/photo.png" />,
    );
    expect(screen.queryByRole("img")).not.toBeInTheDocument();
  });

  it("preserves a named content link when fill hides the title slot", () => {
    render(
      <Card
        id="fill"
        title="Alpine field notes"
        text="A quiet lake."
        url="/l/alpine"
        shape="fill"
      />,
    );
    expect(screen.getByText("A quiet lake.").closest("a")).toHaveAccessibleName(
      "Alpine field notes",
    );
    expect(screen.getByText("A quiet lake.").closest("a")).toHaveAttribute(
      "href",
      "/l/alpine",
    );
  });

  it("expands and collapses a long preview from the keyboard without losing focus", async () => {
    vi.spyOn(HTMLElement.prototype, "scrollHeight", "get").mockReturnValue(220);
    const user = userEvent.setup();
    render(
      <Card id="long" title="Field notes" text="Long preview" url="/t/notes" />,
    );

    const disclosure = screen.getByRole("button", {
      name: "Show more of Field notes",
      expanded: false,
    });
    expect(
      document.getElementById(disclosure.getAttribute("aria-controls")!),
    ).toHaveTextContent("Long preview");
    await user.tab();
    expect(screen.getByRole("link")).toHaveFocus();
    await user.tab();
    expect(disclosure).toHaveFocus();
    await user.keyboard("{Enter}");
    expect(disclosure).toHaveAttribute("aria-expanded", "true");
    expect(disclosure).toHaveAccessibleName("Show less of Field notes");
    expect(disclosure).toHaveFocus();
    await user.keyboard(" ");
    expect(disclosure).toHaveAttribute("aria-expanded", "false");
    expect(disclosure).toHaveFocus();
  });
});
