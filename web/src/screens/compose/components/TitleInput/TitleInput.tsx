import { Controller } from "react-hook-form";

import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";
import { HeadingInput } from "src/theme/components/HeadingInput";

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
