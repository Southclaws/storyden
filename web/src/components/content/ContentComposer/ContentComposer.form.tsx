import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { ContentComposerProps } from "../composer-props";

import { ContentComposer } from "./ContentComposer";

type FormComposeFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> &
  ContentComposerProps & {
    handleEmptyStateChange?: (isEmpty: boolean) => void;
    resetKey: string;
  };

export function FormComposeField<T extends FieldValues>({
  control,
  name,
  handleEmptyStateChange,
  resetKey,

  rules,
  shouldUnregister,
  defaultValue,
  disabled,
  initialValue,
  value,

  ...props
}: FormComposeFieldProps<T>) {
  return (
    <Controller<T>
      render={({ field: { onChange } }) => {
        function handleChange(value: string, isEmpty: boolean) {
          handleEmptyStateChange?.(isEmpty);
          onChange(value);
        }

        return (
          <ContentComposer
            onChange={handleChange}
            resetKey={resetKey}
            {...props}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}
