import {
  Button,
  FormErrorMessage,
  HStack,
  Input,
  VStack,
} from "@chakra-ui/react";
import { FormEvent } from "react";
import { Controller } from "react-hook-form";
import { Editor } from "src/components/Editor";
import { CategorySelect } from "./components/CategorySelect/CategorySelect";
import { useComposeScreen } from "./useComposeScreen";

export function ComposeScreen() {
  const { onSubmit, handleSubmit, control, isValid, errors, isSubmitting } =
    useComposeScreen();

  return (
    <VStack
      alignItems="start" //
      gap={2}
      w="full"
      py={5}
    >
      <VStack
        as="form"
        onSubmit={handleSubmit(onSubmit)}
        alignItems="start"
        w="full"
        gap={2}
      >
        <HStack width="full" justifyContent="end">
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

          <HStack>
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

        <Controller
          render={({ field }) => <Editor onChange={field.onChange} />}
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
