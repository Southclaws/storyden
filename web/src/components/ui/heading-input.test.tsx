import { createEvent, fireEvent, render } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { HeadingInput } from "./heading-input";

describe("HeadingInput", () => {
  it("emits text updates on input", () => {
    const onValueChange = vi.fn();

    const { getByText } = render(
      <HeadingInput onValueChange={onValueChange} defaultValue="Draft title" />,
    );

    const field = getByText("Draft title");
    fireEvent.input(field, { target: { textContent: "Updated title" } });

    expect(onValueChange).toHaveBeenCalledWith("Updated title");
  });

  it("prevents new lines from enter key presses", () => {
    const { getByText } = render(
      <HeadingInput onValueChange={() => {}} defaultValue="Draft title" />,
    );
    const field = getByText("Draft title");

    const event = createEvent.keyDown(field, { key: "Enter", code: "Enter" });
    fireEvent(field, event);

    expect(event.defaultPrevented).toBe(true);
  });

  it("sanitizes pasted text by removing line breaks", () => {
    const execCommand = vi
      .spyOn(document as Document & { execCommand: typeof document.execCommand }, "execCommand")
      .mockImplementation(() => true);

    const { getByText } = render(
      <HeadingInput onValueChange={() => {}} defaultValue="Draft title" />,
    );
    const field = getByText("Draft title");

    const paste = createEvent.paste(field);
    Object.defineProperty(paste, "clipboardData", {
      value: {
        getData: () => "hello\nworld\r\nagain",
      },
    });

    fireEvent(field, paste);

    expect(execCommand).toHaveBeenCalledWith(
      "insertText",
      false,
      "hello world again",
    );
  });

  it("syncs content with controlled value updates", () => {
    const { container, rerender } = render(
      <HeadingInput onValueChange={() => {}} value="First title" />,
    );

    const field = container.querySelector(
      '[contenteditable="true"]',
    ) as HTMLSpanElement;

    expect(field.textContent).toBe("First title");

    rerender(<HeadingInput onValueChange={() => {}} value="Second title" />);

    expect(field.textContent).toBe("Second title");
  });
});
