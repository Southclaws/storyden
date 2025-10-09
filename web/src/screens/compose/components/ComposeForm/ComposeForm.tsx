import { FormProvider } from "react-hook-form";

import { CategorySelect } from "@/components/category/CategorySelect/CategorySelect";
import { TagListField } from "@/components/thread/ThreadTagList";
import { Button } from "@/components/ui/button";
import { HStack, WStack, styled } from "@/styled-system/jsx";

import { BodyInput } from "../BodyInput/BodyInput";
import { TitleInput } from "../TitleInput/TitleInput";

import { Props, useComposeForm } from "./useComposeForm";

export function ComposeForm(props: Props) {
  const { form, state, handlers } = useComposeForm(props);

  return (
    <styled.form
      display="flex"
      flexDir="column"
      alignItems="start"
      w="full"
      h="full"
      gap="2"
      onSubmit={handlers.handlePublish}
    >
      <FormProvider {...form}>
        <WStack
          flexDir={{
            base: "column-reverse",
            md: "row",
          }}
          alignItems={{
            base: "end",
            md: "center",
          }}
        >
          <HStack width="full">
            <CategorySelect control={form.control} name="category" />
            <TagListField
              name="tags"
              control={form.control}
              initialTags={props.initialDraft?.tags}
            />
          </HStack>

          <HStack>
            <Button
              variant="ghost"
              size="xs"
              disabled={!form.formState.isValid || state.isSavingDraft}
              onClick={handlers.handleSaveDraft}
              loading={state.isSavingDraft}
            >
              Save draft
            </Button>

            <Button
              variant="subtle"
              size="xs"
              type="submit"
              disabled={!form.formState.isValid || state.isPublishing}
              loading={state.isPublishing}
            >
              Post
            </Button>
          </HStack>
        </WStack>

        <HStack width="full" justifyContent="space-between" alignItems="start">
          <HStack width="full">
            <TitleInput />
          </HStack>
        </HStack>

        <BodyInput onAssetUpload={handlers.handleAssetUpload} />
      </FormProvider>
    </styled.form>
  );
}
