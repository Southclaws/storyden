import { Controller, useFormContext } from "react-hook-form";

import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Heading } from "@/components/ui/heading";
import { HeadingInput } from "@/components/ui/heading-input";
import { LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { Form } from "../../form";
import { useEditState } from "../../useEditState";

import { useLibraryPageTitleBlock } from "./useLibraryPageTitleBlock";

export function LibraryPageTitleBlock() {
  const { node } = useLibraryPageContext();
  const { editing } = useEditState();

  if (editing) {
    return <LibraryPageTitleBlockEditing />;
  }

  return (
    <Heading fontSize="heading.2" fontWeight="bold">
      {node.name || "(untitled)"}
    </Heading>
  );
}

function LibraryPageTitleBlockEditing() {
  const {
    isTitleSuggestEnabled,
    value,
    isLoading,
    handleReset,
    handleSuggest,
  } = useLibraryPageTitleBlock();

  return (
    <LStack gap="2">
      <LStack minW="0">
        <WStack alignItems="end">
          <TitleInput
            imperativeValue={value}
            onResetImperativeValue={handleReset}
          />
          {isTitleSuggestEnabled && (
            <IntelligenceAction
              title="Suggest a title for this page"
              onClick={handleSuggest}
              variant="subtle"
              h="full"
              loading={isLoading}
            />
          )}
        </WStack>
      </LStack>
    </LStack>
  );
}

type Props = {
  imperativeValue?: string;
  onResetImperativeValue?: () => void;
};

export function TitleInput({ imperativeValue, onResetImperativeValue }: Props) {
  const { control, formState } = useFormContext<Form>();

  const fieldError = formState.errors?.["name"];

  return (
    <FormControl>
      <Controller
        render={({ field: { onChange, ...field }, formState }) => {
          function handleChangeAndReset(event: any) {
            onChange(event);
            onResetImperativeValue?.();
          }

          return (
            <HeadingInput
              id="name-input"
              size={"2xl" as any}
              fontWeight="bold"
              placeholder="Name..."
              onValueChange={handleChangeAndReset}
              defaultValue={formState.defaultValues?.["name"]}
              {...field}
              value={imperativeValue ?? field.value}
            />
          );
        }}
        control={control}
        name="name"
      />

      <FormErrorText>{fieldError?.message?.toString()}</FormErrorText>
    </FormControl>
  );
}
