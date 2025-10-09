import { ListCollection, SelectValueChangeDetails } from "@ark-ui/react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { HStack, styled } from "@/styled-system/jsx";

export type SelectFieldProps<
  T extends FieldValues,
  Item extends ListCollectionItem,
> = Omit<ControllerProps<T>, "render"> & {
  collection: ListCollection<Item>;
  placeholder: string;
};

export type ListCollectionItem = {
  label: string;
  value: string;
  icon?: string;
};

export function SelectField<
  T extends FieldValues,
  Item extends ListCollectionItem,
>({ collection, placeholder, ...controllerProps }: SelectFieldProps<T, Item>) {
  const defaultValue = controllerProps.defaultValue as string;

  return (
    <Controller<T>
      {...controllerProps}
      render={({ field, formState, fieldState }) => {
        function handleChange({ value }: SelectValueChangeDetails) {
          const [v] = value;
          if (!v) return;
          field.onChange(v);
        }

        return (
          <Select.Root
            minW="max"
            size="sm"
            value={[field.value]}
            defaultValue={[defaultValue]}
            collection={collection}
            positioning={{ sameWidth: false }}
            onValueChange={handleChange}
            disabled={controllerProps.disabled}
          >
            <Select.Control>
              <Select.Trigger>
                <Select.ValueText placeholder={placeholder} />
                <SelectIcon />
              </Select.Trigger>
            </Select.Control>
            <Select.Positioner>
              <Select.Content>
                {collection.items.map((item) => (
                  <Select.Item key={item.value} item={item}>
                    <Select.ItemText mr="2">
                      <HStack gap="1">
                        {item.icon && (
                          <styled.span w="4">{item.icon}</styled.span>
                        )}
                        <styled.span>{item.label}</styled.span>
                      </HStack>
                    </Select.ItemText>
                    <Select.ItemIndicator>
                      <CheckIcon />
                    </Select.ItemIndicator>
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Positioner>
          </Select.Root>
        );
      }}
    />
  );
}
