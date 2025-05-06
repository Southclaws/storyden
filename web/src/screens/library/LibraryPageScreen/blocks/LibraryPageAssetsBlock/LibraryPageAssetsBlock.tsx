import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { Box, HStack } from "@/styled-system/jsx";

import { useLibraryPageAssetsBlock } from "./useLibraryPageAssetsBlock";

export function LibraryPageAssetsBlock() {
  const { editing, shouldShow, assets, handleUpload, handleRemove } =
    useLibraryPageAssetsBlock();

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
        {assets.map((a, i) => (
          // Sizing for next/image is measured in px, size tokens are basically
          // 4X, so size token 20 used above is equal to 80px, so we pass 80 here.
          <Box
            key={a.id}
            position="relative"
            scrollSnapAlign="start"
            scrollSnapStop="always"
          >
            <AssetThumbnail
              asset={a}
              set={assets}
              setIndex={i}
              width={80}
              height={80}
            />
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
