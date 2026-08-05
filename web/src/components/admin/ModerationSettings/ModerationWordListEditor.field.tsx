import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { ModerationWordListEditor } from "./ModerationWordListEditor";

type WordListFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string[]
>;

export type ModerationWordListEditorFieldProps<
  TFieldValues extends FieldValues,
  TName extends WordListFieldName<TFieldValues> =
    WordListFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render">;

export function ModerationWordListEditorField<
  TFieldValues extends FieldValues,
  TName extends WordListFieldName<TFieldValues> =
    WordListFieldName<TFieldValues>,
>(controllerProps: ModerationWordListEditorFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field }) => (
        <ModerationWordListEditor
          value={Array.isArray(field.value) ? field.value : []}
          onChange={field.onChange}
          disabled={field.disabled}
        />
      )}
    />
  );
}
