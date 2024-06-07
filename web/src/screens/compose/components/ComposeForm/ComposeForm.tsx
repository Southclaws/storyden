import { FormProvider } from "react-hook-form";

import { BodyInput } from "../BodyInput/BodyInput";
import { CategorySelect } from "../CategorySelect/CategorySelect";
import { TitleInput } from "../TitleInput/TitleInput";

import { Button } from "@/components/ui/button";
import { HStack, styled } from "@/styled-system/jsx";

import { Props, useComposeForm } from "./useComposeForm";

export function ComposeForm(props: Props) {
  const { formContext, onPublish, onSave, onAssetUpload } =
    useComposeForm(props);

  return (
    <styled.form
      display="flex"
      flexDir="column"
      alignItems="start"
      w="full"
      h="full"
      gap="2"
      onSubmit={onPublish}
    >
      <FormProvider {...formContext}>
        <HStack width="full" justifyContent="space-between" alignItems="start">
          <HStack width="full">
            <TitleInput />
          </HStack>

          <HStack
            display={{ base: "none", md: "flex" }}
            flexGrow="1"
            flexShrink="0"
            maxWidth="min"
            flexDir={{ base: "column-reverse", md: "row" }}
            gap="2"
            alignItems="end"
          >
            <Button
              variant="ghost"
              size="xs"
              disabled={!formContext.formState.isValid}
              onClick={onSave}
            >
              Save
            </Button>

            <Button
              variant="subtle"
              size="xs"
              type="submit"
              disabled={!formContext.formState.isValid}
              // isLoading={formContext.formState.isSubmitting}
            >
              Post
            </Button>
          </HStack>
        </HStack>

        <HStack width="full">
          <CategorySelect />
        </HStack>

        <BodyInput onAssetUpload={onAssetUpload} />
      </FormProvider>
    </styled.form>
  );
}
