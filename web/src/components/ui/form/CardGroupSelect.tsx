import { CheckboxCheckedChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import * as Checkbox from "@/components/ui/checkbox";
import { Box, CardBox } from "@/styled-system/jsx";
import { hstack, lstack } from "@/styled-system/patterns";

import { CheckIcon } from "../icons/Check";

type CollectionItem = {
  label: string;
  description: string;
  value: string;
};

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render"> & {
  items: CollectionItem[];
};

export function CardGroupSelect<T extends FieldValues>({
  items,
  ...props
}: Props<T>) {
  const defaultValue =
    (props.control?._defaultValues[props.name] as string[]) ?? [];

  return (
    <Controller<T>
      {...props}
      render={({ formState, field }) => {
        const defaultValue = formState.defaultValues![props.name];

        return (
          <Checkbox.Group className={lstack()} defaultValue={defaultValue}>
            {items.map((item) => {
              function handleChange({ checked }: CheckboxCheckedChangeDetails) {
                const current = field.value;
                if (checked) {
                  const next = [...current, item.value];

                  field.onChange(next);
                } else {
                  const next = current.filter((v) => v !== item.value);

                  field.onChange(next);
                }
              }

              return (
                <CardBox
                  key={item.value}
                  _hover={{
                    background: "bg.emphasized",
                  }}
                >
                  <Checkbox.Root
                    className={hstack({
                      alignItems: "start",
                      gap: "2",
                    })}
                    value={item.value}
                    cursor="pointer"
                    onCheckedChange={handleChange}
                  >
                    <Box p="0.5">
                      <Checkbox.Control>
                        <Checkbox.Indicator>
                          <CheckIcon />
                        </Checkbox.Indicator>
                      </Checkbox.Control>
                    </Box>

                    <Box>
                      <Checkbox.Label>{item.label}</Checkbox.Label>
                      <p>{item.description}</p>
                    </Box>

                    <Checkbox.HiddenInput />
                  </Checkbox.Root>
                </CardBox>
              );
            })}
          </Checkbox.Group>
        );
      }}
    />
  );
}
