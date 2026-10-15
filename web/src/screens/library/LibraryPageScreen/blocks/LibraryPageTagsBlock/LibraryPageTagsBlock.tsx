import { Suspense, lazy } from "react";

import { TagBadgeList } from "@/components/tag/TagBadgeList";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

const LibraryPageTagsBlockEditing = lazy(() =>
  import("./LibraryPageTagsBlockEditing").then((module) => ({
    default: module.LibraryPageTagsBlockEditing,
  })),
);

export function LibraryPageTagsBlock() {
  const { isDirectEditing } = useEditState();
  const { store } = useLibraryPageContext();
  const { tags } = store.getState().draft;
  const badges = <TagBadgeList tags={tags} />;

  if (isDirectEditing) {
    return (
      <Suspense fallback={badges}>
        <LibraryPageTagsBlockEditing />
      </Suspense>
    );
  }

  return badges;
}
