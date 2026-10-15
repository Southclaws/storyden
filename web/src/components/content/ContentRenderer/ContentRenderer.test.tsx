import { fireEvent, render, screen } from "@testing-library/react";
import { renderToString } from "react-dom/server";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ContentLinkCard } from "./ContentLinkCard";
import { ContentRenderer } from "./ContentRenderer";
import { ServerRenderedContent } from "./ServerRenderedContent";

vi.mock("server-only", () => ({}));

function ReadingFixture({ html }: { html: string }) {
  return (
    <ServerRenderedContent values={[html]}>
      <ContentRenderer initialValue={html} />
    </ServerRenderedContent>
  );
}

const trigger = vi.hoisted(() => vi.fn());
const linkState = vi.hoisted(() => ({ error: undefined as Error | undefined }));

vi.mock("@/api/openapi-client/links", () => ({
  useLinkCreate: () => ({
    trigger,
    data: undefined,
    error: linkState.error,
    isMutating: false,
  }),
}));

describe("ContentRenderer", () => {
  beforeEach(() => {
    linkState.error = new Error("Unavailable");
  });

  it("server-renders a link-card destination before its preview has loaded", () => {
    linkState.error = undefined;
    const html = renderToString(
      <ReadingFixture html='<div data-display="card" data-href="https://example.com/card"></div>' />,
    );
    const document = new DOMParser().parseFromString(html, "text/html");
    expect(document.querySelector("a")?.getAttribute("href")).toBe(
      "https://example.com/card",
    );
    expect(document.querySelector("a")?.textContent).toContain(
      "https://example.com/card",
    );
    expect(trigger).not.toHaveBeenCalled();
  });

  it("does not expose an unsafe link-card destination while loading its preview", () => {
    linkState.error = undefined;
    const html = renderToString(
      <ReadingFixture html='<div data-display="card" data-href="javascript:alert(1)"></div>' />,
    );
    const document = new DOMParser().parseFromString(html, "text/html");
    expect(document.querySelector("a")).toBeNull();
  });

  it("renders stored TipTap HTML with the same document structure on the server and client", () => {
    const initialValue =
      '<body><h2>Field notes</h2><p><strong>Bold</strong> and <em>italic</em> with <s>struck</s> and <u>underlined</u> text, <code>inline code</code> and <a href="https://example.com/guide">a guide</a>.</p><ul><li><p>First</p></li></ul><ol start="3"><li><p>Third</p></li></ol><blockquote><p>A quotation</p></blockquote><pre><code class="language-js">const value = 1;</code></pre><hr><p>Line one<br>Line two</p><img src="https://example.com/photo.png" alt="Field photograph"></body>';
    const html = renderToString(<ReadingFixture html={initialValue} />);
    expect(html).toContain("Field notes");
    expect(html).toContain('<ol start="3">');
    expect(html).not.toContain("contenteditable");

    const { container } = render(<ReadingFixture html={initialValue} />);
    expect(
      screen.getByRole("heading", { name: "Field notes", level: 2 }),
    ).toBeVisible();
    expect(container.querySelector("strong")).toHaveTextContent("Bold");
    expect(container.querySelector("em")).toHaveTextContent("italic");
    expect(container.querySelector("s")).toHaveTextContent("struck");
    expect(container.querySelector("u")).toHaveTextContent("underlined");
    expect(container.querySelector("ol")).toHaveAttribute("start", "3");
    expect(container.querySelector("blockquote")).toHaveTextContent(
      "A quotation",
    );
    expect(container.querySelector("pre code")).toHaveClass("language-js");
    expect(screen.getByRole("link", { name: "a guide" })).toHaveAttribute(
      "href",
      "https://example.com/guide",
    );
    expect(
      screen
        .getByRole("img", { name: "Field photograph" })
        .closest("[data-node-view-wrapper]"),
    ).not.toBeNull();
    expect(container.querySelector("[contenteditable]")).toBeNull();
    expect(screen.queryByRole("textbox")).not.toBeInTheDocument();
    expect(
      container.querySelector(".content-composer__content"),
    ).not.toBeNull();
    expect(container.querySelector(".content-composer__editor")).toBeNull();
  });

  it("updates read-only content including an explicitly empty value", async () => {
    const { rerender } = render(
      <ContentRenderer initialValue="<p>Original</p>" />,
    );
    rerender(
      <ContentRenderer
        initialValue="<p>Original</p>"
        value="<h3>Updated</h3>"
      />,
    );
    expect(
      await screen.findByRole("heading", { name: "Updated" }),
    ).toBeVisible();
    expect(screen.queryByText("Original")).not.toBeInTheDocument();
    rerender(<ContentRenderer initialValue="<p>Original</p>" value="" />);
    expect(screen.queryByText("Updated")).not.toBeInTheDocument();
    expect(screen.queryByText("Original")).not.toBeInTheDocument();
  });

  it("preserves link-card fallback navigation without an editor", () => {
    render(
      <ReadingFixture html='<div data-display="card" data-href="https://example.com/card"></div>' />,
    );
    expect(
      screen.getByRole("link", { name: "https://example.com/card" }),
    ).toHaveAttribute("href", "https://example.com/card");
    expect(screen.getByText("Link preview failed to load")).toBeVisible();
    expect(trigger).toHaveBeenCalledWith({ url: "https://example.com/card" });
  });

  it("preserves the editable link-card error and retry action", () => {
    render(<ContentLinkCard href="https://example.com/card" isEditable />);
    expect(screen.getByRole("alert")).toHaveTextContent(
      "Link preview failed: Unavailable",
    );
    fireEvent.click(screen.getByRole("button", { name: "Retry" }));
    expect(trigger).toHaveBeenCalledTimes(2);
    expect(trigger).toHaveBeenLastCalledWith({
      url: "https://example.com/card",
    });
  });

  it("keeps TipTap schema filtering of executable markup and unsupported attributes", () => {
    const { container } = render(
      <ReadingFixture html='<p onclick="alert(1)">Safe<script>alert(1)</script><a href="javascript:alert(1)">unsafe link</a></p><img src="https://example.com/photo.png" onerror="alert(1)">' />,
    );
    expect(container.querySelector("script,[onclick],[onerror]")).toBeNull();
    expect(container.querySelector('a[href^="javascript:"]')).toBeNull();
    expect(container.querySelector("p")).toHaveTextContent("Safe");
  });
});
