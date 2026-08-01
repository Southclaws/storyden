import { RadioGroupValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import * as RadioGroup from "./RadioGroup.internal";

export type FormRadioGroupFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  items: FormRadioGroupFieldItem[];
};

export type FormRadioGroupFieldItem = {
  label: string;
  value: string;
};

export function FormRadioGroupField<T extends FieldValues>({
  items,
  ...controllerProps
}: FormRadioGroupFieldProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ field, formState }) => {
        const defaultValue = formState.defaultValues![
          controllerProps.name
        ] as string;

        function handleChange({ value }: RadioGroupValueChangeDetails) {
          field.onChange(value);
        }

        return (
          <RadioGroup.Root
            defaultValue={defaultValue}
            onValueChange={handleChange}
          >
            {items.map((item) => (
              <RadioGroup.Item key={item.value} value={item.value}>
                <RadioGroup.ItemControl />
                <RadioGroup.ItemText>{item.label}</RadioGroup.ItemText>
                <RadioGroup.ItemHiddenInput />
              </RadioGroup.Item>
            ))}
          </RadioGroup.Root>
        );
      }}
    />
  );
}
