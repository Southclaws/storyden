import { DatePickerValueChangeDetails } from "@ark-ui/react";
import { type DateValue, parseDate } from "@internationalized/date";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { DatePicker } from "./DatePicker.internal";

type FormDatePickerFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
>;

export function FormDatePickerField<T extends FieldValues>({
  ...controllerProps
}: FormDatePickerFieldProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ field }) => {
        function handleChange({ value }: DatePickerValueChangeDetails) {
          const first = value[0];
          if (!first) {
            field.onChange("");
            return;
          }

          const changeValue = first.toDate("utc").toISOString();

          field.onChange(changeValue);
        }

        function parseValue(value: unknown): DateValue[] | undefined {
          if (typeof value !== "string" || value.trim() === "")
            return undefined;

          try {
            const datePart = value.split("T")[0];
            if (!datePart) return undefined;
            return [parseDate(datePart)];
          } catch {
            return undefined;
          }
        }
        const fieldValue = parseValue(field.value);

        return (
          <DatePicker value={fieldValue ?? []} onValueChange={handleChange} />
        );
      }}
    />
  );
}
