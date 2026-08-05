import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { HeadingInputField } from "./HeadingInput.field";

type FormValues = {
  title: string;
};

function Harness() {
  const form = useForm<FormValues>({ defaultValues: { title: "Original" } });

  return (
    <>
      <HeadingInputField
        control={form.control}
        name="title"
        aria-label="Title"
      />
      <button type="button" onClick={() => form.reset({ title: "Reset" })}>
        Reset title
      </button>
      <output>{form.watch("title")}</output>
    </>
  );
}

describe("HeadingInputField", () => {
  it("updates form state and follows reset values", async () => {
    const user = userEvent.setup();
    render(<Harness />);

    const input = screen.getByRole("textbox", { name: "Title" });
    input.textContent = "Changed";
    fireEvent.input(input);
    expect(screen.getByRole("status")).toHaveTextContent("Changed");

    await user.click(screen.getByRole("button", { name: "Reset title" }));
    expect(input).toHaveTextContent("Reset");
  });
});
