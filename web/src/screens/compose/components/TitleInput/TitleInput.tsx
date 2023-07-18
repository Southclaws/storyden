import { FormControl, FormErrorMessage, Input } from "@chakra-ui/react";
import { FormEvent } from "react";
import { Controller } from "react-hook-form";

import { useTitleInput } from "./useTitleInput";

export function TitleInput() {
  const { control, fieldError } = useTitleInput();

  return (
    <>
      <FormControl width="full" isInvalid={!!fieldError}>
        <Controller
          render={({ field: { onChange, ...field }, formState }) => {
            function onInput(e: FormEvent<HTMLElement>) {
              // NOTE: not sure which event type to use here...
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              onChange((e.target as any).textContent);
            }

            return (
              <Input
                id="title-input"
                as="span"
                display="inline-block"
                contentEditable
                //
                // NOTE: We're doing a bit of a hack here in order to make this
                // field look nice and behave like the Substack title editor.
                //
                // More info:
                //
                // https://medium.com/programming-essentials/good-to-know-about-the-state-management-of-a-contenteditable-element-in-react-adb4f933df12
                //
                suppressContentEditableWarning
                variant="unstyled"
                width="full"
                fontSize="3xl"
                overflowWrap="break-word"
                wordBreak="break-word"
                fontWeight="semibold"
                placeholder="Thread title"
                onInput={onInput}
                {...field}
              >
                {formState.defaultValues?.["title"]}
              </Input>
            );
          }}
          control={control}
          name="title"
        />

        <FormErrorMessage>{fieldError?.message?.toString()}</FormErrorMessage>
      </FormControl>

      <style jsx global>{`
        /* Sets placeholder text for the title input. */
        #title-input[contenteditable="true"]:empty:before {
          content: "Thread title...";
          opacity: 0.5;
          color: var(--chakra-colors-chakra-placeholder-color);
        }
      `}</style>
    </>
  );
}
