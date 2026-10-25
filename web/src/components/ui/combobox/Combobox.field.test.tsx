import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { describe, expect, it, vi } from "vitest";

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
      <textarea aria-label="Instruction" />
    </>
  );
}

describe("ComboboxField", () => {
  it.each(["pointer", "keyboard"])(
    "keeps focus on the next field after %s selection",
    async (selection) => {
      const user = userEvent.setup();
      render(<Harness />);
      await user.click(screen.getByRole("button", { name: "Load models" }));
      const input = screen.getByRole("combobox", { name: "Select model" });
      await user.click(screen.getByRole("button", { name: "Select model" }));
      const option = await screen.findByRole("option", {
        name: "GPT-5.6 Luna",
      });
      await new Promise(requestAnimationFrame);

      const frames = new Map<number, FrameRequestCallback>();
      let frameID = 0;
      vi.spyOn(globalThis, "requestAnimationFrame").mockImplementation(
        (callback) => {
          frames.set(++frameID, callback);
          return frameID;
        },
      );
      vi.spyOn(globalThis, "cancelAnimationFrame").mockImplementation((id) => {
        frames.delete(id);
      });
      if (selection === "pointer") {
        await user.click(option);
      } else {
        await user.keyboard("{ArrowDown}{Enter}");
      }
      expect(screen.getByRole("status")).toHaveTextContent(
        "openai/gpt-5.6-luna",
      );
      expect(frames.size).toBeGreaterThan(0);
      const instruction = screen.getByRole("textbox", { name: "Instruction" });
      await act(async () => {
        instruction.focus();
      });
      act(() => {
        const pending = [...frames.values()];
        frames.clear();
        for (const callback of pending) callback(performance.now());
      });
      expect(instruction).toHaveFocus();
      await user.keyboard("block trail action");
      expect(instruction).toHaveValue("block trail action");
      expect(input).toHaveValue("GPT-5.6 Luna");
    },
  );

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

    await waitFor(() => expect(input).toHaveFocus());
    expect(input).toHaveValue("GPT-5.6 Luna");
    expect(screen.getByRole("status")).toHaveTextContent("openai/gpt-5.6-luna");
  });
});
