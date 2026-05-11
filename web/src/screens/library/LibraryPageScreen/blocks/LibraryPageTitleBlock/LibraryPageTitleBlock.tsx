import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { Heading } from "@/components/ui/heading";
import { HeadingInput } from "@/components/ui/heading-input";
import { useI18n } from "@/i18n/provider";
import { LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { useLibraryPageTitleBlock } from "./useLibraryPageTitleBlock";

export function LibraryPageTitleBlock() {
  const { t } = useI18n();
  const { store } = useLibraryPageContext();
  const { draft } = store.getState();
  const { editing } = useEditState();

  if (editing) {
    return <LibraryPageTitleBlockEditing />;
  }

  return (
    <Heading fontSize="heading.2" fontWeight="bold">
      {draft.name || t("(untitled)")}
    </Heading>
  );
}

function LibraryPageTitleBlockEditing() {
  const { t } = useI18n();
  const {
    defaultValue,
    isTitleSuggestEnabled,
    value,
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
            id="name-input"
            size={"2xl" as any}
            fontWeight="bold"
            placeholder={t("Name...")}
            onValueChange={handleChangeAndReset}
            defaultValue={defaultValue}
            value={value}
          />
          {isTitleSuggestEnabled && (
            <IntelligenceAction
              title={t("Suggest a title for this page")}
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
