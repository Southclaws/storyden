import { type CSSProperties, useState } from "react";

import { assetUpload } from "@/api/openapi-client/assets";
import { Asset } from "@/api/openapi-schema";
import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";
import { ImageCropper } from "@/components/ui/image-cropper";
import { type LibraryPageBlockTypeCover } from "@/lib/library/metadata";
import { css } from "@/styled-system/css";
import { getAssetURL } from "@/utils/asset";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useBlock } from "../useBlock";

import { LibraryPageCoverBlockMenuItems } from "./LibraryPageCoverBlockMenuItems";
import type { useLibraryPageCoverBlock } from "./useLibraryPageCoverBlock";

type Props = ReturnType<typeof useLibraryPageCoverBlock>;
const frameStyles = {
  banner: css({ width: "full" }),
  image: css({
    width: "full",
    maxHeight: "[min(70svh,40rem)]",
    minHeight: "9",
  }),
  artwork: css({ width: "[min(100%,20rem)]" }),
};

export function LibraryPageCoverBlockEditing(props: Props) {
  const { store } = useLibraryPageContext();
  const primary = useWatch((s) => s.draft.primary_image);
  const block = useBlock("cover");

  async function handleUpload(asset: Asset) {
    const empty = !store.getState().draft.primary_image;
    store.getState().setPrimaryImage({ asset, isReplacement: true });
    if (empty && block && !block.config?.layout)
      store.getState().overwriteBlock({
        ...block,
        config: { ...block.config, layout: "image" },
      });
  }

  if (!primary || !block)
    return (
      <AssetUploadAction
        operation="add"
        onFinish={handleUpload}
        accept={["image/png", "image/jpeg", "image/gif", "image/webp"]}
        w="full"
      >
        <Button type="button" w="full" variant="subtle">
          Upload main image
        </Button>
      </AssetUploadAction>
    );
  const source = primary.parent ?? primary;
  return (
    <MainImageCropEditor
      key={`${source.id}:${block.config?.layout ?? "banner"}`}
      {...props}
      primary={primary}
      block={block}
      onUpload={handleUpload}
    />
  );
}

function MainImageCropEditor({
  primary,
  block,
  cropperRef,
  handleInteractionEnd,
  busy,
  onUpload,
}: Props & {
  primary: Asset;
  block: LibraryPageBlockTypeCover;
  onUpload: (asset: Asset) => Promise<void>;
}) {
  const { store, setImageSaving } = useLibraryPageContext();
  const { crop, layout = "banner" } = block.config ?? {};
  const legacyPosition = useWatch((s) => s.draft.meta.coverImage);
  // Crops replace the primary asset. The workspace must retain the original's
  // dimensions until a new original or page layout remounts this editor.
  const [source] = useState(() => primary.parent ?? primary);
  const size =
    layout === "banner"
      ? { width: 1536, height: 384 }
      : layout === "artwork"
        ? { width: 1024, height: 1024 }
        : { width: source.width || 1536, height: source.height || 384 };

  return (
    <ImageCropper
      ref={cropperRef}
      src={getAssetURL(source.path)}
      {...size}
      cropMode={layout === "image" ? "resizable" : "fixed"}
      maxOutputSize={{ width: 1920, height: 1920 }}
      cropOnUpload={false}
      toolbar="expanded"
      showHint={false}
      accept={["image/png", "image/jpeg", "image/gif", "image/webp"]}
      menuLabel="Main image options"
      className={frameStyles[layout]}
      style={{ marginInline: "auto" } as CSSProperties}
      data-main-image-layout={layout}
      saving={busy}
      defaultCoordinates={
        crop
          ? (imageSize) => ({
              left: crop.left * imageSize.width,
              top: crop.top * imageSize.height,
              width: crop.width * imageSize.width,
              height: crop.height * imageSize.height,
            })
          : undefined
      }
      defaultPosition={legacyPosition ?? undefined}
      onChange={handleInteractionEnd}
      onFileChange={async (file) => {
        setImageSaving(true);
        await assetUpload(file, { filename: file.name })
          .then(onUpload)
          .finally(() => setImageSaving(false));
      }}
      onRemove={async () => store.getState().removePrimaryImage()}
    >
      <LibraryPageCoverBlockMenuItems />
    </ImageCropper>
  );
}
