import { handle } from "@/api/client";
import { nodeAddAsset, nodeRemoveAsset } from "@/api/openapi-client/nodes";
import { Asset, NodeWithChildren } from "@/api/openapi-schema";

export function useAssets(node: NodeWithChildren) {
  async function handleAssetUpload(asset: Asset) {
    await handle(async () => {
      await nodeAddAsset(node.id, asset.id);
    });
  }

  async function handleAssetRemove(asset: Asset) {
    await handle(async () => {
      await nodeRemoveAsset(node.id, asset.id);
    });
  }

  return {
    handleAssetUpload,
    handleAssetRemove,
  };
}
