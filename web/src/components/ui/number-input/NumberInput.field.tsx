import type { NumberInputValueChangeDetails } from "@ark-ui/react";
import type { ComponentProps } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { NumberInput } from "./NumberInput.internal";

type NumberFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  number | undefined
>;

export type NumberInputFieldProps<
  TFieldValues extends FieldValues,
  TName extends NumberFieldName<TFieldValues> = NumberFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  ariaLabel: string;
  inputProps?: Omit<
    NonNullable<ComponentProps<typeof NumberInput>["inputProps"]>,
    "aria-label"
  >;
  min?: number;
  max?: number;
  step?: number;
  scrubber?: boolean;
  formatOptions?: Intl.NumberFormatOptions;
};

export function NumberInputField<
  TFieldValues extends FieldValues,
  TName extends NumberFieldName<TFieldValues> = NumberFieldName<TFieldValues>,
>({
  inputProps,
  ariaLabel,
  min,
  max,
  step = 1,
  scrubber,
  formatOptions,
  ...controllerProps
}: NumberInputFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <NumberInput
          ref={field.ref}
          name={field.name}
          value={typeof field.value === "number" ? field.value.toString() : ""}
          onValueChange={({ value }: NumberInputValueChangeDetails) => {
            field.onChange(value === "" ? undefined : Number(value));
          }}
          onBlur={field.onBlur}
          disabled={field.disabled}
          invalid={fieldState.invalid}
          min={min}
          max={max}
          step={step}
          scrubber={scrubber}
          formatOptions={formatOptions}
          inputProps={{ ...inputProps, "aria-label": ariaLabel }}
        />
      )}
    />
  );
}
