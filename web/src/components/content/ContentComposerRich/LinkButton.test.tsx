import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Editor } from "@tiptap/react";
import { describe, expect, it, vi } from "vitest";

import { LinkButton } from "./LinkButton";

function createEditorState(options: { active: boolean; href?: string }) {
  const chain = {
    focus: vi.fn(() => chain),
    extendMarkRange: vi.fn(() => chain),
    unsetLink: vi.fn(() => chain),
    setLink: vi.fn(() => chain),
    run: vi.fn(() => true),
  };

  const editor = {
    isActive: vi.fn(() => options.active),
    getAttributes: vi.fn(() => ({ href: options.href ?? "" })),
    chain: vi.fn(() => chain),
  } as unknown as Editor;

  return { editor, chain };
}

describe("LinkButton", () => {
  it("adds a normalized link for inactive selections", async () => {
    const user = userEvent.setup();
    const { editor, chain } = createEditorState({ active: false });

    render(<LinkButton editor={editor} />);

    await user.click(screen.getByTitle("Add link"));
    await user.type(screen.getByLabelText("Link URL"), "example.com{enter}");

    expect(chain.setLink).toHaveBeenCalledWith({ href: "https://example.com/" });
    expect(chain.extendMarkRange).not.toHaveBeenCalled();
    expect(chain.run).toHaveBeenCalled();
  });

  it("removes an existing link when the input is emptied", async () => {
    const user = userEvent.setup();
    const { editor, chain } = createEditorState({
      active: true,
      href: "https://storyden.org",
    });

    render(<LinkButton editor={editor} />);

    await user.click(screen.getByTitle("Edit link"));

    const input = screen.getByLabelText("Link URL");
    await user.clear(input);
    await user.click(screen.getByRole("button", { name: "Update" }));

    expect(chain.unsetLink).toHaveBeenCalled();
    expect(chain.run).toHaveBeenCalled();
  });

  it("keeps the popover open for invalid link values", async () => {
    const user = userEvent.setup();
    const { editor, chain } = createEditorState({ active: false });

    render(<LinkButton editor={editor} />);

    await user.click(screen.getByTitle("Add link"));
    await user.type(
      screen.getByLabelText("Link URL"),
      "javascript:alert(1).com",
    );
    await user.click(screen.getByRole("button", { name: "Add" }));

    expect(chain.setLink).not.toHaveBeenCalled();
    expect(screen.getByLabelText("Link URL")).toBeInTheDocument();
  });

  it("shows remove button for active links", async () => {
    const user = userEvent.setup();
    const { editor, chain } = createEditorState({
      active: true,
      href: "https://storyden.org",
    });

    render(<LinkButton editor={editor} />);

    await user.click(screen.getByTitle("Edit link"));
    await user.click(screen.getByTitle("Remove link"));

    expect(chain.unsetLink).toHaveBeenCalled();
  });
});
