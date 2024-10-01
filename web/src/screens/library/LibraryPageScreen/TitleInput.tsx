import { Controller, useFormContext } from "react-hook-form";

import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { HeadingInput } from "@/components/ui/heading-input";

import { Form } from "./useLibraryPageScreen";

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
              fontWeight="bold"
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
