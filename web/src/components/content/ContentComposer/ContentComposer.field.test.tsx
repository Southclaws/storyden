import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it, vi } from "vitest";

import { ContentComposerField } from "./ContentComposer.field";

vi.mock("./ContentComposer", () => ({
  ContentComposer: ({
    disabled,
    initialValue,
    onChange,
  }: {
    disabled?: boolean;
    initialValue?: string;
    onChange?: (value: string, isEmpty: boolean) => void;
  }) => (
    <button
      type="button"
      disabled={disabled}
      onClick={() => onChange?.("Updated", false)}
    >
      {initialValue}
    </button>
  ),
}));

type FormValues = {
  body: string;
};

function Harness({ onEmptyStateChange }: { onEmptyStateChange: () => void }) {
  const form = useForm<FormValues>({ defaultValues: { body: "Initial" } });

  return (
    <>
      <ContentComposerField
        control={form.control}
        name="body"
        onEmptyStateChange={onEmptyStateChange}
      />
      <output>{form.watch("body")}</output>
    </>
  );
}

describe("ContentComposerField", () => {
  it("uses the field value as its initial content and reports changes", async () => {
    const user = userEvent.setup();
    const onEmptyStateChange = vi.fn();
    render(<Harness onEmptyStateChange={onEmptyStateChange} />);

    await user.click(screen.getByRole("button", { name: "Initial" }));

    expect(screen.getByRole("status")).toHaveTextContent("Updated");
    expect(onEmptyStateChange).toHaveBeenCalledWith(false);
  });
});
