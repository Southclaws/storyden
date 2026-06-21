import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { ContentComposerMarkdown } from "./ContentComposerMarkdown";

describe("ContentComposerMarkdown", () => {
  it("renders standalone SDR Markdown links as inline references", () => {
    render(
      <ContentComposerMarkdown
        disabled
        initialValue="Sure: [Documentation Hub](sdr:node/d8818ueot5pfij6bvm90)."
        initialValueFormat="markdown"
      />,
    );

    const link = screen.getByRole("link", { name: /Documentation Hub/ });

    expect(link).toHaveAttribute(
      "href",
      "/_/resolve/node/d8818ueot5pfij6bvm90",
    );
    expect(link.querySelector("svg")).not.toBeNull();
  });
});
