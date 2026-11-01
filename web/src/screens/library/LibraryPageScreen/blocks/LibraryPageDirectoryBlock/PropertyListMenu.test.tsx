import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { LibraryPageBlockTypeDirectory } from "@/lib/library/metadata";

import { PropertyListMenu } from "./PropertyListMenu/PropertyListMenu";

const overwriteBlock = vi.fn();
const setChildPropertyHiddenState = vi.fn();
let block: LibraryPageBlockTypeDirectory;
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: {
      getState: () => ({ overwriteBlock, setChildPropertyHiddenState }),
    },
  }),
}));
vi.mock("../../store", () => ({
  useWatch: () => [{ fid: "status", name: "Status", type: "text", sort: "0" }],
}));
vi.mock("./useDirectoryBlock", () => ({ useDirectoryBlock: () => block }));

describe("Directory display controls", () => {
  beforeEach(() => {
    block = {
      type: "directory",
      config: {
        layout: "grid",
        card_style: "box",
        columns: [
          { fid: "fixed:name", hidden: false },
          { fid: "fixed:primary_image", hidden: false },
          { fid: "status", hidden: false },
        ],
      },
    };
  });

  it("places a property without changing card style or unrelated visibility", async () => {
    const user = userEvent.setup();
    render(
      <PropertyListMenu>
        <button>Display</button>
      </PropertyListMenu>,
    );
    await user.click(screen.getByRole("button", { name: "Display" }));
    await waitFor(() => expect(screen.getByRole("menu")).toHaveFocus());
    await user.keyboard("{End}{ArrowRight}");
    await screen.findByRole("menuitemradio", { name: "Bottom-left badge" });
    await waitFor(() =>
      expect(screen.getAllByRole("menu").at(-1)).toHaveFocus(),
    );
    await user.keyboard("{Home}{ArrowDown}{ArrowDown}{Enter}");
    await waitFor(() =>
      expect(overwriteBlock).toHaveBeenCalledWith({
        type: "directory",
        config: {
          layout: "grid",
          card_style: "box",
          columns: [
            { fid: "fixed:name", hidden: false },
            { fid: "fixed:primary_image", hidden: false },
            { fid: "status", hidden: false, placement: "footer-start" },
          ],
        },
      }),
    );
    expect(setChildPropertyHiddenState).not.toHaveBeenCalled();
  });

  it("toggles rich content independently of the card style and columns", async () => {
    const user = userEvent.setup();
    const { rerender } = render(
      <PropertyListMenu>
        <button>Display</button>
      </PropertyListMenu>,
    );
    await user.click(screen.getByRole("button", { name: "Display" }));
    const toggle = await screen.findByRole("menuitemcheckbox", {
      name: "Rich content",
    });
    expect(toggle).toHaveAttribute("aria-checked", "false");
    await user.click(toggle);
    expect(overwriteBlock).toHaveBeenLastCalledWith({
      ...block,
      config: { ...block.config, show_content: true },
    });

    block = { ...block, config: { ...block.config!, show_content: true } };
    rerender(
      <PropertyListMenu>
        <button>Display</button>
      </PropertyListMenu>,
    );
    const enabled = await screen.findByRole("menuitemcheckbox", {
      name: "Rich content",
    });
    expect(enabled).toHaveAttribute("aria-checked", "true");
    await user.click(enabled);
    expect(overwriteBlock).toHaveBeenLastCalledWith({
      ...block,
      config: { ...block.config, show_content: false },
    });
  });

  it.each(["table", "grid"] as const)(
    "does not offer rich content in the %s minimal presentation",
    async (layout) => {
      block.config = { ...block.config!, layout, card_style: "minimal" };
      const user = userEvent.setup();
      render(
        <PropertyListMenu>
          <button>Display</button>
        </PropertyListMenu>,
      );
      await user.click(screen.getByRole("button", { name: "Display" }));
      await screen.findByRole("menu");
      expect(
        screen.queryByRole("menuitemcheckbox", { name: "Rich content" }),
      ).toBeNull();
    },
  );

  it("selects Rows while retaining the saved card style and placements", async () => {
    const user = userEvent.setup();
    render(
      <PropertyListMenu>
        <button>Display</button>
      </PropertyListMenu>,
    );
    await user.click(screen.getByRole("button", { name: "Display" }));
    await waitFor(() => expect(screen.getByRole("menu")).toHaveFocus());
    await user.keyboard("{Home}{ArrowRight}");
    await screen.findByRole("menuitemradio", { name: /Rows/ });
    await waitFor(() =>
      expect(screen.getAllByRole("menu").at(-1)).toHaveFocus(),
    );
    await user.keyboard("{Home}{ArrowDown}{Enter}");
    await waitFor(() =>
      expect(overwriteBlock).toHaveBeenCalledWith({
        ...block,
        config: { ...block.config, layout: "rows" },
      }),
    );
  });

  it("only offers masonry-compatible styles in Gallery", async () => {
    block.config!.layout = "gallery";
    const user = userEvent.setup();
    render(
      <PropertyListMenu>
        <button>Display</button>
      </PropertyListMenu>,
    );
    await user.click(screen.getByRole("button", { name: "Display" }));
    await waitFor(() => expect(screen.getByRole("menu")).toHaveFocus());
    await user.keyboard("{Home}{ArrowDown}{ArrowRight}");
    await waitFor(() =>
      expect(
        screen.getByRole("menuitemradio", { name: /Editorial/ }),
      ).toBeVisible(),
    );
    expect(
      screen.queryByRole("menuitemradio", { name: /Image overlay/ }),
    ).toBeNull();
    expect(screen.queryByRole("menuitemradio", { name: /Compact/ })).toBeNull();
    await waitFor(() =>
      expect(screen.getAllByRole("menu").at(-1)).toHaveFocus(),
    );
    await user.keyboard("{Home}{Enter}");
    await waitFor(() =>
      expect(overwriteBlock).toHaveBeenCalledWith({
        ...block,
        config: { ...block.config, card_style: "minimal" },
      }),
    );
  });
});
