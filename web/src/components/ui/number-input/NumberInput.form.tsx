import { NumberInputValueChangeDetails } from "@ark-ui/react";
import { ComponentProps } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { NumberInput } from "./NumberInput.internal";

export type FormNumberInputFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  inputProps?: ComponentProps<typeof NumberInput>["inputProps"];
  min?: number;
  max?: number;
  step?: number;
  scrubber?: boolean;
  formatOptions?: Intl.NumberFormatOptions;
};

export function FormNumberInputField<T extends FieldValues>({
  inputProps,
  min,
  max,
  step = 1,
  scrubber,
  formatOptions,
  ...controllerProps
}: FormNumberInputFieldProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ field }) => {
        function handleChange(details: NumberInputValueChangeDetails) {
          const numValue =
            details.value === "" ? undefined : Number(details.value);
          field.onChange(numValue);
        }

        return (
          <NumberInput
            name={field.name}
            value={field.value?.toString() ?? ""}
            onValueChange={handleChange}
            onBlur={field.onBlur}
            disabled={controllerProps.disabled}
            min={min}
            max={max}
            step={step}
            scrubber={scrubber}
            formatOptions={formatOptions}
            inputProps={inputProps}
          />
        );
      }}
    />
  );
}
