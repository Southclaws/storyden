import {
  Box,
  Button,
  Flex,
  FormErrorMessage,
  HStack,
  Input,
  VStack,
} from "@chakra-ui/react";
import { FormEvent } from "react";
import { Controller } from "react-hook-form";
import {
  Back,
  Bold,
  Italic,
  Save,
  Send,
  Underline,
} from "src/components/Action/Action";
import { Editor } from "src/components/Editor";
import { Toolpill } from "src/components/Toolpill/Toolpill";
import { CategorySelect } from "./components/CategorySelect/CategorySelect";
import { Props, useComposeScreen } from "./useComposeScreen";

export function ComposeScreen(props: Props) {
  const { onBack, onSave, onPublish, control, isValid, errors, isSubmitting } =
    useComposeScreen(props);

  return (
    <VStack
      alignItems="start" //
      gap={2}
      w="full"
      py={5}
    >
      <VStack
        as="form"
        onSubmit={onPublish}
        alignItems="start"
        w="full"
        gap={2}
      >
        <HStack width="full" justifyContent="space-between" alignItems="start">
          <HStack>
            <Controller
              render={({ field: { onChange, ...field } }) => {
                function onInput(e: FormEvent<HTMLElement>) {
                  // NOTE: not sure which event type to use here...
                  // eslint-disable-next-line @typescript-eslint/no-explicit-any
                  onChange((e.target as any).textContent);
                }

                return (
                  <Input
                    as="span"
                    contentEditable
                    variant="unstyled"
                    fontSize="3xl"
                    overflowWrap="break-word"
                    wordBreak="break-word"
                    fontWeight="semibold"
                    placeholder="Thread title"
                    onInput={onInput}
                    {...field}
                  />
                );
              }}
              control={control}
              name="title"
            />
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
              isLoading={isSubmitting}
            >
              Post
            </Button>
          </HStack>
        </HStack>
        <FormErrorMessage>{errors.title?.message}</FormErrorMessage>

        <HStack width="full">
          <Box>
            <Controller
              render={({ field }) => (
                <>
                  <CategorySelect {...field} />
                  <FormErrorMessage>
                    {errors.category?.message}
                  </FormErrorMessage>
                </>
              )}
              control={control}
              name="category"
            />
          </Box>

          <Flex flex="1 1 auto" gap={2} overflow="hidden">
            {/* TODO: tag select */}
          </Flex>
        </HStack>

        <Controller
          render={({ field }) => (
            <Editor onChange={field.onChange}>
              <HStack display={{ base: "none", md: "flex" }}>
                <Bold />
                <Italic />
                <Underline />
              </HStack>

              <Toolpill w="min-content" display={{ base: "flex", md: "none" }}>
                <VStack>
                  <HStack>
                    <Bold />
                    <Italic />
                    <Underline />
                  </HStack>
                  <HStack>
                    <Back onClick={onBack} />
                    <Send onClick={onPublish} />
                    <Save onClick={onSave} />
                  </HStack>
                </VStack>
              </Toolpill>
            </Editor>
          )}
          control={control}
          name="body"
        />
        <FormErrorMessage>{errors.body?.message}</FormErrorMessage>
      </VStack>

      <style jsx global>{`
        [contenteditable="true"]:empty:before {
          content: "Thread title...";
          color: gray;
        }

        /* prevents the user from being able to make a (visible) newline */
        [contenteditable="true"] br {
          display: none;
        }

        .remirror-editor {
          box-shadow: none !important;
        }
      `}</style>
    </VStack>
  );
}
