import type { PinInputValueChangeDetails } from "@ark-ui/react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { PinInput, type PinInputProps } from "./PinInput.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type PinInputFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> &
  Omit<
    PinInputProps,
    "defaultValue" | "disabled" | "name" | "onBlur" | "onValueChange" | "value"
  >;

export function PinInputField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  control,
  name,
  disabled,
  ...pinInputProps
}: PinInputFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      control={control}
      name={name}
      disabled={disabled}
      render={({ field, fieldState }) => (
        <PinInput
          {...pinInputProps}
          ref={field.ref}
          name={field.name}
          value={typeof field.value === "string" ? field.value.split("") : []}
          onValueChange={({ value }: PinInputValueChangeDetails) =>
            field.onChange(value.join(""))
          }
          onBlur={field.onBlur}
          disabled={field.disabled}
          invalid={fieldState.invalid}
        />
      )}
    />
  );
}
