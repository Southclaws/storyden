import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useRef } from "react";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset } from "@/api/openapi-schema";
import type { ImageCropperRef } from "@/components/ui/image-cropper";

import { useLibraryPageContext } from "../../Context";

import { useLibraryCoverEvent } from "./events";

export function useLibraryPageCoverBlock() {
  const { store, imageSaving, setImageSaving } = useLibraryPageContext();
  const cropperRef = useRef<ImageCropperRef>(null);

  // Imported and uploaded images remain whole until the author chooses a crop.
  useLibraryCoverEvent("library-cover:update-from-asset", (asset: Asset) => {
    store.getState().setPrimaryImage({ asset, isReplacement: true });
  });

  // Capture the pixels before debouncing so navigation cannot discard the ref.
  function captureCrop() {
    const crop = cropperRef.current?.getCrop();
    const primary = store.getState().draft.primary_image;
    if (!crop || !primary) return;
    return { ...crop, source: primary.parent ?? primary };
  }

  const uploadCrop = useCallback(
    async (
      snapshot: NonNullable<ReturnType<typeof captureCrop>>,
    ): Promise<boolean> => {
      const result = await handle(
        async () => {
          const { coordinates, imageSize, width, height, getBlob, source } =
            snapshot;
          const blob = await getBlob();
          const uploaded = await assetUpload(blob, {
            filename: "cropped-cover",
            parent_asset_id: source.id,
          });
          // Dimensions are known locally before the variant job completes.
          const asset = {
            ...uploaded,
            width,
            height,
          };
          const state = store.getState();
          const primary = state.draft.primary_image;
          if ((primary?.parent ?? primary)?.id !== source.id) return false;
          state.setPrimaryImage({
            asset,
            config: { left: coordinates.left, top: coordinates.top },
            isReplacement: false,
          });
          const block = store
            .getState()
            .draft.meta.layout?.blocks.find((b) => b.type === "cover");
          if (block?.type === "cover") {
            state.overwriteBlock({
              ...block,
              config: {
                ...block.config,
                crop: {
                  left: coordinates.left / imageSize.width,
                  top: coordinates.top / imageSize.height,
                  width: coordinates.width / imageSize.width,
                  height: coordinates.height / imageSize.height,
                },
              },
            });
          }
          return true;
        },
        {
          errorToast: true,
          cleanup: async () => {
            setImageSaving(false);
          },
        },
      );
      return result === true;
    },
    [store, setImageSaving],
  );

  const save = useMemo(() => debounce(uploadCrop, 500), [uploadCrop]);
  useEffect(
    () => () => {
      void save.flush();
    },
    [save],
  );

  function handleInteractionEnd() {
    const snapshot = captureCrop();
    if (!snapshot) return;
    setImageSaving(true);
    void save(snapshot);
  }

  return {
    cropperRef,
    handleInteractionEnd,
    busy: imageSaving,
  };
}
