import { SliderValueChangeDetails } from "@ark-ui/react";
import { ComponentProps } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { Slider } from "@/components/ui/slider";

export type SliderFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  min?: number;
  max?: number;
  step?: number;
  marks?: ComponentProps<typeof Slider>["marks"];
  label?: string;
  sliderDefaultValue?: number;
};

export function SliderField<T extends FieldValues>({
  min = 0,
  max = 100,
  step = 1,
  marks,
  label,
  sliderDefaultValue,
  ...controllerProps
}: SliderFieldProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ field, formState }) => {
        function handleChange(details: SliderValueChangeDetails) {
          const value = details.value[0];
          field.onChange(value);
        }

        return (
          <Slider
            name={field.name}
            value={field.value != null ? [field.value] : [min]}
            onValueChange={handleChange}
            onBlur={field.onBlur}
            disabled={controllerProps.disabled}
            min={min}
            max={max}
            step={step}
            marks={marks}
            defaultValue={
              sliderDefaultValue !== undefined
                ? [sliderDefaultValue]
                : undefined
            }
          >
            {label}
          </Slider>
        );
      }}
    />
  );
}
