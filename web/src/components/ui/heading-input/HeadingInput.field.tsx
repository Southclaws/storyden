import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { HeadingInput, type HeadingInputProps } from "./HeadingInput.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type HeadingInputFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> &
  Omit<
    HeadingInputProps,
    "defaultValue" | "disabled" | "name" | "onBlur" | "onValueChange" | "value"
  >;

export function HeadingInputField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  control,
  name,
  disabled,
  ...headingInputProps
}: HeadingInputFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      control={control}
      name={name}
      disabled={disabled}
      render={({ field, fieldState }) => (
        <HeadingInput
          {...headingInputProps}
          ref={field.ref}
          name={field.name}
          value={typeof field.value === "string" ? field.value : ""}
          onValueChange={field.onChange}
          onBlur={field.onBlur}
          aria-invalid={fieldState.invalid || undefined}
          role="textbox"
          contentEditable={!field.disabled}
        />
      )}
    />
  );
}
