import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";
import { MediaRemoveIcon } from "@/components/ui/icons/Media";

import { useLibraryPageCoverImageControl } from "./useLibraryPageCoverImageControl";

export function LibraryPageCoverImageControl() {
  const {
    hasCoverImage,
    handlers: { handleUploadCoverImage, handleRemoveCoverImage },
  } = useLibraryPageCoverImageControl();

  // TODO: When we add more options to this editing toolbar, group these up.
  return (
    <>
      {hasCoverImage && (
        <Button
          type="button"
          size="xs"
          variant="outline"
          onClick={handleRemoveCoverImage}
        >
          <MediaRemoveIcon />
          remove cover
        </Button>
      )}

      <AssetUploadAction
        operation={hasCoverImage ? "update" : "add"}
        onFinish={handleUploadCoverImage}
        accept={["image/png", "image/jpeg", "image/gif"]}
      />
    </>
  );
}
