import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import type { ContentComposerProps } from "../composer-props";

import { ContentComposer } from "./ContentComposer";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type ContentComposerFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> &
  Omit<ContentComposerProps, "disabled" | "onChange" | "value"> & {
    onEmptyStateChange?: (isEmpty: boolean) => void;
  };

export function ContentComposerField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  initialValue,
  onEmptyStateChange,
  ...controllerProps
}: ContentComposerFieldProps<TFieldValues, TName>) {
  const {
    control,
    name,
    rules,
    shouldUnregister,
    defaultValue,
    disabled,
    ...composerProps
  } = controllerProps;

  return (
    <Controller
      control={control}
      name={name}
      rules={rules}
      shouldUnregister={shouldUnregister}
      defaultValue={defaultValue}
      disabled={disabled}
      render={({ field }) => (
        <ContentComposer
          {...composerProps}
          initialValue={
            initialValue ??
            (typeof field.value === "string" ? field.value : undefined)
          }
          disabled={field.disabled}
          onChange={(value, isEmpty) => {
            onEmptyStateChange?.(isEmpty);
            field.onChange(value);
          }}
        />
      )}
    />
  );
}
