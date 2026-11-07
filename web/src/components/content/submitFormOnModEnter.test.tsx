import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { ContentComposerMarkdown } from "./ContentComposerMarkdown/ContentComposerMarkdown";
import { ContentComposerRich } from "./ContentComposerRich/ContentComposerRich";

describe.each([
  { mode: "rich text", Composer: ContentComposerRich },
  { mode: "Markdown", Composer: ContentComposerMarkdown },
])("$mode form submission", ({ Composer }) => {
  it.each([{ metaKey: true }, { ctrlKey: true }])(
    "submits once with %o through the form's submit button without changing content",
    async (modifier) => {
      const onSubmit = vi.fn((event) => event.preventDefault());
      const onChange = vi.fn();
      render(
        <form onSubmit={onSubmit}>
          <button type="button">Cancel</button>
          <Composer initialValue="<p>Example content</p>" onChange={onChange} />
          <button type="submit" name="action" value="publish">
            Publish
          </button>
        </form>,
      );

      const editor = await screen.findByRole("textbox");
      fireEvent.keyDown(editor, { key: "Enter", ...modifier });

      expect(onSubmit).toHaveBeenCalledOnce();
      expect(onSubmit.mock.calls[0]![0].nativeEvent.submitter).toBe(
        screen.getByRole("button", { name: "Publish" }),
      );
      expect(onChange).not.toHaveBeenCalled();
    },
  );

  it("respects a disabled submit button and ignores held shortcut keys", async () => {
    const onSubmit = vi.fn((event) => event.preventDefault());
    render(
      <form onSubmit={onSubmit}>
        <Composer initialValue="<p>Example content</p>" />
        <button type="submit" disabled>
          Post
        </button>
      </form>,
    );

    const editor = await screen.findByRole("textbox");
    fireEvent.keyDown(editor, { key: "Enter", metaKey: true });
    expect(onSubmit).not.toHaveBeenCalled();

    screen.getByRole("button", { name: "Post" }).removeAttribute("disabled");
    fireEvent.keyDown(editor, { key: "Enter", metaKey: true, repeat: true });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("runs native form validation", async () => {
    const onSubmit = vi.fn((event) => event.preventDefault());
    render(
      <form onSubmit={onSubmit}>
        <input aria-label="Title" required />
        <Composer initialValue="<p>Example content</p>" />
        <button type="submit">Publish</button>
      </form>,
    );

    const editor = await screen.findByRole("textbox", { name: /editor/i });
    fireEvent.keyDown(editor, { key: "Enter", ctrlKey: true });
    expect(onSubmit).not.toHaveBeenCalled();

    fireEvent.change(screen.getByRole("textbox", { name: "Title" }), {
      target: { value: "A title" },
    });
    fireEvent.keyDown(editor, { key: "Enter", ctrlKey: true });
    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it("does not submit for newlines, other modifiers or IME composition", async () => {
    const onSubmit = vi.fn((event) => event.preventDefault());
    render(
      <form onSubmit={onSubmit}>
        <Composer initialValue="<p>Example content</p>" />
        <button type="submit">Post</button>
      </form>,
    );

    const editor = await screen.findByRole("textbox");
    const ignoredKeys = [
      {},
      { shiftKey: true },
      { metaKey: true, shiftKey: true },
      { ctrlKey: true, altKey: true },
      { metaKey: true, isComposing: true },
      { metaKey: true, keyCode: 229 },
    ];

    for (const modifiers of ignoredKeys) {
      fireEvent.keyDown(editor, { key: "Enter", ...modifiers });
    }

    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("leaves shortcuts alone when there is no form submit action", async () => {
    const { rerender } = render(<Composer />);
    const standaloneEditor = await screen.findByRole("textbox");
    const onEditorKeyDown = vi.fn();
    standaloneEditor.addEventListener("keydown", onEditorKeyDown);
    fireEvent.keyDown(standaloneEditor, { key: "Enter", ctrlKey: true });
    expect(onEditorKeyDown).toHaveBeenCalledOnce();

    const onSubmit = vi.fn((event) => event.preventDefault());
    rerender(
      <form onSubmit={onSubmit}>
        <Composer />
      </form>,
    );
    const editor = await screen.findByRole("textbox");
    const onFormEditorKeyDown = vi.fn();
    editor.addEventListener("keydown", onFormEditorKeyDown);
    fireEvent.keyDown(editor, { key: "Enter", ctrlKey: true });
    expect(onFormEditorKeyDown).toHaveBeenCalledOnce();
    expect(onSubmit).not.toHaveBeenCalled();
  });
});
