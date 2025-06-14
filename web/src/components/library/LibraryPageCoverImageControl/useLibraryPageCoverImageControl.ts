import { Asset } from "@/api/openapi-schema";
import { useLibraryPageContext } from "@/screens/library/LibraryPageScreen/Context";

export function useLibraryPageCoverImageControl() {
  const { store } = useLibraryPageContext();
  const node = store.getState().draft;
  const { setPrimaryImage, removePrimaryImage } = store.getState();

  const hasCoverImage = Boolean(node.primary_image);

  async function handleUploadCoverImage(asset: Asset) {
    setPrimaryImage({
      asset,
      isReplacement: true,
    });
  }

  async function handleRemoveCoverImage() {
    removePrimaryImage();
  }

  return {
    hasCoverImage,
    handlers: {
      handleUploadCoverImage,
      handleRemoveCoverImage,
    },
  };
}
