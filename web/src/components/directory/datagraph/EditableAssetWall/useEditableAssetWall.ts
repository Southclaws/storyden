import { Asset } from "src/api/openapi/schemas";
import { useFileUpload } from "src/components/content/FileDrop/useFileDrop";

export type Props = {
  assets: Asset[];
  editing?: boolean;
  onUpload: (asset: Asset) => void;
};

export function useEditableAssetWall({ onUpload }: Props) {
  const { upload } = useFileUpload();

  async function handleFile(event: React.ChangeEvent<HTMLInputElement>) {
    if (!event.target.files) return;

    for (const file of event.target.files) {
      handleAssetUpload(await upload(file));
    }
  }

  const handleAssetUpload = async (asset: Asset) => {
    onUpload(asset);
  };

  return {
    handlers: {
      handleFile,
      handleAssetUpload,
    },
  };
}
