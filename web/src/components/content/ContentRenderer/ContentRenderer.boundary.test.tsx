import { render, screen } from "@testing-library/react";
import { expect, it, vi } from "vitest";

import { ContentRenderer } from "./ContentRenderer";
import { RenderedContentProvider } from "./RenderedContentContext";

const renderClient = vi.hoisted(() => vi.fn());
vi.mock("./ContentRendererBody", () => ({
  ContentRendererBody: ({ html }: { html: string }) => {
    renderClient(html);
    return <p>Updated client content</p>;
  },
}));

it("uses the server-rendered node and only loads the parser for changed content", async () => {
  const content = { "<h2>Original</h2>": <h2>Original</h2> };
  const { rerender } = render(
    <RenderedContentProvider content={content}>
      <ContentRenderer initialValue="<h2>Original</h2>" />
    </RenderedContentProvider>,
  );
  expect(screen.getByRole("heading", { name: "Original" })).toBeVisible();
  expect(renderClient).not.toHaveBeenCalled();

  rerender(
    <RenderedContentProvider content={content}>
      <ContentRenderer
        initialValue="<h2>Original</h2>"
        value="<p>Updated</p>"
      />
    </RenderedContentProvider>,
  );
  expect(await screen.findByText("Updated client content")).toBeVisible();
  expect(
    screen.queryByRole("heading", { name: "Original" }),
  ).not.toBeInTheDocument();
  expect(renderClient).toHaveBeenCalledWith("<p>Updated</p>");

  rerender(
    <RenderedContentProvider content={content}>
      <ContentRenderer initialValue="<h2>Original</h2>" value="" />
    </RenderedContentProvider>,
  );
  expect(screen.queryByText("Updated client content")).not.toBeInTheDocument();
  expect(screen.queryByRole("heading")).not.toBeInTheDocument();
});

it("keeps the empty composer shell available without loading a parser", () => {
  const { container } = render(<ContentRenderer className="theme-hook" />);
  expect(
    container.querySelector(".content-composer--readonly.theme-hook"),
  ).not.toBeNull();
  expect(
    container.querySelector(".content-composer__content p br"),
  ).not.toBeNull();
  expect(renderClient).not.toHaveBeenCalled();
});
