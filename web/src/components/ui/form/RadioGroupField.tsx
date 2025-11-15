import { RadioGroupValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import * as RadioGroup from "@/components/ui/radio-group";

export type RadioGroupFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  items: RadioGroupItem[];
};

export type RadioGroupItem = {
  label: string;
  value: string;
};

export function RadioGroupField<T extends FieldValues>({
  items,
  ...controllerProps
}: RadioGroupFieldProps<T>) {
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
