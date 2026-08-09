import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { NumberInputField } from "./NumberInput.field";

type FormValues = {
  count?: number;
};

function Harness() {
  const form = useForm<FormValues>({
    defaultValues: { count: 12 },
  });

  return (
    <>
      <NumberInputField control={form.control} name="count" ariaLabel="Count" />
      <button type="button" onClick={() => form.reset({ count: 7 })}>
        Reset count
      </button>
      <output>{form.watch("count")}</output>
    </>
  );
}

describe("NumberInputField", () => {
  it("maps numeric input and reset values without losing its accessible name", async () => {
    const user = userEvent.setup();
    render(<Harness />);

    const input = screen.getByRole("spinbutton", { name: "Count" });
    await user.clear(input);
    await user.type(input, "24");

    await waitFor(() => {
      expect(screen.getByRole("status")).toHaveTextContent("24");
    });

    await user.click(screen.getByRole("button", { name: "Reset count" }));

    await waitFor(() => {
      expect(input).toHaveValue("7");
    });
  });
});
