import { type DateValue, parseDate } from "@internationalized/date";
import { DatePickerValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { DatePicker } from "../date-picker";

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render">;

export function DatePickerInputField<T extends FieldValues>({
  ...controllerProps
}: Props<T>) {
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
          if (typeof value !== "string" || value.trim() === "") return undefined;

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
          <DatePicker
            value={fieldValue ?? []}
            onValueChange={handleChange}
          />
        );
      }}
    />
  );
}
