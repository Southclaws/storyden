import Image from "next/image";
import { FixedCropper, ImageRestriction } from "react-advanced-cropper";

import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Asset } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { parseNodeMetadata } from "@/lib/library/metadata";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import "react-advanced-cropper/dist/style.css";

import {
  CROP_STENCIL_HEIGHT,
  CROP_STENCIL_WIDTH,
  useLibraryPageCoverBlock,
} from "./useLibraryPageCoverBlock";

export function LibraryPageCoverBlock() {
  const { editing } = useEditState();
  const primary_image = useWatch((s) => s.draft.primary_image);

  if (editing) {
    return <LibraryPageCoverBlockEditing />;
  }

  const primaryAssetURL = getAssetURL(primary_image?.path);

  if (!primaryAssetURL) {
    return null;
  }

  return (
    <Box height="64" width="full">
      <Image
        className={css({
          width: "full",
          height: "full",
          borderRadius: "md",
          objectFit: "cover",
          objectPosition: "center",
        })}
        src={primaryAssetURL}
        alt=""
        width={CROP_STENCIL_WIDTH}
        height={CROP_STENCIL_HEIGHT}
      />
    </Box>
  );
}

function LibraryPageCoverBlockEditing() {
  const { store } = useLibraryPageContext();
  const { setPrimaryImage } = store.getState();
  const { cropperRef, handleInteractionEnd } = useLibraryPageCoverBlock();

  const primary_image = useWatch((s) => s.draft.primary_image);
  const meta = useWatch((s) => s.draft.meta);

  async function handleUploadCoverImage(asset: Asset) {
    setPrimaryImage({
      asset,
      isReplacement: true,
    });
  }

  // This URL is used for the crop editor, it will always be the original image
  // depending on whether the current primary image has any new versions set.
  // The parent is always set to the originally uploaded image while the actual
  // `primary_image` field has whichever version is currently set as the cover.
  const primaryAssetEditingURL = getAssetURL(
    primary_image?.parent?.path ?? primary_image?.path,
  );

  const initialCoverCoordinates = parseNodeMetadata(meta).coverImage;

  if (!primaryAssetEditingURL) {
    return (
      <AssetUploadAction
        operation="add"
        onFinish={handleUploadCoverImage}
        accept={["image/png", "image/jpeg", "image/gif"]}
        w="full"
        size="xs"
        variant="subtle"
      >
        <Button
          type="button"
          w="full"
          size="xs"
          variant="subtle"
        >
          Upload cover image
        </Button>
      </AssetUploadAction>
    );
  }

  return (
    <Box width="full" height="64">
      <FixedCropper
        ref={cropperRef}
        className={css({
          maxWidth: "full",
          maxHeight: "64",
          borderRadius: "md",
          // TODO: Remove black background when empty
          backgroundColor: "bg.default",
        })}
        onInteractionEnd={handleInteractionEnd}
        defaultPosition={
          initialCoverCoordinates
            ? {
                top: initialCoverCoordinates.top,
                left: initialCoverCoordinates.left,
              }
            : undefined
        }
        backgroundWrapperProps={{
          scaleImage: false,
        }}
        stencilProps={{
          handlers: false,
          lines: false,
          movable: false,
          resizable: false,
        }}
        stencilSize={{
          width: CROP_STENCIL_WIDTH,
          height: CROP_STENCIL_HEIGHT,
        }}
        imageRestriction={ImageRestriction.stencil}
        src={primaryAssetEditingURL}
      />
    </Box>
  );
}
