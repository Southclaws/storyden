import { NumberInputValueChangeDetails } from "@ark-ui/react";
import { ComponentProps } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { NumberInput } from "@/components/ui/number-input";

export type NumberInputFieldProps<T extends FieldValues> = Omit<
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

export function NumberInputField<T extends FieldValues>({
  inputProps,
  min,
  max,
  step = 1,
  scrubber,
  formatOptions,
  ...controllerProps
}: NumberInputFieldProps<T>) {
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
