import type { DatePickerValueChangeDetails } from "@ark-ui/react";
import { type DateValue, parseDate } from "@internationalized/date";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { DatePicker } from "./DatePicker.internal";

type DateFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type DatePickerFieldProps<
  TFieldValues extends FieldValues,
  TName extends DateFieldName<TFieldValues> = DateFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render">;

export function DatePickerField<
  TFieldValues extends FieldValues,
  TName extends DateFieldName<TFieldValues> = DateFieldName<TFieldValues>,
>(controllerProps: DatePickerFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field }) => (
        <DatePicker
          name={field.name}
          value={parseValue(field.value)}
          onValueChange={({ value }: DatePickerValueChangeDetails) => {
            const first = value[0];
            field.onChange(first ? first.toDate("utc").toISOString() : "");
          }}
          onOpenChange={({ open }) => {
            if (!open) field.onBlur();
          }}
          disabled={field.disabled}
        />
      )}
    />
  );
}

function parseValue(value: unknown): DateValue[] {
  if (typeof value !== "string" || value.trim() === "") return [];

  try {
    const datePart = value.split("T")[0];
    return datePart ? [parseDate(datePart)] : [];
  } catch {
    return [];
  }
}
