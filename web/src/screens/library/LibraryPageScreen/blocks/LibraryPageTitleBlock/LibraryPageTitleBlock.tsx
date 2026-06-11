import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { Heading } from "@/components/ui/heading";
import { HeadingInput } from "@/components/ui/heading-input";
import { LStack, WStack } from "@/styled-system/jsx";

import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { useLibraryPageTitleBlock } from "./useLibraryPageTitleBlock";

export function LibraryPageTitleBlock() {
  const name = useWatch((s) => s.draft.name);
  const { editing } = useEditState();

  if (editing) {
    return <LibraryPageTitleBlockEditing />;
  }

  return (
    <Heading fontSize="heading.2" fontWeight="bold">
      {name || "(untitled)"}
    </Heading>
  );
}

function LibraryPageTitleBlockEditing() {
  const {
    defaultValue,
    isTitleSuggestEnabled,
    titleInputKey,
    isLoading,
    handleSuggest,
    handleChange,
  } = useLibraryPageTitleBlock();

  function handleChangeAndReset(value: string) {
    handleChange(value);
  }

  return (
    <LStack gap="2">
      <LStack minW="0">
        <WStack alignItems="end">
          <HeadingInput
            key={`title:${titleInputKey}`}
            id="name-input"
            size={"2xl" as any}
            fontWeight="bold"
            placeholder="Name..."
            onValueChange={handleChangeAndReset}
            defaultValue={defaultValue}
          />
          {isTitleSuggestEnabled && (
            <IntelligenceAction
              title="Suggest a title for this page"
              onClick={handleSuggest}
              variant="subtle"
              h="full"
              minH="8"
              loading={isLoading}
            />
          )}
        </WStack>
      </LStack>
    </LStack>
  );
}
