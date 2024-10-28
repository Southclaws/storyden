import { handle } from "@/api/client";
import { Asset, Node } from "@/api/openapi-schema";
import { useLibraryMutation } from "@/lib/library/library";

export type Props = {
  node: Node;
};

export function useLibraryPageCoverImageControl(props: Props) {
  const { updateNode, removeNodeCoverImage, revalidate } = useLibraryMutation(
    props.node,
  );

  const hasCoverImage = Boolean(props.node.primary_image);

  async function handleUploadCoverImage(asset: Asset) {
    await handle(
      async () => {
        await updateNode(
          props.node.slug,
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
        await removeNodeCoverImage(props.node.slug);
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
