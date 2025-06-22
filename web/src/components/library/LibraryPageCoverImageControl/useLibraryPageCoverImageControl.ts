import { Asset } from "@/api/openapi-schema";
import { useLibraryPageContext } from "@/screens/library/LibraryPageScreen/Context";
import { useWatch } from "@/screens/library/LibraryPageScreen/store";

export function useLibraryPageCoverImageControl() {
  const { store } = useLibraryPageContext();
  const { setPrimaryImage, removePrimaryImage } = store.getState();

  const hasCoverImage = useWatch((s) => s.draft.primary_image);

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
