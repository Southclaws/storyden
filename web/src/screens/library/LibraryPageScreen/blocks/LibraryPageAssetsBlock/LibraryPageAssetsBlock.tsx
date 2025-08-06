import { CSSProperties } from "react";

import { AssetThumbnail } from "@/components/asset/AssetThumbnail";
import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { css } from "@/styled-system/css";
import { Box, Grid, GridItem, HStack } from "@/styled-system/jsx";

import { useLibraryPageAssetsBlock } from "./useLibraryPageAssetsBlock";

export function LibraryPageAssetsBlock() {
  const { editing, config, shouldShow, assets, handleUpload, handleRemove } =
    useLibraryPageAssetsBlock();

  if (!shouldShow) {
    return null;
  }

  const layout = config?.layout ?? "strip";

  switch (layout) {
    case "strip": {
      const size = config?.gridSize ?? 2;

      const style = {
        "--thumbnail-size-token": size,
        "--thumbnail-size": `calc(var(--thumbnail-size-token) * 4rem)`,
      } as CSSProperties;

      return (
        <HStack
          id="gallery-strip"
          w="full"
          overflowX="scroll"
          overflowY="hidden"
          mb="-scrollGutter"
          scrollSnapType="x"
          scrollSnapStrictness="mandatory"
          style={style}
        >
          <HStack
            // w="full"
            // maxW="full"
            height={"var(--thumbnail-size)"}
          >
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
                  className={css({
                    objectFit: "cover",
                    width: "var(--thumbnail-size)",
                    height: "var(--thumbnail-size)",
                    minWidth: "var(--thumbnail-size)",
                  })}
                  asset={a}
                  set={assets}
                  setIndex={i}
                  showDeleteButton={editing}
                  handleDelete={() => handleRemove(a)}
                />
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
                  className={css({
                    width: "var(--thumbnail-size)",
                    height: "var(--thumbnail-size)",
                  })}
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
    case "grid": {
      const size = config?.gridSize ?? 3;

      const style = {
        gridTemplateColumns: `repeat(${size}, 1fr)`,
      };

      return (
        <Grid style={style}>
          {assets.map((a, i) => (
            <GridItem
              key={a.id}
              position="relative"
              scrollSnapAlign="start"
              scrollSnapStop="always"
            >
              <AssetThumbnail
                asset={a}
                set={assets}
                setIndex={i}
                showDeleteButton={editing}
                handleDelete={() => handleRemove(a)}
              />
            </GridItem>
          ))}

          {editing && (
            <AssetUploadAction
              title="Upload media"
              operation="add"
              onFinish={handleUpload}
              hideLabel
            >
              <Button
                w="full"
                h="full"
                aspectRatio="1"
                type="button"
                variant="outline"
              >
                <MediaAddIcon />
              </Button>
            </AssetUploadAction>
          )}
        </Grid>
      );
    }
  }
}
