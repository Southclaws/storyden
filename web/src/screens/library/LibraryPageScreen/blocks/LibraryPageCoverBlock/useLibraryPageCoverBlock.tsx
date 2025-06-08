import { useRef } from "react";
import { FixedCropperRef } from "react-advanced-cropper";

import { assetUpload } from "@/api/openapi-client/assets";
import { CoverImage } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";

import "react-advanced-cropper/dist/style.css";

export const CROP_STENCIL_WIDTH = 1536;
export const CROP_STENCIL_HEIGHT = 384;

export function useLibraryPageCoverBlock() {
  const { store } = useLibraryPageContext();
  const { draft } = store.getState();

  const cropperRef = useRef<FixedCropperRef>(null);

  async function handleUploadCroppedCover() {
    if (!cropperRef.current) {
      return;
    }

    const canvas = cropperRef.current.getCanvas();
    if (!canvas) {
      throw new Error("An unexpected error occurred with the image editor.");
    }

    const coordinates =
      cropperRef.current.getCoordinates() satisfies CoverImage | null;
    if (!coordinates) {
      throw new Error(
        "An unexpected error occurred with the image editor: unable to get crop coordinates.",
      );
    }

    const blob = await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (blob == null) {
          reject("An unexpected error occurred with the image editor.");
          return;
        }

        resolve(blob);
      });
    });

    if (draft.primary_image) {
      // TODO: Delete the original asset maybe?
    }

    // The cover image is determined to be a copy of an original if it has a
    // parent asset associated with it. Original assets do not have parents.
    const isCopy = draft.primary_image?.parent?.id !== undefined;

    const parent_asset_id = isCopy
      ? // If the primary image is already a copy, use the existing parent.
        draft.primary_image?.parent?.id
      : // Otherwise, use the primary image asset ID, which will result in
        // this ID becoming the parent.
        draft.primary_image?.id;

    // The result of this is that when the asset revalidates, the original
    // asset will be present in the primary image asset object. This means
    // that users will always edit the originally uploaded image and when
    // they stop editing, the saved asset will be the cropped version where
    // the original is still present within the asset's parent property.
    const asset = await assetUpload(blob, {
      // TODO: Split filename from ID in API side (Marks) use original name.
      filename: "cropped-cover",
      parent_asset_id,
    });

    return {
      isReplacement: false,
      config: coordinates,
      asset,
    };
  }

  return {
    cropperRef,
    handleUploadCroppedCover,
  };
}
