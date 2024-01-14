import { Controller, useFormContext } from "react-hook-form";

import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";
import { HeadingInput } from "src/theme/components/HeadingInput";

import { Form } from "./useClusterScreen";

export function TitleInput() {
  const { control, formState } = useFormContext<Form>();

  const fieldError = formState.errors?.["name"];

  return (
    <FormControl>
      <Controller
        render={({ field: { onChange, ...field }, formState }) => {
          return (
            <HeadingInput
              id="name-input"
              size={"2xl" as any}
              placeholder="Name..."
              onValueChange={onChange}
              defaultValue={formState.defaultValues?.["name"]}
              {...field}
            />
          );
        }}
        control={control}
        name="name"
      />

      <FormErrorText>{fieldError?.message?.toString()}</FormErrorText>
    </FormControl>
  );
}
