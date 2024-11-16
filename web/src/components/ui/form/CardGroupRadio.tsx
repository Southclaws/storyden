import { RadioGroupValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import * as RadioGroup from "@/components/ui/radio-group";
import { Box, CardBox } from "@/styled-system/jsx";
import { hstack, lstack } from "@/styled-system/patterns";

type CollectionItem = {
  label: string;
  description: string;
  value: string;
  disabled?: boolean;
};

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render"> & {
  items: CollectionItem[];
};

export function CardGroupRadio<T extends FieldValues>({
  items,
  ...props
}: Props<T>) {
  return (
    <Controller<T>
      {...props}
      render={({ formState, field }) => {
        const defaultValue = formState.defaultValues![props.name];

        function handleChange({ value }: RadioGroupValueChangeDetails) {
          field.onChange(value);
        }

        return (
          <RadioGroup.Root
            className={lstack()}
            defaultValue={defaultValue}
            onValueChange={handleChange}
          >
            {items.map((item) => {
              return (
                <CardBox
                  key={item.value}
                  _hover={{
                    background: "bg.emphasized",
                  }}
                >
                  <RadioGroup.Item
                    className={hstack({
                      alignItems: "start",
                      gap: "2",
                    })}
                    value={item.value}
                    cursor="pointer"
                    disabled={item.disabled}
                  >
                    <Box p="0.5">
                      <RadioGroup.ItemControl>
                        <RadioGroup.Indicator />
                      </RadioGroup.ItemControl>
                    </Box>

                    <Box>
                      <RadioGroup.ItemText>{item.label}</RadioGroup.ItemText>
                      <p>{item.description}</p>
                    </Box>

                    <RadioGroup.ItemHiddenInput />
                  </RadioGroup.Item>
                </CardBox>
              );
            })}
          </RadioGroup.Root>
        );
      }}
    />
  );
}
