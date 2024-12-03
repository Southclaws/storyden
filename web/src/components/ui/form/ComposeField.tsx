import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { ContentComposerProps } from "@/components/content/ContentComposer/useContentComposer";

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render"> &
  ContentComposerProps & {
    handleEmptyStateChange?: (isEmpty: boolean) => void;
    resetKey: string;
  };

export function ComposeField<T extends FieldValues>({
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
}: Props<T>) {
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
