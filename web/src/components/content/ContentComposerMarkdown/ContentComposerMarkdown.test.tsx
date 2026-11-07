import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { ContentComposerMarkdown } from "./ContentComposerMarkdown";

describe("ContentComposerMarkdown", () => {
  it("starts a fresh buffer when keyed, while preserving spaces in the current draft", () => {
    const { rerender } = render(
      <ContentComposerMarkdown key="first" initialValue="<p>First draft</p>" />,
    );
    const textbox = screen.getByRole("textbox");
    fireEvent.change(textbox, { target: { value: "hello " } });

    rerender(
      <ContentComposerMarkdown key="first" initialValue="<p>hello </p>" />,
    );
    expect(screen.getByRole("textbox")).toBe(textbox);
    expect(textbox).toHaveValue("hello ");

    rerender(<ContentComposerMarkdown key="second" initialValue="" />);
    expect(screen.getByRole("textbox")).not.toBe(textbox);
    expect(screen.getByRole("textbox")).toHaveValue("");
  });

  it.each(["markdown", "html"] as const)(
    "renders live read-only %s props, including cleared content",
    (format) => {
      const onChange = vi.fn();
      const { container, rerender } = render(
        <ContentComposerMarkdown
          disabled
          initialValue=""
          initialValueFormat={format}
          onChange={onChange}
        />,
      );
      const renderer = container.querySelector(".markdown-editor-readonly");

      for (const text of [
        "First chunk",
        "First chunk and more text",
        "Final response",
        "",
      ]) {
        let content = text;
        if (format === "html") {
          content = `<p>${text}</p>`;
        }

        rerender(
          <ContentComposerMarkdown
            disabled
            initialValue={content}
            initialValueFormat={format}
            onChange={onChange}
          />,
        );
        expect(container.querySelector(".markdown-editor-readonly")).toBe(
          renderer,
        );
        expect(renderer?.textContent).toBe(text);
      }
      expect(onChange).not.toHaveBeenCalled();
    },
  );

  it("updates the read-only value without replacing an editable draft", () => {
    const { container, rerender } = render(
      <ContentComposerMarkdown
        initialValue="Original"
        initialValueFormat="markdown"
      />,
    );
    fireEvent.change(screen.getByRole("textbox"), {
      target: { value: "My draft " },
    });

    rerender(
      <ContentComposerMarkdown
        disabled
        initialValue="Original"
        value="Latest saved content"
        initialValueFormat="markdown"
      />,
    );
    expect(
      container.querySelector(".markdown-editor-readonly"),
    ).toHaveTextContent("Latest saved content");

    rerender(
      <ContentComposerMarkdown
        initialValue="Updated initial content"
        value="Ignored while editing"
        initialValueFormat="markdown"
      />,
    );
    expect(screen.getByRole("textbox")).toHaveValue("My draft ");
  });

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
