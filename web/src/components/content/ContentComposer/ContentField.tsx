import { Control, Controller, FieldValues, Path } from "react-hook-form";

import { ContentComposer } from "./ContentComposer";
import { ContentComposerProps } from "./useContentComposer";

export type Props<T extends FieldValues> = ContentComposerProps & {
  control: Control<T>;
  name: Path<T>;
};

export function ContentFormField<T extends FieldValues>({
  control,
  name,
  ...props
}: Props<T>) {
  return (
    <Controller
      render={({ field }) => (
        <ContentComposer onChange={field.onChange} {...props} />
      )}
      control={control}
      name={name}
    />
  );
}
