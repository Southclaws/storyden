import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { ComboboxField, type ComboboxFieldItem } from "./Combobox.field";

type FormValues = {
  model?: string;
};

function Harness() {
  const form = useForm<FormValues>({
    defaultValues: { model: "openai/gpt-4-0613" },
  });
  const [items, setItems] = useState<ComboboxFieldItem[]>([]);

  return (
    <>
      <ComboboxField
        ariaLabel="Select model"
        control={form.control}
        items={items}
        name="model"
        placeholder="Select a model"
      />
      <button
        type="button"
        onClick={() =>
          setItems([
            { label: "GPT-4 0613", value: "openai/gpt-4-0613" },
            { label: "GPT-5.6 Luna", value: "openai/gpt-5.6-luna" },
          ])
        }
      >
        Load models
      </button>
      <output>{form.watch("model")}</output>
    </>
  );
}

describe("ComboboxField", () => {
  it("restores a persisted value before asynchronously loaded items arrive", async () => {
    const user = userEvent.setup();
    render(<Harness />);

    const input = screen.getByPlaceholderText("Select a model");
    expect(input).toHaveValue("openai/gpt-4-0613");

    await user.click(screen.getByRole("button", { name: "Load models" }));

    expect(input).toHaveValue("GPT-4 0613");

    await user.click(screen.getByRole("button", { name: "Select model" }));
    await user.click(
      await screen.findByRole("option", { name: "GPT-5.6 Luna" }),
    );

    expect(input).toHaveValue("GPT-5.6 Luna");
    expect(screen.getByRole("status")).toHaveTextContent("openai/gpt-5.6-luna");
  });
});
