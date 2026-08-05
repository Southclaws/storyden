import { createListCollection } from "@ark-ui/react";
import { useEffect, useMemo, useState } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { ChevronUpDownIcon } from "@/components/ui/icons/Chevron";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";

import * as Combobox from "./Combobox.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type ComboboxFieldItem = {
  label: string;
  value: string;
};

export type ComboboxFieldProps<
  TFieldValues extends FieldValues,
  Item extends ComboboxFieldItem,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  items: Item[];
  placeholder: string;
  ariaLabel: string;
};

export function ComboboxField<
  TFieldValues extends FieldValues,
  Item extends ComboboxFieldItem,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  items,
  placeholder,
  ariaLabel,
  ...controllerProps
}: ComboboxFieldProps<TFieldValues, Item, TName>) {
  const initialCollection = useMemo(
    () => createListCollection({ items }),
    [items],
  );
  const [collection, setCollection] = useState(initialCollection);

  useEffect(() => setCollection(initialCollection), [initialCollection]);

  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <Combobox.Root
          collection={collection}
          value={typeof field.value === "string" ? [field.value] : []}
          onValueChange={({ value }) => field.onChange(value[0] ?? undefined)}
          onInputValueChange={({ inputValue }) => {
            const query = inputValue.toLocaleLowerCase();
            setCollection(
              createListCollection({
                items: items.filter(
                  (item) =>
                    item.label.toLocaleLowerCase().includes(query) ||
                    item.value.toLocaleLowerCase().includes(query),
                ),
              }),
            );
          }}
          onOpenChange={({ open }) => {
            setCollection(initialCollection);
            if (!open) field.onBlur();
          }}
          positioning={{ sameWidth: true, fitViewport: true }}
          disabled={field.disabled}
          invalid={fieldState.invalid}
          size="sm"
        >
          <Combobox.Control>
            <Combobox.Input
              ref={field.ref}
              name={field.name}
              placeholder={placeholder}
              asChild
            >
              <Input size="sm" />
            </Combobox.Input>
            <Combobox.Trigger asChild>
              <IconButton
                type="button"
                variant="ghost"
                aria-label={ariaLabel}
                size="sm"
              >
                <ChevronUpDownIcon />
              </IconButton>
            </Combobox.Trigger>
          </Combobox.Control>

          <Combobox.Positioner>
            <Combobox.Content>
              <Combobox.List>
                <Combobox.ItemGroup>
                  {collection.items.map((item) => (
                    <Combobox.Item key={item.value} item={item}>
                      <Combobox.ItemText>
                        <styled.span lineClamp="1">{item.label}</styled.span>
                      </Combobox.ItemText>
                      <Combobox.ItemIndicator>
                        <CheckIcon />
                      </Combobox.ItemIndicator>
                    </Combobox.Item>
                  ))}
                </Combobox.ItemGroup>
              </Combobox.List>
            </Combobox.Content>
          </Combobox.Positioner>
        </Combobox.Root>
      )}
    />
  );
}
