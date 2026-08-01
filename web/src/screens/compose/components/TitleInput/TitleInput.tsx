import { Controller } from "react-hook-form";

import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { HeadingInput } from "@/components/ui/heading-input";

import { useTitleInput } from "./useTitleInput";

export function TitleInput() {
  const { control, fieldError } = useTitleInput();

  return (
    <>
      <FormControl>
        <Controller
          render={({ field: { onChange, ...field }, formState }) => {
            return (
              <HeadingInput
                id="title-input"
                placeholder="Thread title..."
                onValueChange={onChange}
                defaultValue={formState.defaultValues?.["title"]}
                {...field}
              />
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
