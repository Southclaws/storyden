import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { FormControl } from "@/components/ui/form-control";

import { HueSelector } from "./HueSelector.internal";

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type HueSelectorFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  onUpdate: (value: string) => void;
};

export function HueSelectorField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  onUpdate,
  ...controllerProps
}: HueSelectorFieldProps<TFieldValues, TName>) {
  return (
    <FormControl>
      <Controller
        {...controllerProps}
        render={({ field }) => (
          <HueSelector
            onChange={field.onChange}
            onUpdate={onUpdate}
            value={typeof field.value === "string" ? field.value : ""}
          />
        )}
      />
    </FormControl>
  );
}
