import { handle } from "@/api/client";
import { nodeAddAsset, nodeRemoveAsset } from "@/api/openapi-client/nodes";
import { Asset } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";
import { useBlock } from "../useBlock";

export function useLibraryPageAssetsBlock() {
  const { editing } = useEditState();
  const { nodeID, store } = useLibraryPageContext();
  const block = useBlock("assets");

  const { addAsset, removeAsset, overwriteBlock } = store.getState();

  const assets = useWatch((s) => s.draft.assets);

  const isEmpty = assets.length === 0;
  const shouldShow = editing || !isEmpty;

  const { revalidate } = useLibraryMutation();

  if (block === undefined) {
    throw new Error(
      "useLibraryPageAssetsBlock called without an Assets block.",
    );
  }

  async function handleUpload(a: Asset) {
    await handle(
      async () => {
        await nodeAddAsset(nodeID, a.id);
        addAsset(a);
      },
      {
        promiseToast: {
          loading: "Uploading...",
          success: "Added new media",
        },
        cleanup: async () => await revalidate(),
      },
    );
  }

  async function handleRemove(a: Asset) {
    await handle(
      async () => {
        await nodeRemoveAsset(nodeID, a.id);
        removeAsset(a);
      },
      {
        promiseToast: {
          loading: "Removing...",
          success: "Removed media",
        },
        cleanup: async () => await revalidate(),
      },
    );
  }

  const handleChangeSize = (size: number) => {
    overwriteBlock({
      type: "assets",
      config: {
        layout: block.config?.layout ?? "grid",
        gridSize: size,
      },
    });
  };

  return {
    editing,
    shouldShow,
    assets,
    config: block.config,
    handleUpload,
    handleRemove,
    handleChangeSize,
  };
}
