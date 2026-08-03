import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { PageLayout } from ".";

describe("PageLayout", () => {
  it("uses the constrained content width by default", () => {
    render(<PageLayout data-testid="layout" />);

    expect(screen.getByTestId("layout")).toHaveAttribute(
      "data-page-width",
      "content",
    );
  });

  it.each(["wide", "full"] as const)(
    "exposes the %s shell width marker",
    (width) => {
      render(<PageLayout data-testid="layout" width={width} />);

      expect(screen.getByTestId("layout")).toHaveAttribute(
        "data-page-width",
        width,
      );
    },
  );
});
