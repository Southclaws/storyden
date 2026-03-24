import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Controller, useForm } from "react-hook-form";
import { describe, expect, it, vi } from "vitest";

import { Switch } from "./switch";

describe("Switch", () => {
  it("emits checked state changes when toggled", async () => {
    const user = userEvent.setup();
    const onCheckedChange = vi.fn();

    render(<Switch onCheckedChange={onCheckedChange}>Hidden category</Switch>);

    const input = screen.getByRole("checkbox", { name: "Hidden category" });

    await user.click(input);
    expect(onCheckedChange).toHaveBeenNthCalledWith(
      1,
      expect.objectContaining({ checked: true }),
    );

    await user.click(input);
    expect(onCheckedChange).toHaveBeenNthCalledWith(
      2,
      expect.objectContaining({ checked: false }),
    );
  });

  it("works with react-hook-form controlled state", async () => {
    const user = userEvent.setup();

    function Harness() {
      const form = useForm<{ is_hidden: boolean }>({
        defaultValues: { is_hidden: false },
      });

      return (
        <Controller
          control={form.control}
          name="is_hidden"
          render={({ field }) => (
            <Switch
              checked={field.value}
              onCheckedChange={({ checked }) => field.onChange(checked)}
            >
              Hidden category
            </Switch>
          )}
        />
      );
    }

    render(<Harness />);

    const input = screen.getByRole("checkbox", { name: "Hidden category" });
    expect(input).not.toBeChecked();

    await user.click(input);
    expect(input).toBeChecked();
  });
});
