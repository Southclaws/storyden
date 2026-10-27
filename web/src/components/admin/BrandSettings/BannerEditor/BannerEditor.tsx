import { useRef, useState } from "react";

import { handle } from "@/api/client";
import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { assetUpload } from "@/api/openapi-client/assets";
import { bannerGet } from "@/api/openapi-client/misc";
import type { Asset } from "@/api/openapi-schema";
import { UnreadyBanner } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { SaveIcon } from "@/components/ui/icons/Save";
import {
  ImageCropper,
  type ImageCropperRef,
} from "@/components/ui/image-cropper";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { parseAdminSettings } from "@/lib/settings/settings";
import { css } from "@/styled-system/css";
import { LStack } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { getBannerURL } from "@/utils/icon";

export const CROP_STENCIL_WIDTH = 1200;
export const CROP_STENCIL_HEIGHT = 630;

export function BannerEditor() {
  const { data, error } = useAdminSettingsGet();
  const { updateSettings } = useSettingsMutation();
  const cropperRef = useRef<ImageCropperRef>(null);
  const fileRef = useRef<File | undefined>(undefined);
  const uploadedSourceRef = useRef<Asset | undefined>(undefined);
  const settings = data ? parseAdminSettings(data) : undefined;
  const assigned = settings?.banner_asset;
  const source = assigned?.parent ?? assigned;
  const framing = settings?.metadata.branding?.banner;
  const initialCrop =
    framing?.asset_id === assigned?.id ? framing?.crop : undefined;
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  async function handleSave() {
    setSaving(true);
    setSaved(false);
    await handle(
      async () => {
        const crop = cropperRef.current?.getCrop();
        if (!crop)
          throw new Error("Wait for the banner image to finish loading.");
        let original = uploadedSourceRef.current;
        if (!original) {
          original = fileRef.current
            ? await assetUpload(fileRef.current, {
                filename: fileRef.current.name,
              })
            : (source ??
              (await assetUpload(await bannerGet({ variant: "original" }), {
                filename: "banner-original.png",
              })));
          uploadedSourceRef.current = original;
        }
        const image = await assetUpload(await crop.getBlob(), {
          filename: "banner-crop.png",
          parent_asset_id: original.id,
        });
        const { coordinates, imageSize } = crop;
        await updateSettings({
          banner_asset_id: image.id,
          metadata: {
            branding: {
              ...settings?.metadata.branding,
              banner: {
                asset_id: image.id,
                crop: {
                  left: coordinates.left / imageSize.width,
                  top: coordinates.top / imageSize.height,
                  width: coordinates.width / imageSize.width,
                  height: coordinates.height / imageSize.height,
                },
              },
            },
          },
        });
        setSaved(true);
      },
      { cleanup: async () => setSaving(false) },
    );
  }

  if (!settings) return <UnreadyBanner error={error} />;

  return (
    <LStack gap="2">
      <Button
        type="button"
        alignSelf="start"
        onClick={handleSave}
        loading={saving}
        loadingText="Saving banner..."
      >
        <SaveIcon /> Save banner
      </Button>
      <ImageCropper
        ref={cropperRef}
        aria-label="Banner image"
        src={getAssetURL(source?.path) ?? getBannerURL()}
        defaultCoordinates={
          initialCrop
            ? (size) => ({
                left: initialCrop.left * size.width,
                top: initialCrop.top * size.height,
                width: initialCrop.width * size.width,
                height: initialCrop.height * size.height,
              })
            : undefined
        }
        onFileChange={async (file) => {
          fileRef.current = file;
          uploadedSourceRef.current = undefined;
          setSaved(false);
        }}
        width={CROP_STENCIL_WIDTH}
        height={CROP_STENCIL_HEIGHT}
        className={css({ maxWidth: "[calc(16rem * 1200 / 630)]" })}
        toolbar="expanded"
        saving={saving}
        saved={saved}
        onChange={() => setSaved(false)}
      />
    </LStack>
  );
}
