import { FormEvent } from "react";
import { Controller } from "react-hook-form";

import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";
import { TitleInput as TitleInputComponent } from "src/theme/components/TitleInput";

import { useTitleInput } from "./useTitleInput";

export function TitleInput() {
  const { control, fieldError } = useTitleInput();

  return (
    <>
      <FormControl>
        <Controller
          render={({ field: { onChange, ...field }, formState }) => {
            function onInput(e: FormEvent<HTMLElement>) {
              // NOTE: not sure which event type to use here...
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              onChange((e.target as any).textContent);
            }

            return (
              <TitleInputComponent
                id="title-input"
                placeholder="Thread title..."
                onInput={onInput}
                {...field}
              >
                {formState.defaultValues?.["title"]}
              </TitleInputComponent>
            );
          }}
          control={control}
          name="title"
        />

        <FormErrorText>{fieldError?.message?.toString()}</FormErrorText>
      </FormControl>
    </>
  );
}
