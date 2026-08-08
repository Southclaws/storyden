import type { CheckboxCheckedChangeDetails } from "@ark-ui/react";
import type { ReactNode } from "react";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { CardBox } from "@/components/ui/card-box";
import { Box } from "@/styled-system/jsx";
import { hstack, lstack } from "@/styled-system/patterns";

import { CheckIcon } from "../icons/Check";

import {
  Checkbox,
  type CheckboxProps,
  Control,
  Group,
  HiddenInput,
  Indicator,
  Label,
  Root,
} from "./Checkbox.internal";

type BooleanFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  boolean | undefined
>;

export type CheckboxFieldProps<
  TFieldValues extends FieldValues,
  TName extends BooleanFieldName<TFieldValues> = BooleanFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> &
  Omit<
    CheckboxProps,
    | "checked"
    | "defaultChecked"
    | "disabled"
    | "inputRef"
    | "name"
    | "onBlur"
    | "onCheckedChange"
  >;

export function CheckboxField<
  TFieldValues extends FieldValues,
  TName extends BooleanFieldName<TFieldValues> = BooleanFieldName<TFieldValues>,
>({
  control,
  name,
  disabled,
  ...checkboxProps
}: CheckboxFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      control={control}
      name={name}
      disabled={disabled}
      render={({ field, fieldState }) => (
        <Checkbox
          {...checkboxProps}
          checked={field.value === true}
          disabled={field.disabled}
          name={field.name}
          onBlur={field.onBlur}
          onCheckedChange={({ checked }: CheckboxCheckedChangeDetails) => {
            field.onChange(checked === true);
          }}
          inputRef={field.ref}
          aria-invalid={fieldState.invalid || undefined}
        />
      )}
    />
  );
}

type StringArrayFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string[] | undefined
>;

export type CheckboxCardGroupFieldItem = {
  label: ReactNode;
  description: ReactNode;
  value: string;
  disabled?: boolean;
};

export type CheckboxGroupFieldItem = {
  label: ReactNode;
  value: string;
  disabled?: boolean;
};

export type CheckboxGroupFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  items: CheckboxGroupFieldItem[];
};

export function CheckboxGroupField<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
>({ items, ...controllerProps }: CheckboxGroupFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => {
        const value: string[] = Array.isArray(field.value)
          ? field.value.filter(
              (entry): entry is string => typeof entry === "string",
            )
          : [];

        return (
          <div className={lstack()}>
            {items.map((item, index) => (
              <Checkbox
                key={item.value}
                size="sm"
                name={field.name}
                checked={value.includes(item.value)}
                disabled={field.disabled || item.disabled}
                onBlur={field.onBlur}
                onCheckedChange={({ checked }) => {
                  field.onChange(
                    checked
                      ? [
                          ...value.filter((entry) => entry !== item.value),
                          item.value,
                        ]
                      : value.filter((entry) => entry !== item.value),
                  );
                }}
                inputRef={index === 0 ? field.ref : undefined}
                aria-invalid={fieldState.invalid || undefined}
              >
                {item.label}
              </Checkbox>
            ))}
          </div>
        );
      }}
    />
  );
}

export type CheckboxCardGroupFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  items: CheckboxCardGroupFieldItem[];
};

export function CheckboxCardGroupField<
  TFieldValues extends FieldValues,
  TName extends StringArrayFieldName<TFieldValues> =
    StringArrayFieldName<TFieldValues>,
>({
  items,
  ...controllerProps
}: CheckboxCardGroupFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => {
        const value: string[] = Array.isArray(field.value)
          ? field.value.filter(
              (entry): entry is string => typeof entry === "string",
            )
          : [];

        function handleChange(
          itemValue: string,
          { checked }: CheckboxCheckedChangeDetails,
        ) {
          const next = checked
            ? [...value.filter((item) => item !== itemValue), itemValue]
            : value.filter((item) => item !== itemValue);

          field.onChange(next);
        }

        return (
          <Group className={lstack()} value={value}>
            {items.map((item, index) => (
              <CardBox
                key={item.value}
                _hover={{ background: "interactive.emphasized.surface" }}
              >
                <Root
                  className={hstack({ alignItems: "start", gap: "2" })}
                  value={item.value}
                  cursor="pointer"
                  disabled={field.disabled || item.disabled}
                  onBlur={field.onBlur}
                  onCheckedChange={(details) =>
                    handleChange(item.value, details)
                  }
                  aria-invalid={fieldState.invalid || undefined}
                >
                  <Box p="0.5">
                    <Control>
                      <Indicator>
                        <CheckIcon />
                      </Indicator>
                    </Control>
                  </Box>

                  <Box>
                    <Label>{item.label}</Label>
                    <p>{item.description}</p>
                  </Box>

                  <HiddenInput
                    name={field.name}
                    ref={index === 0 ? field.ref : undefined}
                  />
                </Root>
              </CardBox>
            ))}
          </Group>
        );
      }}
    />
  );
}
