import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { describe, expect, it, vi } from "vitest";

import { ModerationWordListEditor } from "./ModerationWordListEditor";

function Harness({ onSubmit }: { onSubmit: () => void }) {
  const [words, setWords] = useState<string[]>([]);

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        onSubmit();
      }}
    >
      <ModerationWordListEditor value={words} onChange={setWords} />
    </form>
  );
}

describe("ModerationWordListEditor", () => {
  it("stages the current word on Enter without submitting the parent form", async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<Harness onSubmit={onSubmit} />);

    const input = screen.getByRole("textbox");
    await user.type(input, "spoiler{Enter}");

    expect(screen.getByText("spoiler")).toBeVisible();
    expect(input).toHaveValue("");
    expect(onSubmit).not.toHaveBeenCalled();
  });
});
