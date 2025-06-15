import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags } from "@/components/ui/combotags";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { useLibraryPageTagsBlockEditing } from "./useLibraryPageTagsBlock";

export function LibraryPageTagsBlock() {
  const { editing } = useEditState();
  const { store } = useLibraryPageContext();

  if (editing) {
    return <LibraryPageTagsBlockEditing />;
  }

  const { tags } = store.getState().draft;

  return <TagBadgeList tags={tags} />;
}

export function LibraryPageTagsBlockEditing() {
  const {
    ref,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
    handleChange,
  } = useLibraryPageTagsBlockEditing();

  return (
    <HStack w="full" gap="1" alignItems="start">
      <Combotags
        ref={ref}
        initialValue={currentTags}
        onQuery={handleQuery}
        onChange={handleChange}
      />
      {isSuggestEnabled && (
        <IntelligenceAction
          title="Suggest tags for this page"
          onClick={handleSuggestTags}
          variant="subtle"
          loading={loadingTags}
        />
      )}
    </HStack>
  );
}
