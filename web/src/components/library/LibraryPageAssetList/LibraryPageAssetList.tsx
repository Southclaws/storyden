import { parseAsBoolean, useQueryState } from "nuqs";

import { handle } from "@/api/client";
import { Asset, NodeWithChildren } from "@/api/openapi-schema";
import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { useLibraryMutation } from "@/lib/library/library";
import { Box, HStack } from "@/styled-system/jsx";

export type Props = {
  node: NodeWithChildren;
};

export function LibraryPageAssetList(props: Props) {
  const [editing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const { assets } = props.node;
  const isEmpty = assets.length === 0;
  const shouldShow = editing || !isEmpty;

  const { revalidate, addAsset, removeAsset } = useLibraryMutation(props.node);

  async function handleUpload(a: Asset) {
    await handle(
      async () => {
        await addAsset(props.node.slug, a);
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
        await removeAsset(props.node.slug, a.id);
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

  if (!shouldShow) {
    return null;
  }

  return (
    <HStack
      w="full"
      overflowX="scroll"
      overflowY="hidden"
      mb="-scrollGutter"
      scrollSnapType="x"
      scrollSnapStrictness="mandatory"
    >
      <HStack w="full" h="20" maxW="full">
        {assets.map((a) => (
          // Sizing for next/image is measured in px, size tokens are basically
          // 4X, so size token 20 used above is equal to 80px, so we pass 80 here.
          <Box
            key={a.id}
            position="relative"
            scrollSnapAlign="start"
            scrollSnapStop="always"
          >
            <AssetThumbnail asset={a} width={80} height={80} />
            {editing && (
              <IconButton
                type="button"
                position="absolute"
                top="1"
                right="1"
                colorPalette="red"
                variant="subtle"
                size="xs"
                w="5"
                h="5"
                minW="5"
                title="Remove media from page"
                onClick={() => handleRemove(a)}
              >
                <DeleteIcon />
              </IconButton>
            )}
          </Box>
        ))}

        {editing && (
          <AssetUploadAction
            title="Upload media"
            operation="add"
            onFinish={handleUpload}
            hideLabel
          >
            <Button
              w="20"
              h="20"
              minW="20"
              minH="20"
              type="button"
              variant="outline"
            >
              <MediaAddIcon />
            </Button>
          </AssetUploadAction>
        )}
      </HStack>
    </HStack>
  );
}
