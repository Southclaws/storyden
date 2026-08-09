import type { SliderValueChangeDetails } from "@ark-ui/react";
import type { ComponentProps } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { Slider } from "./Slider.internal";

type NumberFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  number | undefined
>;

export type SliderFieldProps<
  TFieldValues extends FieldValues,
  TName extends NumberFieldName<TFieldValues> = NumberFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  min?: number;
  max?: number;
  step?: number;
  marks?: ComponentProps<typeof Slider>["marks"];
  label?: string;
  sliderDefaultValue?: number;
};

export function SliderField<
  TFieldValues extends FieldValues,
  TName extends NumberFieldName<TFieldValues> = NumberFieldName<TFieldValues>,
>({
  min = 0,
  max = 100,
  step = 1,
  marks,
  label,
  sliderDefaultValue,
  ...controllerProps
}: SliderFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <Slider
          ref={field.ref}
          name={field.name}
          value={[typeof field.value === "number" ? field.value : min]}
          onValueChange={({ value }: SliderValueChangeDetails) =>
            field.onChange(value[0])
          }
          onBlur={field.onBlur}
          disabled={field.disabled}
          invalid={fieldState.invalid}
          min={min}
          max={max}
          step={step}
          marks={marks}
          defaultValue={
            sliderDefaultValue === undefined ? undefined : [sliderDefaultValue]
          }
        >
          {label}
        </Slider>
      )}
    />
  );
}
