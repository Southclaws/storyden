import { handle } from "@/api/client";
import { Asset } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";
import { useLibraryPageContext } from "@/screens/library/LibraryPageScreen/Context";

export function useLibraryPageCoverImageControl() {
  const { store } = useLibraryPageContext();
  const node = store.getState().draft;

  const { updateNode, removeNodeCoverImage, revalidate } =
    useLibraryMutation(node);

  const hasCoverImage = Boolean(node.primary_image);

  async function handleUploadCoverImage(asset: Asset) {
    await handle(
      async () => {
        await updateNode(
          node.slug,
          {},
          {
            asset,
            isReplacement: true,
          },
        );
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  async function handleRemoveCoverImage() {
    await handle(
      async () => {
        await removeNodeCoverImage(node.slug);
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  return {
    hasCoverImage,
    handlers: {
      handleUploadCoverImage,
      handleRemoveCoverImage,
    },
  };
}
