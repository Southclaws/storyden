import { isValid } from "date-fns";
import { FormProvider } from "react-hook-form";

import { Bold } from "src/components/content/ContentComposer/controls/Bold";
import { Italic } from "src/components/content/ContentComposer/controls/Italic";
import { Back, Save, Send } from "src/components/site/Action/Action";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";
import { Button } from "src/theme/components/Button";

import { BodyInput } from "../BodyInput/BodyInput";
import { CategorySelect } from "../CategorySelect/CategorySelect";
import { LinkInput } from "../LinkInput/LinkInput";
import { TitleInput } from "../TitleInput/TitleInput";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useComposeForm } from "./useComposeForm";

export function ComposeForm(props: Props) {
  const { formContext, onBack, onPublish, onSave, onAssetUpload } =
    useComposeForm(props);

  return (
    <styled.form
      display="flex"
      flexDir="column"
      alignItems="start"
      w="full"
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
            <Button kind="secondary" disabled={!isValid} onClick={onSave}>
              Save
            </Button>

            <Button
              kind="primary"
              type="submit"
              disabled={!isValid}
              // isLoading={formContext.formState.isSubmitting}
            >
              Post
            </Button>
          </HStack>
        </HStack>

        <HStack width="full">
          <CategorySelect />
        </HStack>

        <HStack width="full">
          <LinkInput />
        </HStack>

        <BodyInput onAssetUpload={onAssetUpload}>
          <Toolpill w="min-content" display={{ base: "flex", md: "none" }}>
            <VStack>
              <HStack>
                <Bold />
                <Italic />
              </HStack>
              <HStack>
                <Back onClick={onBack} />
                <Send onClick={onPublish} />
                <Save onClick={onSave} />
              </HStack>
            </VStack>
          </Toolpill>
          <HStack display={{ base: "none", md: "flex" }}>
            <Bold />
            <Italic />
          </HStack>
        </BodyInput>
      </FormProvider>
    </styled.form>
  );
}
