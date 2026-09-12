import { FeedBlock, FeedConfig } from "@/lib/settings/feed";
import {
  IndexPageLayout,
  WebMCPToolImplementations,
} from "@/lib/webmcp/tools.generated";

import type { FeedBlockEditorActions } from "./Context";

type IndexPageLayoutEditor = Pick<
  FeedBlockEditorActions,
  "getFeed" | "addBlock" | "moveBlock" | "overwriteBlock" | "removeBlock"
>;

type IndexPageLayoutToolImplementations = Pick<
  WebMCPToolImplementations,
  | "index_page_layout_get"
  | "index_page_block_add"
  | "index_page_block_remove"
  | "index_page_block_move"
  | "index_page_block_categories_update"
  | "index_page_block_threads_update"
  | "index_page_block_quick_share_update"
  | "index_page_block_library_update"
>;

export function createIndexPageLayoutToolImplementations(
  editor: IndexPageLayoutEditor,
): IndexPageLayoutToolImplementations {
  const implementations: IndexPageLayoutToolImplementations = {
    index_page_layout_get: () => {
      const layout = toWebMCPIndexPageLayout(editor.getFeed());

      return {
        message: `Retrieved ${layout.blocks.length} index page layout blocks.`,
        layout,
      };
    },

    index_page_block_add: async ({ block, after_block }) => {
      const feed = editor.getFeed();
      if (feed.blocks.some((candidate) => candidate.type === block)) {
        throw new Error(
          `The ${block} block is already present. Use index_page_block_move to reposition it.`,
        );
      }

      if (after_block === undefined) {
        const updated = await editor.addBlock(block);
        return toolResult(
          updated,
          `Added the ${block} block at the end of the index page.`,
        );
      }

      const afterIndex = feed.blocks.findIndex(
        (candidate) => candidate.type === after_block,
      );
      if (afterIndex === -1) {
        throw new Error(
          `The reference block ${after_block} is not present. Get the current layout with index_page_layout_get and try again.`,
        );
      }

      const updated = await editor.addBlock(block, afterIndex);
      return toolResult(
        updated,
        `Added the ${block} block after the ${after_block} block.`,
      );
    },

    index_page_block_remove: async ({ block }) => {
      if (
        !editor.getFeed().blocks.some((candidate) => candidate.type === block)
      ) {
        throw new Error(
          `The ${block} block is not present, so there is nothing to remove.`,
        );
      }

      const updated = await editor.removeBlock(block);
      return toolResult(
        updated,
        `Removed the ${block} block from the index page without deleting its content.`,
      );
    },

    index_page_block_move: async ({ block, before_block }) => {
      const feed = editor.getFeed();
      const currentIndex = feed.blocks.findIndex(
        (candidate) => candidate.type === block,
      );
      if (currentIndex === -1) {
        throw new Error(
          `The ${block} block is not present. Add it with index_page_block_add before moving it.`,
        );
      }

      if (before_block === undefined) {
        const last = feed.blocks.at(-1);
        const updated =
          last && last.type !== block
            ? await editor.moveBlock(block, last.type)
            : feed;
        return toolResult(updated, `Moved the ${block} block to the end.`);
      }

      if (before_block === block) {
        throw new Error("A block cannot be moved before itself.");
      }

      const beforeIndex = feed.blocks.findIndex(
        (candidate) => candidate.type === before_block,
      );
      if (beforeIndex === -1) {
        throw new Error(
          `The reference block ${before_block} is not present. Get the current layout with index_page_layout_get and try again.`,
        );
      }

      if (currentIndex + 1 === beforeIndex) {
        return toolResult(
          feed,
          `The ${block} block is already immediately before the ${before_block} block.`,
        );
      }

      const over =
        currentIndex < beforeIndex
          ? feed.blocks[beforeIndex - 1]
          : feed.blocks[beforeIndex];
      if (!over) {
        throw new Error("The index page block order could not be resolved.");
      }

      const updated = await editor.moveBlock(block, over.type);
      return toolResult(
        updated,
        `Moved the ${block} block before the ${before_block} block.`,
      );
    },

    index_page_block_categories_update: async ({ layout }) => {
      getBlock(editor.getFeed(), "categories");
      const updated = await editor.overwriteBlock({
        type: "categories",
        layout,
      });
      return toolResult(
        updated,
        `Updated the categories block to the ${layout} layout.`,
      );
    },

    index_page_block_threads_update: async ({ source }) => {
      getBlock(editor.getFeed(), "threads");
      const updated = await editor.overwriteBlock({ type: "threads", source });
      return toolResult(
        updated,
        `Updated the threads block to show ${source === "all" ? "all" : "uncategorised"} threads.`,
      );
    },

    index_page_block_quick_share_update: async ({ show_category_select }) => {
      getBlock(editor.getFeed(), "quick-share");
      const updated = await editor.overwriteBlock({
        type: "quick-share",
        showCategorySelect: show_category_select,
      });
      return toolResult(
        updated,
        `${show_category_select ? "Shown" : "Hidden"} the category picker in the quick-share block.`,
      );
    },

    index_page_block_library_update: async ({ source, node_id, layout }) => {
      const block = getBlock(editor.getFeed(), "library");

      if (source === "library_page") {
        if (!node_id?.trim()) {
          throw new Error(
            "node_id is required when the Library block source is library_page.",
          );
        }
        if (layout !== undefined) {
          throw new Error(
            "layout only applies to the Library root. A selected Library page controls its own layout.",
          );
        }

        const updated = await editor.overwriteBlock({
          type: "library",
          node: node_id,
          layout: block.layout,
        });
        return toolResult(
          updated,
          `Updated the Library block to show page ${node_id}.`,
        );
      }

      if (node_id !== undefined) {
        throw new Error(
          "node_id must be omitted when the Library block source is library_root.",
        );
      }

      const updated = await editor.overwriteBlock({
        type: "library",
        layout: layout ?? block.layout,
      });
      return toolResult(
        updated,
        `Updated the Library block to show the Library root in the ${layout ?? block.layout} layout.`,
      );
    },
  };

  let mutationQueue = Promise.resolve();
  const enqueueMutation = <T>(
    operation: () => T | PromiseLike<T>,
  ): Promise<T> => {
    const result = mutationQueue.then(operation);
    mutationQueue = result.then(
      () => undefined,
      () => undefined,
    );
    return result;
  };

  return {
    ...implementations,
    index_page_block_add: (input) =>
      enqueueMutation(() => implementations.index_page_block_add(input)),
    index_page_block_remove: (input) =>
      enqueueMutation(() => implementations.index_page_block_remove(input)),
    index_page_block_move: (input) =>
      enqueueMutation(() => implementations.index_page_block_move(input)),
    index_page_block_categories_update: (input) =>
      enqueueMutation(() =>
        implementations.index_page_block_categories_update(input),
      ),
    index_page_block_threads_update: (input) =>
      enqueueMutation(() =>
        implementations.index_page_block_threads_update(input),
      ),
    index_page_block_quick_share_update: (input) =>
      enqueueMutation(() =>
        implementations.index_page_block_quick_share_update(input),
      ),
    index_page_block_library_update: (input) =>
      enqueueMutation(() =>
        implementations.index_page_block_library_update(input),
      ),
  };
}

export function toWebMCPIndexPageLayout(feed: FeedConfig): IndexPageLayout {
  return { blocks: feed.blocks.map(toWebMCPBlock) };
}

function toWebMCPBlock(block: FeedBlock): IndexPageLayout["blocks"][number] {
  switch (block.type) {
    case "categories":
      return { type: block.type, layout: block.layout };
    case "threads":
      return { type: block.type, source: block.source };
    case "quick-share":
      return {
        type: block.type,
        show_category_select: block.showCategorySelect,
      };
    case "library":
      return {
        type: block.type,
        layout: block.layout,
        ...(block.node ? { node_id: block.node } : {}),
      };
    default:
      return { type: block.type };
  }
}

function getBlock<K extends FeedBlock["type"]>(
  feed: FeedConfig,
  type: K,
): Extract<FeedBlock, { type: K }> {
  const block = feed.blocks.find(
    (candidate): candidate is Extract<FeedBlock, { type: K }> =>
      candidate.type === type,
  );

  if (!block) {
    throw new Error(
      `The ${type} block is not present. Add it with index_page_block_add before updating it.`,
    );
  }

  return block;
}

function toolResult(feed: FeedConfig, message: string) {
  return { message, layout: toWebMCPIndexPageLayout(feed) };
}
