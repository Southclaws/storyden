import { handle } from "@/api/client";
import { nodeAddAsset, nodeRemoveAsset } from "@/api/openapi-client/nodes";
import { Asset, Identifier, NodeWithChildren } from "@/api/openapi-schema";

export function useAssets(nodeID: Identifier) {
  async function handleAssetUpload(asset: Asset) {
    await handle(async () => {
      await nodeAddAsset(nodeID, asset.id);
    });
  }

  async function handleAssetRemove(asset: Asset) {
    await handle(async () => {
      await nodeRemoveAsset(nodeID, asset.id);
    });
  }

  return {
    handleAssetUpload,
    handleAssetRemove,
  };
}
