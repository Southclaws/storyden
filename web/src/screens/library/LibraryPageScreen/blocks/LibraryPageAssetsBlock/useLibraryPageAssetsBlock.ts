import { handle } from "@/api/client";
import { Asset } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

export function useLibraryPageAssetsBlock() {
  const { editing } = useEditState();
  const { nodeID } = useLibraryPageContext();

  const assets = useWatch((s) => s.draft.assets);

  const isEmpty = assets.length === 0;
  const shouldShow = editing || !isEmpty;

  const { revalidate, addAsset, removeAsset } = useLibraryMutation();

  async function handleUpload(a: Asset) {
    await handle(
      async () => {
        await addAsset(nodeID, a);
      },
      {
        promiseToast: {
          loading: "Uploading...",
          success: "New media added",
        },
        cleanup: async () => await revalidate(),
      },
    );
  }

  async function handleRemove(a: Asset) {
    await handle(
      async () => {
        await removeAsset(nodeID, a.id);
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
