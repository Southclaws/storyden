import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { DatePickerInputField } from "./DatePickerField";

type FormValues = {
  date?: string;
};

function Harness() {
  const form = useForm<FormValues>({
    defaultValues: {
      date: "2026-03-03T00:00:00.000Z",
    },
  });

  const value = form.watch("date");

  return (
    <>
      <DatePickerInputField<FormValues> name="date" control={form.control} />
      <button type="button" onClick={() => form.setValue("date", "")}>
        Clear
      </button>
      <output data-testid="form-value">{value ?? "undefined"}</output>
    </>
  );
}

describe("DatePickerInputField", () => {
  it("clears the visible input when react-hook-form state is cleared", async () => {
    const user = userEvent.setup();

    render(<Harness />);

    const input = screen.getByPlaceholderText("YYYY-MM-DD");
    expect(input).toHaveValue("2026-03-03");

    await user.click(screen.getByRole("button", { name: "Clear" }));

    await waitFor(() => {
      expect(input).toHaveValue("");
      expect(screen.getByTestId("form-value")).toHaveTextContent("");
    });
  });
});
