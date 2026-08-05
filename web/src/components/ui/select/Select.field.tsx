import type { ListCollection, SelectValueChangeDetails } from "@ark-ui/react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import { HStack, styled } from "@/styled-system/jsx";

import * as Select from "./Select.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type SelectFieldItem = {
  label: string;
  value: string;
  icon?: string;
};

export type SelectFieldProps<
  TFieldValues extends FieldValues,
  Item extends SelectFieldItem,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  collection: ListCollection<Item>;
  placeholder: string;
  ariaLabel?: string;
  positioning?: Select.RootProps["positioning"];
  size?: Select.RootProps["size"];
  variant?: Select.RootProps["variant"];
};

export function SelectField<
  TFieldValues extends FieldValues,
  Item extends SelectFieldItem,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  collection,
  placeholder,
  ariaLabel = placeholder,
  positioning,
  size = "sm",
  variant,
  ...controllerProps
}: SelectFieldProps<TFieldValues, Item, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <Select.Root
          minW="max"
          name={field.name}
          size={size}
          variant={variant}
          value={typeof field.value === "string" ? [field.value] : []}
          collection={collection}
          positioning={{ sameWidth: false, ...positioning }}
          onValueChange={({ value }: SelectValueChangeDetails) => {
            field.onChange(value[0] ?? undefined);
          }}
          onOpenChange={({ open }) => {
            if (!open) field.onBlur();
          }}
          disabled={field.disabled}
          invalid={fieldState.invalid}
        >
          <Select.Label srOnly>{ariaLabel}</Select.Label>
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
                      <span>{item.label}</span>
                    </HStack>
                  </Select.ItemText>
                  <Select.ItemIndicator>
                    <CheckIcon />
                  </Select.ItemIndicator>
                </Select.Item>
              ))}
            </Select.Content>
          </Select.Positioner>
          <Select.HiddenSelect ref={field.ref} />
        </Select.Root>
      )}
    />
  );
}
