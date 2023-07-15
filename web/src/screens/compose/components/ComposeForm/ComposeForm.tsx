import { Button, Flex, HStack, VStack } from "@chakra-ui/react";
import { isValid } from "date-fns";
import { FormProvider } from "react-hook-form";

import { Back, Save, Send } from "src/components/Action/Action";
import { Bold } from "src/components/ContentComposer/controls/Bold";
import { Italic } from "src/components/ContentComposer/controls/Italic";
import { Toolpill } from "src/components/Toolpill/Toolpill";

import { BodyInput } from "../BodyInput/BodyInput";
import { CategorySelect } from "../CategorySelect/CategorySelect";
import { TitleInput } from "../TitleInput/TitleInput";

import { Props, useComposeForm } from "./useComposeForm";

export function ComposeForm(props: Props) {
  const { formContext, onBack, onPublish, onSave, onAssetUpload } =
    useComposeForm(props);

  return (
    <VStack as="form" onSubmit={onPublish} alignItems="start" w="full" gap={2}>
      <FormProvider {...formContext}>
        <HStack width="full" justifyContent="space-between" alignItems="start">
          <HStack width="full">
            <TitleInput />
          </HStack>

          <HStack
            display={{ base: "none", md: "flex" }}
            flex="1 0 auto"
            maxWidth="min-content"
            flexDir={{ base: "column-reverse", md: "row" }}
            gap={2}
            alignItems="end"
          >
            <Button variant="outline" isDisabled={!isValid} onClick={onSave}>
              Save
            </Button>

            <Button
              type="submit"
              isDisabled={!isValid}
              isLoading={formContext.formState.isSubmitting}
            >
              Post
            </Button>
          </HStack>
        </HStack>

        <HStack width="full">
          <CategorySelect />

          <Flex flex="1 1 auto" gap={2} overflow="hidden">
            {/* TODO: tag select */}
          </Flex>
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
    </VStack>
  );
}
