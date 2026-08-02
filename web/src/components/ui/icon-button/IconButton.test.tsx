import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { IconButton } from ".";

describe("IconButton", () => {
  it("uses its accessible label as its title", () => {
    render(<IconButton aria-label="Create">+</IconButton>);

    expect(screen.getByRole("button", { name: "Create" })).toHaveAttribute(
      "title",
      "Create",
    );
  });

  it("preserves an explicit title", () => {
    render(
      <IconButton aria-label="Create page" title="Create a new page">
        +
      </IconButton>,
    );

    expect(screen.getByRole("button", { name: "Create page" })).toHaveAttribute(
      "title",
      "Create a new page",
    );
  });
});
