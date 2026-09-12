import { describe, expect, it, vi } from "vitest";

import {
  FeedBlock,
  FeedBlockType,
  FeedConfig,
  createFeedBlock,
} from "@/lib/settings/feed";

import { createIndexPageLayoutToolImplementations } from "./indexPageLayoutTools";

describe("index page WebMCP tool implementations", () => {
  it("gets each block in its frontend tool representation", async () => {
    const { tools } = createTools({
      blocks: [
        { type: "title" },
        { type: "categories", layout: "grid" },
        { type: "threads", source: "all" },
        { type: "quick-share", showCategorySelect: false },
        { type: "library", layout: "list", node: "library-page-id" },
      ],
    });

    expect(await tools.index_page_layout_get({})).toEqual({
      message: "Retrieved 5 index page layout blocks.",
      layout: {
        blocks: [
          { type: "title" },
          { type: "categories", layout: "grid" },
          { type: "threads", source: "all" },
          { type: "quick-share", show_category_select: false },
          {
            type: "library",
            layout: "list",
            node_id: "library-page-id",
          },
        ],
      },
    });
  });

  it("applies progressive add, configure, move, and remove operations", async () => {
    const { tools } = createTools({
      blocks: [
        { type: "title" },
        { type: "categories", layout: "list" },
        { type: "threads", source: "uncategorised" },
      ],
    });

    await tools.index_page_block_add({
      block: "quick-share",
      after_block: "categories",
    });
    await tools.index_page_block_quick_share_update({
      show_category_select: false,
    });
    await tools.index_page_block_move({
      block: "threads",
      before_block: "categories",
    });
    const removed = await tools.index_page_block_remove({ block: "title" });

    expect(removed.layout.blocks).toEqual([
      { type: "threads", source: "uncategorised" },
      { type: "categories", layout: "list" },
      { type: "quick-share", show_category_select: false },
    ]);
  });

  it("configures each specialised block without replacing its neighbours", async () => {
    const { tools } = createTools({
      blocks: [
        { type: "categories", layout: "list" },
        { type: "threads", source: "uncategorised" },
        { type: "quick-share", showCategorySelect: true },
        { type: "library", layout: "list" },
      ],
    });

    await tools.index_page_block_categories_update({ layout: "grid" });
    await tools.index_page_block_threads_update({ source: "all" });
    await tools.index_page_block_quick_share_update({
      show_category_select: false,
    });
    const library = await tools.index_page_block_library_update({
      source: "library_page",
      node_id: "page-id",
    });

    expect(library.layout.blocks).toEqual([
      { type: "categories", layout: "grid" },
      { type: "threads", source: "all" },
      { type: "quick-share", show_category_select: false },
      { type: "library", layout: "list", node_id: "page-id" },
    ]);
  });

  it("switches a selected Library page back to the Library root", async () => {
    const { tools } = createTools({
      blocks: [
        {
          type: "library",
          layout: "list",
          node: "selected-library-page",
        },
      ],
    });

    const output = await tools.index_page_block_library_update({
      source: "library_root",
      layout: "grid",
    });

    expect(output.layout.blocks).toEqual([{ type: "library", layout: "grid" }]);
  });

  it("returns actionable errors for stale or ambiguous operations", async () => {
    const { tools } = createTools({
      blocks: [{ type: "title" }, { type: "library", layout: "grid" }],
    });

    await expect(
      tools.index_page_block_add({ block: "title" }),
    ).rejects.toThrow(
      "The title block is already present. Use index_page_block_move to reposition it.",
    );
    await expect(
      tools.index_page_block_move({
        block: "title",
        before_block: "title",
      }),
    ).rejects.toThrow("A block cannot be moved before itself.");
    await expect(
      tools.index_page_block_library_update({ source: "library_page" }),
    ).rejects.toThrow(
      "node_id is required when the Library block source is library_page.",
    );
    await expect(
      tools.index_page_block_library_update({
        source: "library_page",
        node_id: "page-id",
        layout: "list",
      }),
    ).rejects.toThrow("layout only applies to the Library root");
  });

  it("serializes mutations so later calls use the saved evolving state", async () => {
    let releaseFirstSave = () => {};
    const firstSave = new Promise<void>((resolve) => {
      releaseFirstSave = resolve;
    });
    const { tools, editor } = createTools(
      { blocks: [{ type: "title" }, { type: "content" }] },
      firstSave,
    );

    const add = tools.index_page_block_add({ block: "cover" });
    const remove = tools.index_page_block_remove({ block: "title" });

    await vi.waitFor(() => expect(editor.addBlock).toHaveBeenCalledTimes(1));
    expect(editor.removeBlock).not.toHaveBeenCalled();

    releaseFirstSave();
    const [added, removed] = await Promise.all([add, remove]);

    expect(added.layout.blocks.map((block) => block.type)).toEqual([
      "title",
      "content",
      "cover",
    ]);
    expect(removed.layout.blocks.map((block) => block.type)).toEqual([
      "content",
      "cover",
    ]);
  });
});

function createTools(initial: FeedConfig, firstSave?: Promise<void>) {
  let feed = structuredClone(initial);
  let saveCount = 0;

  const save = async (updated: FeedConfig) => {
    feed = updated;
    saveCount += 1;
    if (saveCount === 1 && firstSave) {
      await firstSave;
    }
    return feed;
  };

  const editor = {
    getFeed: () => feed,
    addBlock: vi.fn(async (type: FeedBlockType, afterIndex?: number) => {
      const blocks = [...feed.blocks];
      blocks.splice(
        afterIndex === undefined ? blocks.length : afterIndex + 1,
        0,
        createFeedBlock(type),
      );
      return save({ blocks });
    }),
    moveBlock: vi.fn(async (active: FeedBlockType, over: FeedBlockType) => {
      const blocks = [...feed.blocks];
      const activeIndex = blocks.findIndex((block) => block.type === active);
      const overIndex = blocks.findIndex((block) => block.type === over);
      const [moved] = blocks.splice(activeIndex, 1);
      if (!moved) throw new Error("Missing moved block in test editor");
      blocks.splice(overIndex, 0, moved);
      return save({ blocks });
    }),
    overwriteBlock: vi.fn(async (updated: FeedBlock) => {
      const blocks = feed.blocks.map((block) =>
        block.type === updated.type ? updated : block,
      );
      return save({ blocks });
    }),
    removeBlock: vi.fn(async (type: FeedBlockType) => {
      return save({
        blocks: feed.blocks.filter((block) => block.type !== type),
      });
    }),
  };

  return {
    editor,
    tools: createIndexPageLayoutToolImplementations(editor),
  };
}
