import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { describe, expect, it } from "vitest";

import { EditableBadge } from "./EditableBadge";

function Editor() {
  const [value, setValue] = useState("Commercial");
  return (
    <EditableBadge
      editing
      value={value}
      inputLabel="License"
      placeholder="License"
      onValueChange={setValue}
      size="sm"
    />
  );
}

describe("EditableBadge", () => {
  it("edits the value and finishes with Enter without introducing a newline", async () => {
    const user = userEvent.setup();
    render(<Editor />);
    const editor = screen.getByRole("textbox", { name: "License" });
    await user.clear(editor);
    await user.type(editor, "Open source");
    expect(editor).toHaveValue("Open source");
    await user.keyboard("{Enter}");
    expect(editor).not.toHaveFocus();
    expect(editor).toHaveValue("Open source");
  });

  it("keeps empty values editable and pastes plain text without newlines", async () => {
    const user = userEvent.setup();
    render(<Editor />);
    const editor = screen.getByRole("textbox", { name: "License" });
    await user.clear(editor);
    expect(editor).toHaveAttribute("placeholder", "License");
    await user.paste("Open\nsource");
    expect(editor).toHaveValue("Open source");
    (editor as HTMLInputElement).setSelectionRange(0, 4);
    await user.paste("Community\r\nlicensed");
    expect(editor).toHaveValue("Community licensed source");
  });

  it("does not finish editing while composing text", async () => {
    const user = userEvent.setup();
    render(<Editor />);
    const editor = screen.getByRole("textbox", { name: "License" });
    await user.click(editor);
    fireEvent.keyDown(editor, { key: "Enter", isComposing: true });
    expect(editor).toHaveFocus();
  });

  it("renders view content without exposing a hidden editor", () => {
    render(
      <EditableBadge
        editing={false}
        value="Commercial"
        inputLabel="License"
        onValueChange={() => {}}
      >
        <span>License: Commercial</span>
      </EditableBadge>,
    );
    expect(screen.getByText("License: Commercial")).toBeVisible();
    expect(screen.queryByRole("textbox")).toBeNull();
  });
});
