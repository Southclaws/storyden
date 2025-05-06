import Image from "next/image";
import {
  FixedCropper,
  FixedCropperRef,
  ImageRestriction,
} from "react-advanced-cropper";

import { LibraryPageCoverImageControl } from "@/components/library/LibraryPageCoverImageControl/LibraryPageCoverImageControl";
import { parseNodeMetadata } from "@/lib/library/metadata";
import { css } from "@/styled-system/css";
import { Box, HStack } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";
import { CROP_STENCIL_HEIGHT, CROP_STENCIL_WIDTH } from "../../useCoverImage";
import { useEditState } from "../../useEditState";

import "react-advanced-cropper/dist/style.css";

type Props = {
  ref: React.RefObject<FixedCropperRef | null>;
};

export function LibraryPageCoverBlock({ ref }: Props) {
  const { node } = useLibraryPageContext();
  const { editing } = useEditState();

  if (editing) {
    return <LibraryPageCoverBlockEditing ref={ref} />;
  }

  const primaryAssetURL = getAssetURL(node.primary_image?.path);

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

function LibraryPageCoverBlockEditing({ ref }: Props) {
  const { node } = useLibraryPageContext();

  // This URL is used for the crop editor, it will always be the original image
  // depending on whether the current primary image has any new versions set.
  // The parent is always set to the originally uploaded image while the actual
  // `primary_image` field has whichever version is currently set as the cover.
  const primaryAssetEditingURL = getAssetURL(
    node.primary_image?.parent?.path ?? node.primary_image?.path,
  );

  const initialCoverCoordinates = parseNodeMetadata(node.meta).coverImage;

  if (!primaryAssetEditingURL) {
    return (
      <HStack w="full" justify="end">
        {/* TODO: Make this a floating overlay on top of the cropper, even if it's empty */}
        <LibraryPageCoverImageControl node={node} />
      </HStack>
    );
  }

  return (
    <Box width="full" height="64">
      <FixedCropper
        ref={ref}
        className={css({
          maxWidth: "full",
          maxHeight: "64",
          borderRadius: "md",
          // TODO: Remove black background when empty
          backgroundColor: "bg.default",
        })}
        defaultPosition={
          initialCoverCoordinates && {
            top: initialCoverCoordinates.top,
            left: initialCoverCoordinates.left,
          }
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
