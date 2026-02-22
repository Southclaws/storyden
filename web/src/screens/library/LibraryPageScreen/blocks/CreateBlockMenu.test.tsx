import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { CreateBlockMenu } from "./CreateBlockMenu";

const emit = vi.fn();
let metadataMock: any;

vi.mock("../store", () => ({
  useWatch: (selector: (state: any) => unknown) =>
    selector({
      draft: {
        meta: metadataMock,
      },
    }),
}));

vi.mock("@/lib/library/events", () => ({
  useEmitLibraryBlockEvent: () => emit,
}));

describe("CreateBlockMenu", () => {
  beforeEach(() => {
    emit.mockReset();
    metadataMock = {
      layout: {
        blocks: [{ type: "title" }, { type: "cover" }],
      },
    };
  });

  it("hides blocks that already exist in metadata", async () => {
    const user = userEvent.setup();

    render(<CreateBlockMenu trigger={<button>Open menu</button>} />);

    await user.click(screen.getByRole("button", { name: "Open menu" }));

    expect(screen.queryByRole("menuitem", { name: "Title" })).toBeNull();
    expect(screen.queryByRole("menuitem", { name: "Cover image" })).toBeNull();
    expect(screen.getByRole("menuitem", { name: "Directory" })).toBeInTheDocument();
  });

  it("emits add-block with selected block type", async () => {
    const user = userEvent.setup();

    render(<CreateBlockMenu trigger={<button>Open menu</button>} />);

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    await user.click(screen.getByRole("menuitem", { name: "Directory" }));

    await waitFor(() =>
      expect(emit).toHaveBeenCalledWith("library:add-block", {
        type: "directory",
        index: undefined,
      }),
    );
  });

  it("preserves index=0 when emitting add-block", async () => {
    const user = userEvent.setup();

    render(<CreateBlockMenu trigger={<button>Open menu</button>} index={0} />);

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    await user.click(screen.getByRole("menuitem", { name: "Directory" }));

    await waitFor(() =>
      expect(emit).toHaveBeenCalledWith("library:add-block", {
        type: "directory",
        index: 0,
      }),
    );
  });
});
