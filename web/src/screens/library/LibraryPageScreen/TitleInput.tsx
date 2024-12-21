import { Controller, useFormContext } from "react-hook-form";

import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { HeadingInput } from "@/components/ui/heading-input";

import { Form } from "./useLibraryPageScreen";

type Props = {
  imperativeValue?: string;
  onResetImperativeValue?: () => void;
};

export function TitleInput({ imperativeValue, onResetImperativeValue }: Props) {
  const { control, formState } = useFormContext<Form>();

  const fieldError = formState.errors?.["name"];

  return (
    <FormControl>
      <Controller
        render={({ field: { onChange, ...field }, formState }) => {
          function handleChangeAndReset(event: any) {
            onChange(event);
            onResetImperativeValue?.();
          }

          return (
            <HeadingInput
              id="name-input"
              size={"2xl" as any}
              fontWeight="bold"
              placeholder="Name..."
              onValueChange={handleChangeAndReset}
              defaultValue={formState.defaultValues?.["name"]}
              {...field}
              value={imperativeValue ?? field.value}
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
