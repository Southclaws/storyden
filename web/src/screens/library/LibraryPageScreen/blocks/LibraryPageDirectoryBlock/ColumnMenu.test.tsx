import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ColumnMenu } from "./ColumnMenu";

const actions = {
  setChildPropertyHiddenState: vi.fn(),
  setChildPropertyName: vi.fn(),
  removeChildPropertyByID: vi.fn(),
};

let editingMock = true;

vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: {
      getState: () => actions,
    },
  }),
}));

vi.mock("../../useEditState", () => ({
  useEditState: () => ({
    editing: editingMock,
  }),
}));

describe("ColumnMenu", () => {
  beforeEach(() => {
    editingMock = true;
    actions.setChildPropertyHiddenState.mockReset();
    actions.setChildPropertyName.mockReset();
    actions.removeChildPropertyByID.mockReset();
  });

  it("updates name, toggles hidden state, and deletes non-fixed columns", async () => {
    const user = userEvent.setup();
    const column = {
      fid: "p1",
      name: "Priority",
      type: "text",
      fixed: false,
      hidden: false,
    } as const;

    render(
      <ColumnMenu column={column}>
        <button>Open menu</button>
      </ColumnMenu>,
    );

    await user.click(screen.getByRole("button", { name: "Open menu" }));

    const input = screen.getByRole("textbox");
    fireEvent.change(input, { target: { value: "Urgency" } });
    expect(actions.setChildPropertyName).toHaveBeenCalledWith("p1", "Urgency");

    await user.click(screen.getByRole("menuitem", { name: /Hide column/ }));
    expect(actions.setChildPropertyHiddenState).toHaveBeenCalledWith("p1", true);

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    await user.click(screen.getByRole("menuitem", { name: /Delete/ }));
    expect(actions.removeChildPropertyByID).toHaveBeenCalledWith("p1");
  });

  it("does not render rename/delete controls for fixed columns", async () => {
    const user = userEvent.setup();
    const fixedColumn = {
      fid: "fixed:name",
      name: "name",
      type: "text",
      fixed: true,
      hidden: true,
      _fixedFieldName: "name",
    } as const;

    render(
      <ColumnMenu column={fixedColumn}>
        <button>Open menu</button>
      </ColumnMenu>,
    );

    await user.click(screen.getByRole("button", { name: "Open menu" }));

    expect(screen.queryByRole("textbox")).toBeNull();
    expect(screen.queryByRole("menuitem", { name: /Delete/ })).toBeNull();
    expect(screen.getByRole("menuitem", { name: /Hide column/ })).toBeInTheDocument();
  });

  it("stays closed when edit mode is disabled", async () => {
    const user = userEvent.setup();
    editingMock = false;

    const column = {
      fid: "p1",
      name: "Priority",
      type: "text",
      fixed: false,
      hidden: false,
    } as const;

    render(
      <ColumnMenu column={column}>
        <button>Open menu</button>
      </ColumnMenu>,
    );

    await user.click(screen.getByRole("button", { name: "Open menu" }));

    await waitFor(() => {
      expect(screen.queryByRole("menuitem", { name: "Hide column" })).toBeNull();
    });
  });
});
