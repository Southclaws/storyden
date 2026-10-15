import { render, screen } from "@testing-library/react";
import { renderToString } from "react-dom/server";
import { describe, expect, it, vi } from "vitest";

import { RenderedContentProvider } from "../ContentRenderer/RenderedContentContext";

import { ContentComposer } from "./ContentComposer";

const editor = vi.hoisted(() => ({ loaded: vi.fn(), rendered: vi.fn() }));

vi.mock("./ContentEditor", () => {
  editor.loaded();
  return {
    ContentEditor: ({ initialValue }: { initialValue?: string }) => {
      editor.rendered();
      return <textarea aria-label="Editor" defaultValue={initialValue} />;
    },
  };
});

describe("ContentComposer boundary", () => {
  it("server-renders content even when editing is enabled, without loading an editor", () => {
    const html = renderToString(
      <RenderedContentProvider
        content={{ "<h2>Server content</h2>": <h2>Server content</h2> }}
      >
        <ContentComposer initialValue="<h2>Server content</h2>" />
      </RenderedContentProvider>,
    );
    expect(html).toContain("<h2>Server content</h2>");
    expect(html).not.toContain("contenteditable");
    expect(editor.loaded).not.toHaveBeenCalled();
    expect(editor.rendered).not.toHaveBeenCalled();
  });

  it("renders content without loading the editor and loads it when editing starts", async () => {
    const { rerender } = render(
      <ContentComposer disabled initialValue="<h2>Original</h2>" />,
      {
        wrapper: ({ children }) => (
          <RenderedContentProvider
            content={{ "<h2>Original</h2>": <h2>Original</h2> }}
          >
            {children}
          </RenderedContentProvider>
        ),
      },
    );
    expect(screen.getByRole("heading", { name: "Original" })).toBeVisible();
    expect(editor.loaded).not.toHaveBeenCalled();
    expect(editor.rendered).not.toHaveBeenCalled();

    rerender(<ContentComposer initialValue="<h2>Original</h2>" />);
    expect(await screen.findByRole("textbox", { name: "Editor" })).toHaveValue(
      "<h2>Original</h2>",
    );
    expect(editor.loaded).toHaveBeenCalledOnce();

    rerender(<ContentComposer disabled initialValue="<h2>Saved</h2>" />);
    expect(await screen.findByRole("heading", { name: "Saved" })).toBeVisible();
    expect(screen.queryByRole("textbox")).not.toBeInTheDocument();
  });
});
