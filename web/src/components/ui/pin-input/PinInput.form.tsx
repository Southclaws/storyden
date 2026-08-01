import { PinInputValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { PinInput } from "./PinInput.internal";

type FormPinInputFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  length?: number;
};

export function FormPinInputField<T extends FieldValues>({
  length,
  ...controllerProps
}: FormPinInputFieldProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ formState, field }) => {
        const defaultValue = formState.defaultValues![controllerProps.name];

        function handleChange({ value }: PinInputValueChangeDetails) {
          field.onChange(value.join(""));
        }

        return (
          <PinInput
            length={length}
            defaultValue={defaultValue}
            onValueChange={handleChange}
          />
        );
      }}
    />
  );
}
