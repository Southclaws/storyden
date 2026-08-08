import type { RadioGroupValueChangeDetails } from "@ark-ui/react";
import type { ReactNode } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { CardBox } from "@/components/ui/card-box";
import { Text } from "@/components/ui/text";
import { Box } from "@/styled-system/jsx";
import { hstack, lstack } from "@/styled-system/patterns";

import * as RadioGroup from "./RadioGroup.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type RadioGroupFieldItem = {
  label: ReactNode;
  value: string;
  disabled?: boolean;
};

export type RadioGroupFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  ariaLabel?: string;
  items: RadioGroupFieldItem[];
};

export function RadioGroupField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  ariaLabel,
  items,
  ...controllerProps
}: RadioGroupFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <RadioGroup.Root
          name={field.name}
          value={typeof field.value === "string" ? field.value : ""}
          onValueChange={({ value }: RadioGroupValueChangeDetails) =>
            field.onChange(value)
          }
          onBlur={field.onBlur}
          disabled={field.disabled}
          aria-invalid={fieldState.invalid || undefined}
        >
          <RadioGroup.Label srOnly>{ariaLabel ?? field.name}</RadioGroup.Label>
          {items.map((item, index) => (
            <RadioGroup.Item
              key={item.value}
              value={item.value}
              disabled={item.disabled}
            >
              <RadioGroup.ItemControl />
              <RadioGroup.ItemText>{item.label}</RadioGroup.ItemText>
              <RadioGroup.ItemHiddenInput
                ref={index === 0 ? field.ref : undefined}
              />
            </RadioGroup.Item>
          ))}
        </RadioGroup.Root>
      )}
    />
  );
}

export type RadioGroupCardFieldItem = RadioGroupFieldItem & {
  description: ReactNode;
};

export type RadioGroupCardFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  ariaLabel?: string;
  items: RadioGroupCardFieldItem[];
};

export function RadioGroupCardField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  ariaLabel,
  items,
  ...controllerProps
}: RadioGroupCardFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <RadioGroup.Root
          className={lstack()}
          name={field.name}
          value={typeof field.value === "string" ? field.value : ""}
          onValueChange={({ value }: RadioGroupValueChangeDetails) =>
            field.onChange(value)
          }
          onBlur={field.onBlur}
          disabled={field.disabled}
          aria-invalid={fieldState.invalid || undefined}
        >
          <RadioGroup.Label srOnly>{ariaLabel ?? field.name}</RadioGroup.Label>
          {items.map((item, index) => (
            <CardBox
              key={item.value}
              _hover={{ background: "interactive.emphasized.surface" }}
            >
              <RadioGroup.Item
                className={hstack({ alignItems: "start", gap: "2" })}
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
                  <Text
                    variant="supporting"
                    color={item.disabled ? "text.disabled" : "text.muted"}
                  >
                    {item.description}
                  </Text>
                </Box>

                <RadioGroup.ItemHiddenInput
                  ref={index === 0 ? field.ref : undefined}
                />
              </RadioGroup.Item>
            </CardBox>
          ))}
        </RadioGroup.Root>
      )}
    />
  );
}
