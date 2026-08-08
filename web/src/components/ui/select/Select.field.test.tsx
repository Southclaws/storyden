import { createListCollection } from "@ark-ui/react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm } from "react-hook-form";
import { describe, expect, it } from "vitest";

import { SelectField } from "./Select.field";

const collection = createListCollection({
  items: [
    { label: "Discussion", value: "discussion" },
    { label: "Library", value: "library" },
  ],
});

type FormValues = {
  section?: string;
};

function Harness() {
  const form = useForm<FormValues>({
    defaultValues: { section: "discussion" },
  });

  return (
    <>
      <SelectField
        control={form.control}
        name="section"
        collection={collection}
        placeholder="Choose section"
      />
      <button type="button" onClick={() => form.reset({ section: "library" })}>
        Reset library
      </button>
      <output>{form.watch("section")}</output>
    </>
  );
}

describe("SelectField", () => {
  it("keeps the hidden select and form state synchronized after reset", async () => {
    const user = userEvent.setup();
    render(<Harness />);
    expect(
      screen.getByRole("combobox", { name: "Choose section" }),
    ).toBeInTheDocument();
    const select = document.querySelector(
      'select[name="section"]',
    ) as HTMLSelectElement;

    await waitFor(() => {
      expect(select).toHaveValue("discussion");
    });
    await user.click(screen.getByRole("button", { name: "Reset library" }));

    await waitFor(() => {
      expect(select).toHaveValue("library");
    });
    expect(screen.getByRole("status")).toHaveTextContent("library");
  });
});
