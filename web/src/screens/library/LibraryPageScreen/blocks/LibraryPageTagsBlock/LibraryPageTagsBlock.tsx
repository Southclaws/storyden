import { Controller } from "react-hook-form";

import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags } from "@/components/ui/combotags";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { useLibraryPageTagsBlockEditing } from "./useLibraryPageTagsBlock";

export function LibraryPageTagsBlock() {
  const { editing } = useEditState();
  const { node } = useLibraryPageContext();

  if (editing) {
    return <LibraryPageTagsBlockEditing />;
  }

  return <TagBadgeList tags={node.tags} />;
}

export function LibraryPageTagsBlockEditing() {
  const {
    ref,
    form,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
  } = useLibraryPageTagsBlockEditing();

  return (
    // TODO: Remove Controller here and use form state directly.
    <Controller
      name="tags"
      control={form.control}
      render={({ field }) => {
        async function handleChange(values: string[]) {
          field.onChange(values);
        }

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
      }}
    />
  );
}
