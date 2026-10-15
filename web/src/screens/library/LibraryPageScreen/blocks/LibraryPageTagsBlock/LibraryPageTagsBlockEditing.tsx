import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { MultiSelectPicker } from "@/components/ui/multi-select-picker";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageTagsBlockEditing } from "./useLibraryPageTagsBlock";

export function LibraryPageTagsBlockEditing() {
  const {
    currentTagItems,
    queryResults,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
    handleChange,
  } = useLibraryPageTagsBlockEditing();

  return (
    <HStack w="full" gap="1" alignItems="start">
      <MultiSelectPicker
        value={currentTagItems}
        onChange={handleChange}
        onQuery={handleQuery}
        queryResults={queryResults}
        allowNewValues={true}
        inputPlaceholder="Add tags..."
        autoColour={true}
        size="sm"
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
