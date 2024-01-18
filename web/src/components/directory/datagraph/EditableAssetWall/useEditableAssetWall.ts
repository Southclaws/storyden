import { useState } from "react";

import { Asset } from "src/api/openapi/schemas";
import { useFileUpload } from "src/components/content/FileDrop/useFileDrop";

export type Props = {
  initialAssets: Asset[];
  editing?: boolean;
  onUpload: (asset: Asset) => void;
};

export function useEditableAssetWall({ initialAssets, onUpload }: Props) {
  const { upload } = useFileUpload();
  const [assets, setAssets] = useState(initialAssets ?? []);

  async function handleFile(event: React.ChangeEvent<HTMLInputElement>) {
    if (!event.target.files) return;

    for (const file of event.target.files) {
      handleAssetUpload(await upload(file));
    }
  }

  const handleAssetUpload = async (asset: Asset) => {
    onUpload(asset);
    setAssets([...assets, asset]);
  };

  return {
    assets: assets.slice(0, 6),
    handlers: {
      handleFile,
      handleAssetUpload,
    },
  };
}
