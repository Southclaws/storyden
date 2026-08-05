import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { RadioGroupField } from "./RadioGroup.field";

type FormValues = {
  visibility: string;
};

function Harness() {
  const form = useForm<FormValues>({
    defaultValues: { visibility: "public" },
  });

  return (
    <>
      <RadioGroupField
        ariaLabel="Visibility"
        control={form.control}
        name="visibility"
        items={[
          { label: "Public", value: "public" },
          { label: "Private", value: "private" },
        ]}
      />
      <button
        type="button"
        onClick={() => form.reset({ visibility: "private" })}
      >
        Reset private
      </button>
      <output>{form.watch("visibility")}</output>
    </>
  );
}

describe("RadioGroupField", () => {
  it("uses the current form value after interaction and reset", async () => {
    const user = userEvent.setup();
    render(<Harness />);

    expect(
      screen.getByRole("radiogroup", { name: "Visibility" }),
    ).toBeVisible();

    await user.click(screen.getByRole("radio", { name: "Private" }));
    expect(screen.getByRole("status")).toHaveTextContent("private");

    await user.click(screen.getByRole("radio", { name: "Public" }));
    await user.click(screen.getByRole("button", { name: "Reset private" }));
    expect(screen.getByRole("radio", { name: "Private" })).toBeChecked();
  });
});
