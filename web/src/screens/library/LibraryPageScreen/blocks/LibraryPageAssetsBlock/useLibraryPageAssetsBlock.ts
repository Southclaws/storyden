import { handle } from "@/api/client";
import { nodeAddAsset, nodeRemoveAsset } from "@/api/openapi-client/nodes";
import { Asset } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

export function useLibraryPageAssetsBlock() {
  const { editing } = useEditState();
  const { nodeID, store } = useLibraryPageContext();

  const { addAsset, removeAsset } = store.getState();

  const assets = useWatch((s) => s.draft.assets);

  const isEmpty = assets.length === 0;
  const shouldShow = editing || !isEmpty;

  const { revalidate } = useLibraryMutation();

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

  return {
    editing,
    shouldShow,
    assets,
    handleUpload,
    handleRemove,
  };
}
