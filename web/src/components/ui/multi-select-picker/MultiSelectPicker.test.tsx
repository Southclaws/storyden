import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { deriveColour } from "@/utils/colour";

import { MultiSelectPicker, MultiSelectPickerItem } from ".";

function openPicker() {
  return userEvent.click(screen.getByRole("button", { name: /Select items/i }));
}

function getSearchInput() {
  return screen.getByRole("combobox", { name: "Search for items" });
}

describe("MultiSelectPicker", () => {
  const baseValue: MultiSelectPickerItem[] = [
    { label: "Alpha", value: "alpha" },
  ];

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

  it("queries the default options again when the search is cleared", async () => {
    const user = userEvent.setup();
    const onQuery = vi.fn();

    render(
      <MultiSelectPicker
        value={[]}
        onQuery={onQuery}
        onChange={vi.fn().mockResolvedValue(undefined)}
      />,
    );

    await user.type(getSearchInput(), "beta");
    await user.clear(getSearchInput());

    expect(getSearchInput()).toHaveValue("");
    expect(onQuery.mock.calls.at(-1)?.[0]).toBe("");
  });

  it("describes the current selection while the popup is closed", () => {
    const props = {
      onQuery: vi.fn(),
      onChange: vi.fn().mockResolvedValue(undefined),
    };
    const { rerender } = render(
      <MultiSelectPicker {...props} value={baseValue} />,
    );

    expect(getSearchInput()).toHaveAccessibleDescription("Selected: Alpha.");
    expect(screen.queryByRole("listbox")).not.toBeInTheDocument();

    rerender(<MultiSelectPicker {...props} value={[]} />);

    expect(getSearchInput()).toHaveAccessibleDescription("No items selected.");
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
    await user.click(screen.getByRole("option", { name: "Beta" }));

    await waitFor(() =>
      expect(onChange).toHaveBeenCalledWith([
        ...baseValue,
        { ...result, colour: deriveColour("beta") },
      ]),
    );
  });

  it("selects the highlighted result with the keyboard", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);
    const result = { label: "Beta", value: "beta" };

    render(
      <MultiSelectPicker
        value={[]}
        onQuery={vi.fn()}
        onChange={onChange}
        queryResults={[result]}
      />,
    );

    const input = getSearchInput();
    await user.type(input, "bet");
    await user.keyboard("{ArrowDown}");

    expect(input).toHaveFocus();
    expect(input).toHaveAttribute(
      "aria-activedescendant",
      screen.getByRole("option", { name: "Beta" }).id,
    );

    await user.keyboard("{Enter}");

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([result]));
  });

  it("creates a new item from Enter without submitting an enclosing form", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);
    const onSubmit = vi.fn((event: React.FormEvent) => event.preventDefault());

    render(
      <form onSubmit={onSubmit}>
        <MultiSelectPicker
          value={[]}
          onQuery={vi.fn()}
          onChange={onChange}
          allowNewValues
        />
      </form>,
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
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it("removes a selected item by toggling its selected option", async () => {
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
    await user.click(screen.getByRole("option", { name: "Alpha" }));

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([]));
  });

  it("removes the last selected item with Backspace from an empty input", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);
    const value = [...baseValue, { label: "Beta", value: "beta" }];

    render(
      <MultiSelectPicker value={value} onQuery={vi.fn()} onChange={onChange} />,
    );

    const input = getSearchInput();
    await user.click(input);
    await user.keyboard("{Backspace}");

    expect(onChange).toHaveBeenCalledWith(baseValue);
  });

  it("keeps selected items when Backspace edits a query", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn().mockResolvedValue(undefined);

    render(
      <MultiSelectPicker
        value={baseValue}
        onQuery={vi.fn()}
        onChange={onChange}
      />,
    );

    const input = getSearchInput();
    await user.type(input, "a{Backspace}");

    expect(onChange).not.toHaveBeenCalled();
  });

  it("shows loading instead of an empty state until options arrive", async () => {
    const props = {
      value: [],
      onQuery: vi.fn(),
      onChange: vi.fn().mockResolvedValue(undefined),
    };
    const { rerender } = render(<MultiSelectPicker {...props} queryLoading />);

    await openPicker();
    expect(screen.getByRole("status", { name: "Loading" })).toBeVisible();
    expect(screen.getByRole("listbox")).toHaveAttribute("aria-busy", "true");
    expect(screen.queryByText("No options available")).not.toBeInTheDocument();

    rerender(
      <MultiSelectPicker
        {...props}
        queryResults={[{ label: "Beta", value: "beta" }]}
      />,
    );

    expect(screen.getByRole("option", { name: "Beta" })).toBeVisible();
    expect(
      screen.queryByRole("status", { name: "Loading" }),
    ).not.toBeInTheDocument();
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

    expect(screen.getByText("No results for “zzz”")).toBeInTheDocument();
  });

  it("keeps receiving text input after opening from its trigger", async () => {
    const user = userEvent.setup();
    const onQuery = vi.fn();

    render(
      <MultiSelectPicker
        value={baseValue}
        onQuery={onQuery}
        onChange={vi.fn().mockResolvedValue(undefined)}
        queryResults={[{ label: "Beta", value: "beta" }]}
      />,
    );

    await openPicker();
    const input = getSearchInput();
    await user.type(input, "bet");

    expect(input).toHaveValue("bet");
    expect(onQuery.mock.calls.at(-1)?.[0]).toBe("bet");
    expect(screen.getByRole("option", { name: "Beta" })).toBeVisible();
  });

  it("does not constrain the popup from its zero-size first measurement", async () => {
    render(
      <MultiSelectPicker
        value={baseValue}
        onQuery={vi.fn()}
        onChange={vi.fn().mockResolvedValue(undefined)}
        queryResults={[{ label: "Beta", value: "beta" }]}
      />,
    );

    await openPicker();

    await waitFor(() =>
      expect(screen.getByRole("listbox").style.maxHeight).toBe(""),
    );
  });
});
