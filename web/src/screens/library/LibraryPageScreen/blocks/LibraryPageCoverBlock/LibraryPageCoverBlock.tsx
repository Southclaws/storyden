import Image from "next/image";
import { ReactNode, Suspense, lazy } from "react";

import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { CROP_STENCIL_HEIGHT, CROP_STENCIL_WIDTH } from "./constants";
import { useLibraryPageCoverBlock } from "./useLibraryPageCoverBlock";

const LibraryPageCoverBlockEditing = lazy(() =>
  import("./LibraryPageCoverBlockEditing").then((module) => ({
    default: module.LibraryPageCoverBlockEditing,
  })),
);

export function LibraryPageCoverBlock() {
  const { isDirectEditing } = useEditState();
  const primary_image = useWatch((s) => s.draft.primary_image);
  const primaryAssetURL = getAssetURL(primary_image?.path);

  const cover = primaryAssetURL ? (
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
  ) : null;

  return isDirectEditing ? <CoverEditingHost fallback={cover} /> : cover;
}

function CoverEditingHost({ fallback }: { fallback: ReactNode }) {
  // Subscribe before the cropper loads so link-block cover updates aren't lost.
  const editor = useLibraryPageCoverBlock();

  return (
    <Suspense fallback={fallback}>
      <LibraryPageCoverBlockEditing {...editor} />
    </Suspense>
  );
}
