import {
  DefaultLayout,
  LibraryPageBlock,
  NodeLayout,
} from "@/lib/library/metadata";
import {
  LibraryPageLayout,
  WebMCPToolImplementations,
} from "@/lib/webmcp/tools.generated";
import { deepEqual } from "@/utils/equality";

import { NodeStoreAPI } from "./store";

const AUTOSAVE_TIMEOUT_MS = 10_000;

type WaitForSave = (
  store: NodeStoreAPI,
  expected: LibraryPageLayout,
) => Promise<void>;

export function createLibraryPageLayoutToolImplementations(
  store: NodeStoreAPI,
  waitForSave: WaitForSave = waitForLibraryPageLayoutSave,
): WebMCPToolImplementations {
  const implementations: WebMCPToolImplementations = {
    library_page_layout_get: () => {
      const layout = getDraftLayout(store);

      return {
        message: `Retrieved ${layout.blocks.length} Library page layout blocks.`,
        layout,
      };
    },

    library_page_block_add: async ({ block, after_block }) => {
      const blocks = getDraftBlocks(store);
      if (blocks.some((candidate) => candidate.type === block)) {
        throw new Error(
          `The ${block} block is already present. Use library_page_block_move to reposition it.`,
        );
      }

      if (after_block === undefined) {
        store.getState().addBlock(block);
        return saveLayout(
          store,
          waitForSave,
          `Added the ${block} block at the end.`,
        );
      }

      const afterIndex = blocks.findIndex(
        (candidate) => candidate.type === after_block,
      );
      if (afterIndex === -1) {
        throw new Error(
          `The reference block ${after_block} is not present. Get the current layout with library_page_layout_get and try again.`,
        );
      }

      store.getState().addBlock(block, afterIndex);
      return saveLayout(
        store,
        waitForSave,
        `Added the ${block} block after the ${after_block} block.`,
      );
    },

    library_page_block_remove: async ({ block }) => {
      if (
        !getDraftBlocks(store).some((candidate) => candidate.type === block)
      ) {
        throw new Error(
          `The ${block} block is not present, so there is nothing to remove.`,
        );
      }
      store.getState().removeBlock(block);

      return saveLayout(
        store,
        waitForSave,
        `Removed the ${block} block from the page layout without deleting its content.`,
      );
    },

    library_page_block_move: async ({ block, before_block }) => {
      const blocks = getDraftBlocks(store);
      const currentIndex = blocks.findIndex(
        (candidate) => candidate.type === block,
      );
      if (currentIndex === -1) {
        throw new Error(
          `The ${block} block is not present. Add it with library_page_block_add before moving it.`,
        );
      }

      if (before_block === undefined) {
        store.getState().moveBlock(block, blocks.length - 1);
        return saveLayout(
          store,
          waitForSave,
          `Moved the ${block} block to the end.`,
        );
      }

      if (before_block === block) {
        throw new Error("A block cannot be moved before itself.");
      }

      const beforeIndex = blocks.findIndex(
        (candidate) => candidate.type === before_block,
      );
      if (beforeIndex === -1) {
        throw new Error(
          `The reference block ${before_block} is not present. Get the current layout with library_page_layout_get and try again.`,
        );
      }

      const newIndex =
        currentIndex < beforeIndex ? beforeIndex - 1 : beforeIndex;
      store.getState().moveBlock(block, newIndex);

      return saveLayout(
        store,
        waitForSave,
        `Moved the ${block} block before the ${before_block} block.`,
      );
    },

    library_page_block_assets_update: async ({ layout, grid_size }) => {
      const block = getBlock(store, "assets");
      store.getState().overwriteBlock({
        type: "assets",
        config: {
          layout,
          ...(grid_size !== undefined
            ? { gridSize: grid_size }
            : block.config?.gridSize !== undefined
              ? { gridSize: block.config.gridSize }
              : {}),
        },
      });

      return saveLayout(
        store,
        waitForSave,
        `Updated the assets block to the ${layout} layout.`,
      );
    },

    library_page_block_directory_update: async ({ layout }) => {
      const block = getBlock(store, "directory");
      store.getState().overwriteBlock({
        type: "directory",
        config: {
          layout,
          columns: block.config?.columns ?? [],
        },
      });

      return saveLayout(
        store,
        waitForSave,
        `Updated the directory block to the ${layout} layout.`,
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
    library_page_block_add: (input) =>
      enqueueMutation(() => implementations.library_page_block_add(input)),
    library_page_block_remove: (input) =>
      enqueueMutation(() => implementations.library_page_block_remove(input)),
    library_page_block_move: (input) =>
      enqueueMutation(() => implementations.library_page_block_move(input)),
    library_page_block_assets_update: (input) =>
      enqueueMutation(() =>
        implementations.library_page_block_assets_update(input),
      ),
    library_page_block_directory_update: (input) =>
      enqueueMutation(() =>
        implementations.library_page_block_directory_update(input),
      ),
  };
}

export function toWebMCPLayout(layout?: NodeLayout): LibraryPageLayout {
  return {
    blocks: (layout ?? DefaultLayout).blocks.map(toWebMCPBlock),
  };
}

function toWebMCPBlock(
  block: LibraryPageBlock,
): LibraryPageLayout["blocks"][number] {
  if (block.type === "assets") {
    return {
      type: "assets",
      ...(block.config?.layout ? { layout: block.config.layout } : {}),
      ...(block.config?.gridSize !== undefined
        ? { grid_size: block.config.gridSize }
        : {}),
    };
  }

  if (block.type === "directory") {
    return {
      type: "directory",
      ...(block.config?.layout ? { layout: block.config.layout } : {}),
    };
  }

  return { type: block.type };
}

function getDraftLayout(store: NodeStoreAPI): LibraryPageLayout {
  return toWebMCPLayout(store.getState().draft.meta.layout);
}

function getDraftBlocks(store: NodeStoreAPI): LibraryPageBlock[] {
  return store.getState().draft.meta.layout?.blocks ?? DefaultLayout.blocks;
}

function getBlock<K extends LibraryPageBlock["type"]>(
  store: NodeStoreAPI,
  type: K,
): Extract<LibraryPageBlock, { type: K }> {
  const block = getDraftBlocks(store).find(
    (candidate): candidate is Extract<LibraryPageBlock, { type: K }> =>
      candidate.type === type,
  );

  if (!block) {
    throw new Error(
      `The ${type} block is not present. Add it with library_page_block_add before updating it.`,
    );
  }

  return block;
}

async function saveLayout(
  store: NodeStoreAPI,
  waitForSave: WaitForSave,
  message: string,
) {
  const layout = getDraftLayout(store);
  await waitForSave(store, layout);

  return { message, layout };
}

export async function waitForLibraryPageLayoutSave(
  store: NodeStoreAPI,
  expected: LibraryPageLayout,
): Promise<void> {
  if (
    deepEqual(toWebMCPLayout(store.getState().original.meta.layout), expected)
  ) {
    return;
  }

  await new Promise<void>((resolve, reject) => {
    let settled = false;

    const finish = (error?: Error) => {
      if (settled) return;
      settled = true;
      clearTimeout(timeout);
      unsubscribe();
      error ? reject(error) : resolve();
    };

    const unsubscribe = store.subscribe((state) => {
      if (deepEqual(toWebMCPLayout(state.original.meta.layout), expected)) {
        finish();
      }
    });

    const timeout = setTimeout(() => {
      finish(
        new Error(
          "The Library page layout changed locally but was not saved within 10 seconds.",
        ),
      );
    }, AUTOSAVE_TIMEOUT_MS);
  });
}
