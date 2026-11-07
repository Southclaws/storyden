import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { type Editor } from "@tiptap/core";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { describe, expect, it, vi } from "vitest";

import { ContentComposerRich } from "../ContentComposerRich/ContentComposerRich";

import { ContentComposerField } from "./ContentComposer.field";

vi.mock("./ContentEditor", () => ({
  ContentEditor: ContentComposerRich,
}));

type FormValues = {
  body: string;
};

function Harness({ onEmptyStateChange }: { onEmptyStateChange: () => void }) {
  const form = useForm<FormValues>({
    defaultValues: { body: "<p>Initial</p>" },
  });
  const [draftKey, setDraftKey] = useState(0);

  return (
    <>
      <ContentComposerField
        key={draftKey}
        control={form.control}
        name="body"
        onEmptyStateChange={onEmptyStateChange}
      />
      <button
        type="button"
        onClick={() => {
          form.reset({ body: "" });
          setDraftKey((previous) => previous + 1);
        }}
      >
        New draft
      </button>
      <output>{form.watch("body")}</output>
    </>
  );
}

describe("ContentComposerField", () => {
  it("uses the field value as its initial content and reports changes", async () => {
    const onEmptyStateChange = vi.fn();
    render(<Harness onEmptyStateChange={onEmptyStateChange} />);
    const textbox = await screen.findByRole("textbox", {
      name: "Rich text editor",
    });
    const editor = (textbox as HTMLElement & { editor: Editor }).editor;
    expect(editor.getHTML()).toBe("<p>Initial</p>");

    act(() => {
      editor.commands.insertContentAt(8, " updated");
    });

    expect(screen.getByRole("status").textContent).toBe(
      "<p>Initial updated</p>",
    );
    expect(onEmptyStateChange).toHaveBeenCalledWith(false);
  });

  it("preserves trailing spaces and selection through form updates after starting fresh drafts", async () => {
    const user = userEvent.setup();
    render(<Harness onEmptyStateChange={vi.fn()} />);
    const original = await screen.findByRole("textbox", {
      name: "Rich text editor",
    });
    let previousEditor = (original as HTMLElement & { editor: Editor }).editor;

    for (let draft = 0; draft < 2; draft += 1) {
      await user.click(screen.getByRole("button", { name: "New draft" }));
      const textbox = await screen.findByRole("textbox", {
        name: "Rich text editor",
      });
      const editor = (textbox as HTMLElement & { editor: Editor }).editor;
      expect(editor).not.toBe(previousEditor);
      expect(editor.can().undo()).toBe(false);
      expect(editor.getText()).toBe("");
      previousEditor = editor;

      act(() => {
        editor.commands.insertContent("hello");
      });
      act(() => {
        editor.commands.insertContent(" ");
      });
      expect(editor.getText()).toBe("hello ");
      expect(screen.getByRole("status").textContent).toBe("<p>hello </p>");

      act(() => {
        editor.commands.setTextSelection(3);
        editor.commands.insertContent("X");
      });
      expect(editor.getText()).toBe("heXllo ");
      expect(editor.state.selection.from).toBe(4);
    }
  });
});
