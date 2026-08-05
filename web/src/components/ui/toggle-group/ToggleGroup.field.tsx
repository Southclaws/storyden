import { Portal, type ToggleGroupValueChangeDetails } from "@ark-ui/react";
import type { ReactNode } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import * as Tooltip from "@/components/ui/tooltip";
import { HStack } from "@/styled-system/jsx";

import * as ToggleGroup from "./ToggleGroup.internal";

type StringArrayFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string[] | undefined
>;

export type ToggleGroupFieldItem = {
  label: ReactNode;
  icon?: ReactNode;
  description: string;
  value: string;
};

export type ToggleGroupFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  items: ToggleGroupFieldItem[];
};

export function ToggleGroupField<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
>({ items, ...controllerProps }: ToggleGroupFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <ToggleGroup.Root
          multiple
          size="sm"
          value={
            Array.isArray(field.value)
              ? field.value.filter(
                  (entry): entry is string => typeof entry === "string",
                )
              : []
          }
          onValueChange={({ value }: ToggleGroupValueChangeDetails) =>
            field.onChange(value)
          }
          onBlur={field.onBlur}
          disabled={field.disabled}
          aria-invalid={fieldState.invalid || undefined}
        >
          {items.map((item) => (
            <ToggleGroup.Item
              key={item.value}
              value={item.value}
              aria-label={item.description}
            >
              <Tooltip.Root lazyMount openDelay={0}>
                <Tooltip.Trigger asChild>
                  <HStack gap="1">
                    {item.icon} {item.label}
                  </HStack>
                </Tooltip.Trigger>
                <Portal>
                  <Tooltip.Positioner>
                    <Tooltip.Content>{item.description}</Tooltip.Content>
                  </Tooltip.Positioner>
                </Portal>
              </Tooltip.Root>
            </ToggleGroup.Item>
          ))}
        </ToggleGroup.Root>
      )}
    />
  );
}
