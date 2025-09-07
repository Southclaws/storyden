import { useRef } from "react";
import { FixedCropperRef } from "react-advanced-cropper";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset } from "@/api/openapi-schema";
import { CoverImage } from "@/lib/library/metadata";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";

import { useLibraryCoverEvent } from "./events";

export const CROP_STENCIL_WIDTH = 1536;
export const CROP_STENCIL_HEIGHT = 384;

export function useLibraryPageCoverBlock() {
  const { store } = useLibraryPageContext();
  const { draft, setPrimaryImage } = store.getState();

  const cropperRef = useRef<FixedCropperRef>(null);

  // Listen for external cover image updates
  useLibraryCoverEvent("library-cover:update-from-asset", (asset: Asset) => {
    handleSetCoverFromAsset(asset);
  });

  async function handleSetCoverFromAsset(asset: Asset) {
    await handle(
      async () => {
        const assetURL = getAssetURL(asset.path);
        if (!assetURL) {
          throw new Error("Asset URL could not be generated");
        }

        const response = await fetch(assetURL);
        if (!response.ok) {
          throw new Error("Failed to download asset");
        }

        const blob = await response.blob();

        const img = new Image();
        const imageLoaded = new Promise<void>((resolve, reject) => {
          img.onload = () => resolve();
          img.onerror = () => reject(new Error("Failed to load image"));
        });

        img.src = URL.createObjectURL(blob);
        await imageLoaded;

        const aspectRatio = CROP_STENCIL_WIDTH / CROP_STENCIL_HEIGHT;
        const imageAspectRatio = img.width / img.height;

        let cropX = 0;
        let cropY = 0;
        let cropWidth = img.width;
        let cropHeight = img.height;

        if (imageAspectRatio > aspectRatio) {
          cropWidth = img.height * aspectRatio;
          cropX = (img.width - cropWidth) / 2;
        } else {
          cropHeight = img.width / aspectRatio;
          cropY = (img.height - cropHeight) / 2;
        }

        const canvas = document.createElement("canvas");
        canvas.width = CROP_STENCIL_WIDTH;
        canvas.height = CROP_STENCIL_HEIGHT;
        const ctx = canvas.getContext("2d");

        if (!ctx) {
          throw new Error("Failed to get canvas context");
        }

        ctx.drawImage(
          img,
          cropX,
          cropY,
          cropWidth,
          cropHeight,
          0,
          0,
          CROP_STENCIL_WIDTH,
          CROP_STENCIL_HEIGHT,
        );

        const croppedBlob = await new Promise<Blob>((resolve, reject) => {
          canvas.toBlob((blob) => {
            if (blob == null) {
              reject(new Error("Failed to create blob from canvas"));
              return;
            }
            resolve(blob);
          });
        });

        const croppedAsset = await assetUpload(croppedBlob, {
          filename: "auto-cropped-cover",
          parent_asset_id: asset.id,
        });

        const coordinates: CoverImage = {
          left: cropX / img.width,
          top: cropY / img.height,
        };

        setPrimaryImage({
          isReplacement: false,
          config: coordinates,
          asset: croppedAsset,
        });

        URL.revokeObjectURL(img.src);
      },
      {
        errorToast: true,
      },
    );
  }

  async function uploadCroppedImageState() {
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

  async function handleInteractionEnd() {
    await handle(
      async () => {
        const result = await uploadCroppedImageState();
        if (!result) {
          throw new Error(
            "An unexpected error occurred with the image editor.",
          );
        }

        setPrimaryImage(result);
      },
      {
        errorToast: true,
      },
    );
  }

  return {
    cropperRef,
    handleInteractionEnd,
  };
}
