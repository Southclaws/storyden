import { PinInputValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { PinInput } from "../pin-input";

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render"> & {
  length?: number;
};

export function PinInputField<T extends FieldValues>({
  length,
  ...controllerProps
}: Props<T>) {
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
