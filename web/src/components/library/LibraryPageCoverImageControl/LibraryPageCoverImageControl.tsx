import { ImageMinusIcon } from "lucide-react";

import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { Button } from "@/components/ui/button";

import {
  Props,
  useLibraryPageCoverImageControl,
} from "./useLibraryPageCoverImageControl";

export function LibraryPageCoverImageControl(props: Props) {
  const {
    hasCoverImage,
    handlers: { handleUploadCoverImage, handleRemoveCoverImage },
  } = useLibraryPageCoverImageControl(props);

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
          <ImageMinusIcon />
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
