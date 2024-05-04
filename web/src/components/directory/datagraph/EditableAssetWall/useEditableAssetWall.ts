import { pull } from "lodash";
import { useState } from "react";

import { Asset } from "src/api/openapi/schemas";
import { useImageUpload } from "src/components/content/useImageUpload";

export type Props = {
  initialAssets: Asset[];
  editing?: boolean;
  onUpload: (asset: Asset) => void;
  onRemove: (asset: Asset) => void;
};

export function useEditableAssetWall({
  initialAssets,
  onUpload,
  onRemove,
}: Props) {
  const { upload } = useImageUpload();
  const [assets, setAssets] = useState(initialAssets ?? []);

  async function handleFile(event: React.ChangeEvent<HTMLInputElement>) {
    if (!event.target.files) return;

    for (const file of event.target.files) {
      handleAssetUpload(await upload(file));
    }
  }

  const handleAssetUpload = async (asset: Asset) => {
    setAssets([...assets, asset]);
    onUpload(asset);
  };

  const handleAssetRemove = async (asset: Asset) => {
    setAssets(pull(assets, asset));
    onRemove(asset);
  };

  return {
    assets: assets.slice(0, 6),
    handlers: {
      handleFile,
      handleAssetUpload,
      handleAssetRemove,
    },
  };
}
