import { createListCollection } from "@ark-ui/react";
import { useState } from "react";
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
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <ComboboxFieldControl
          ariaLabel={ariaLabel}
          disabled={field.disabled}
          fieldRef={field.ref}
          invalid={fieldState.invalid}
          items={items}
          name={field.name}
          onBlur={field.onBlur}
          onChange={field.onChange}
          placeholder={placeholder}
          value={typeof field.value === "string" ? field.value : undefined}
        />
      )}
    />
  );
}

type ComboboxFieldControlProps<Item extends ComboboxFieldItem> = {
  ariaLabel: string;
  disabled?: boolean;
  fieldRef: (instance: HTMLInputElement | null) => void;
  invalid: boolean;
  items: Item[];
  name: string;
  onBlur: () => void;
  onChange: (value: string | undefined) => void;
  placeholder: string;
  value: string | undefined;
};

function ComboboxFieldControl<Item extends ComboboxFieldItem>({
  ariaLabel,
  disabled,
  fieldRef,
  invalid,
  items,
  name,
  onBlur,
  onChange,
  placeholder,
  value,
}: ComboboxFieldControlProps<Item>) {
  const initialCollection = createListCollection({ items });
  const [query, setQuery] = useState<string | null>(null);
  const selectedInputValue =
    value === undefined ? "" : (initialCollection.stringify(value) ?? value);
  const normalizedQuery = query?.toLocaleLowerCase();
  const collection =
    normalizedQuery === undefined
      ? initialCollection
      : createListCollection({
          items: items.filter(
            (item) =>
              item.label.toLocaleLowerCase().includes(normalizedQuery) ||
              item.value.toLocaleLowerCase().includes(normalizedQuery),
          ),
        });

  return (
    <Combobox.Root
      collection={collection}
      inputValue={query ?? selectedInputValue}
      value={value === undefined ? [] : [value]}
      onValueChange={({ value: nextValue }) => {
        setQuery(null);
        onChange(nextValue[0] ?? undefined);
      }}
      onInputValueChange={({ inputValue }) => {
        setQuery(inputValue);
      }}
      onOpenChange={({ open }) => {
        if (!open) {
          setQuery(null);
          onBlur();
        }
      }}
      positioning={{ sameWidth: true, fitViewport: true }}
      disabled={disabled}
      invalid={invalid}
      size="sm"
    >
      <Combobox.Control>
        <Combobox.Input
          ref={fieldRef}
          name={name}
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
  );
}
