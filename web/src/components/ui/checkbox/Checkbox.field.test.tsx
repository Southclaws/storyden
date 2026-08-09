import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { CheckboxField } from "./Checkbox.field";

type FormValues = {
  enabled: boolean;
};

function Harness({ disabled = false }: { disabled?: boolean }) {
  const form = useForm<FormValues>({
    defaultValues: { enabled: false },
  });

  return (
    <>
      <CheckboxField control={form.control} name="enabled" disabled={disabled}>
        Enable notifications
      </CheckboxField>
      <button type="button" onClick={() => form.reset({ enabled: true })}>
        Reset enabled
      </button>
      <output>{String(form.watch("enabled"))}</output>
    </>
  );
}

describe("CheckboxField", () => {
  it("maps Ark checked details into a boolean field value", async () => {
    const user = userEvent.setup();
    render(<Harness />);

    await user.click(
      screen.getByRole("checkbox", { name: "Enable notifications" }),
    );

    expect(screen.getByRole("status")).toHaveTextContent("true");
  });

  it("tracks reset values and propagates disabled state", async () => {
    const user = userEvent.setup();
    const { rerender } = render(<Harness />);

    await user.click(screen.getByRole("button", { name: "Reset enabled" }));
    expect(
      screen.getByRole("checkbox", { name: "Enable notifications" }),
    ).toBeChecked();

    rerender(<Harness disabled />);
    expect(
      screen.getByRole("checkbox", { name: "Enable notifications" }),
    ).toBeDisabled();
  });
});
