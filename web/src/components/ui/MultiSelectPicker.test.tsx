import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { deriveColour } from "@/utils/colour";

import { MultiSelectPicker, MultiSelectPickerItem } from "./MultiSelectPicker";

function openPicker() {
  return userEvent.click(screen.getByRole("button", { name: /Select items/i }));
}

function getSearchInput() {
  return screen.getByRole("textbox", { name: "Search for items" });
}

function rect({
  left,
  right,
  width,
}: {
  left: number;
  right: number;
  width: number;
}): DOMRect {
  return {
    x: left,
    y: 0,
    top: 0,
    bottom: 20,
    left,
    right,
    width,
    height: 20,
    toJSON: () => ({}),
  } as DOMRect;
}

describe("MultiSelectPicker", () => {
  const baseValue: MultiSelectPickerItem[] = [{ label: "Alpha", value: "alpha" }];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calls onQuery as the search input changes", async () => {
    const user = userEvent.setup();
    const onQuery = vi.fn();

    render(
      <MultiSelectPicker
        value={[]}
        onQuery={onQuery}
        onChange={vi.fn().mockResolvedValue(undefined)}
      />,
    );

    await openPicker();
    await user.type(getSearchInput(), "ab");

    expect(onQuery).toHaveBeenCalled();
    expect(onQuery.mock.calls.at(-1)?.[0]).toBe("ab");
  });

  it("adds a selected result and auto-assigns colour when enabled", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);
    const result = { label: "Beta", value: "beta" };

    render(
      <MultiSelectPicker
        value={baseValue}
        onQuery={vi.fn()}
        onChange={onChange}
        queryResults={[result]}
        autoColour
      />,
    );

    await openPicker();
    await user.click(screen.getByRole("menuitem", { name: "Beta" }));

    await waitFor(() =>
      expect(onChange).toHaveBeenCalledWith([
        ...baseValue,
        { ...result, colour: deriveColour("beta") },
      ]),
    );
  });

  it("creates a new item from Enter and clears the query input", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);

    render(
      <MultiSelectPicker
        value={[]}
        onQuery={vi.fn()}
        onChange={onChange}
        allowNewValues
      />,
    );

    await openPicker();
    const input = getSearchInput();
    await user.type(input, "Gamma{enter}");

    await waitFor(() =>
      expect(onChange).toHaveBeenCalledWith([
        { label: "Gamma", value: "Gamma", colour: undefined },
      ]),
    );
    await waitFor(() => expect(input).toHaveValue(""));
  });

  it("removes a selected item from the selected section", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);

    render(
      <MultiSelectPicker
        value={baseValue}
        onQuery={vi.fn()}
        onChange={onChange}
      />,
    );

    await openPicker();
    await user.click(screen.getByRole("button", { name: "Remove Alpha" }));

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([]));
  });

  it("shows no-results state when searching with no matches", async () => {
    const user = userEvent.setup();

    render(
      <MultiSelectPicker
        value={[]}
        onQuery={vi.fn()}
        onChange={vi.fn().mockResolvedValue(undefined)}
        queryResults={[]}
      />,
    );

    await openPicker();
    await user.type(getSearchInput(), "zzz");

    expect(screen.getByText("No results found")).toBeInTheDocument();
  });

  it("shows hidden badge count when badges overflow the trigger", async () => {
    const value = [
      { label: "One", value: "one" },
      { label: "Two", value: "two" },
      { label: "Three", value: "three" },
    ];

    const { rerender } = render(
      <MultiSelectPicker
        value={value}
        onQuery={vi.fn()}
        onChange={vi.fn().mockResolvedValue(undefined)}
      />,
    );

    const trigger = screen.getByRole("button", {
      name: /Select items, 3 selected/i,
    });
    const container = trigger.querySelector("div");
    const one = screen.getByText("One");
    const two = screen.getByText("Two");
    const three = screen.getByText("Three");

    expect(container).not.toBeNull();

    Object.defineProperty(container!, "getBoundingClientRect", {
      configurable: true,
      value: () => rect({ left: 0, right: 100, width: 100 }),
    });
    Object.defineProperty(one, "getBoundingClientRect", {
      configurable: true,
      value: () => rect({ left: 0, right: 40, width: 40 }),
    });
    Object.defineProperty(two, "getBoundingClientRect", {
      configurable: true,
      value: () => rect({ left: 101, right: 131, width: 30 }),
    });
    Object.defineProperty(three, "getBoundingClientRect", {
      configurable: true,
      value: () => rect({ left: 90, right: 150, width: 60 }),
    });

    rerender(
      <MultiSelectPicker
        value={[...value]}
        onQuery={vi.fn()}
        onChange={vi.fn().mockResolvedValue(undefined)}
      />,
    );

    await waitFor(() => {
      expect(within(trigger).getAllByText("+2").length).toBeGreaterThan(0);
    });
  });
});
