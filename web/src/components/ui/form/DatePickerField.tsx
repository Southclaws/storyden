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
      render={({ formState, field }) => {
        const defaultValue = formState.defaultValues![controllerProps.name];

        function handleChange({ value }: DatePickerValueChangeDetails) {
          const first = value[0];
          if (!first) {
            field.onChange(undefined);
            return;
          }

          const changeValue = first.toDate("utc").toISOString();

          field.onChange(changeValue);
        }

        return (
          <DatePicker
            defaultValue={defaultValue}
            onValueChange={handleChange}
          />
        );
      }}
    />
  );
}
