import { describe, expect, it, vi } from "vitest";

import { NodeWithChildren } from "@/api/openapi-schema";
import { hydrateNode } from "@/lib/library/metadata";

import { createLibraryPageLayoutToolImplementations } from "./libraryPageLayoutTools";
import { createNodeStore } from "./store";

describe("Library page WebMCP tool implementations", () => {
  it("gets the current layout in its frontend tool representation", async () => {
    const tools = createTools();

    expect(await tools.library_page_layout_get({})).toEqual({
      message: "Retrieved 3 Library page layout blocks.",
      layout: {
        blocks: [
          { type: "title" },
          { type: "assets", layout: "grid", grid_size: 3 },
          { type: "directory", layout: "table" },
        ],
      },
    });
  });

  it("adds a hidden block after a visible reference block", async () => {
    const store = createStore();
    const waitForSave = vi.fn(async () => {});
    const tools = createLibraryPageLayoutToolImplementations(
      store,
      waitForSave,
    );

    const output = await tools.library_page_block_add({
      block: "tags",
      after_block: "assets",
    });

    expect(output).toEqual({
      message: "Added the tags block after the assets block.",
      layout: {
        blocks: [
          { type: "title" },
          { type: "assets", layout: "grid", grid_size: 3 },
          { type: "tags" },
          { type: "directory", layout: "table" },
        ],
      },
    });
    expect(waitForSave).toHaveBeenCalledWith(store, output.layout);
  });

  it("removes a visible block without changing other block configuration", async () => {
    const tools = createTools();

    const output = await tools.library_page_block_remove({ block: "title" });

    expect(output.layout.blocks).toEqual([
      { type: "assets", layout: "grid", grid_size: 3 },
      { type: "directory", layout: "table" },
    ]);
  });

  it("moves a visible block before another visible block", async () => {
    const tools = createTools();

    const output = await tools.library_page_block_move({
      block: "directory",
      before_block: "assets",
    });

    expect(output.layout.blocks).toEqual([
      { type: "title" },
      { type: "directory", layout: "table" },
      { type: "assets", layout: "grid", grid_size: 3 },
    ]);
  });

  it("moves a visible block to the end when no reference block is provided", async () => {
    const tools = createTools();

    const output = await tools.library_page_block_move({ block: "title" });

    expect(output.layout.blocks).toEqual([
      { type: "assets", layout: "grid", grid_size: 3 },
      { type: "directory", layout: "table" },
      { type: "title" },
    ]);
  });

  it("updates assets presentation without discarding its existing grid size", async () => {
    const tools = createTools();

    const output = await tools.library_page_block_assets_update({
      layout: "strip",
    });

    expect(output.layout.blocks[1]).toEqual({
      type: "assets",
      layout: "strip",
      grid_size: 3,
    });
  });

  it("updates directory presentation without exposing or discarding columns", async () => {
    const store = createStore();
    const tools = createLibraryPageLayoutToolImplementations(
      store,
      async () => {},
    );

    const output = await tools.library_page_block_directory_update({
      layout: "grid",
    });

    expect(output.layout.blocks[2]).toEqual({
      type: "directory",
      layout: "grid",
    });
    expect(
      store
        .getState()
        .draft.meta.layout?.blocks.find((block) => block.type === "directory"),
    ).toEqual({
      type: "directory",
      config: {
        layout: "grid",
        columns: [{ fid: "field-1", hidden: false }],
      },
    });
  });

  it("returns actionable errors for invalid progressive changes", async () => {
    const tools = createTools();

    await expect(
      tools.library_page_block_add({ block: "title" }),
    ).rejects.toThrow(
      "The title block is already present. Use library_page_block_move to reposition it.",
    );
    await expect(
      tools.library_page_block_remove({ block: "tags" }),
    ).rejects.toThrow(
      "The tags block is not present. Add it with library_page_block_add before updating it.",
    );
    await expect(
      tools.library_page_block_move({
        block: "title",
        before_block: "title",
      }),
    ).rejects.toThrow("A block cannot be moved before itself.");
  });

  it("rejects a block-specific update when that block is hidden", async () => {
    const store = createStore();
    const tools = createLibraryPageLayoutToolImplementations(
      store,
      async () => {},
    );
    store.getState().removeBlock("assets");

    await expect(
      tools.library_page_block_assets_update({ layout: "grid" }),
    ).rejects.toThrow(
      "The assets block is not present. Add it with library_page_block_add before updating it.",
    );
  });

  it("serializes mutations so each call waits for its own saved layout", async () => {
    const store = createStore();
    let releaseFirstSave = () => {};
    const firstSave = new Promise<void>((resolve) => {
      releaseFirstSave = resolve;
    });
    const waitForSave = vi
      .fn()
      .mockImplementationOnce(() => firstSave)
      .mockResolvedValueOnce(undefined);
    const tools = createLibraryPageLayoutToolImplementations(
      store,
      waitForSave,
    );

    const add = tools.library_page_block_add({ block: "tags" });
    const remove = tools.library_page_block_remove({ block: "title" });

    await vi.waitFor(() => expect(waitForSave).toHaveBeenCalledTimes(1));
    expect(
      store
        .getState()
        .draft.meta.layout?.blocks.some((block) => block.type === "title"),
    ).toBe(true);

    releaseFirstSave();
    const [added, removed] = await Promise.all([add, remove]);

    expect(waitForSave).toHaveBeenNthCalledWith(1, store, added.layout);
    expect(waitForSave).toHaveBeenNthCalledWith(2, store, removed.layout);
    expect(added.layout.blocks.map((block) => block.type)).toEqual([
      "title",
      "assets",
      "directory",
      "tags",
    ]);
    expect(removed.layout.blocks.map((block) => block.type)).toEqual([
      "assets",
      "directory",
      "tags",
    ]);
  });
});

function createTools() {
  return createLibraryPageLayoutToolImplementations(
    createStore(),
    async () => {},
  );
}

function createStore() {
  const original = hydrateNode(node());

  return createNodeStore({
    original,
    draft: structuredClone(original),
  });
}

function node(): NodeWithChildren {
  return {
    id: "library",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    name: "Library",
    slug: "library",
    description: "",
    content: "",
    owner: {
      id: "author",
      handle: "author",
      name: "Author",
      joined: "2026-01-01T00:00:00Z",
      roles: [],
    },
    primary_image: undefined,
    assets: [],
    tags: [],
    link: undefined,
    visibility: "published",
    hide_child_tree: false,
    meta: {
      layout: {
        blocks: [
          { type: "title" },
          { type: "assets", config: { layout: "grid", gridSize: 3 } },
          {
            type: "directory",
            config: {
              layout: "table",
              columns: [{ fid: "field-1", hidden: false }],
            },
          },
        ],
      },
    },
    properties: [],
    child_property_schema: [],
    children: [],
    recomentations: [],
  };
}
