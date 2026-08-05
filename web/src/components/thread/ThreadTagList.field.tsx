import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import type { TagReferenceList } from "@/api/openapi-schema";
import { deriveColour } from "@/utils/colour";

import { ThreadTagList } from "./ThreadTagList";

type TagFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string[] | undefined
>;

export type ThreadTagListFieldProps<
  TFieldValues extends FieldValues,
  TName extends TagFieldName<TFieldValues> = TagFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  initialTags?: TagReferenceList;
};

export function ThreadTagListField<
  TFieldValues extends FieldValues,
  TName extends TagFieldName<TFieldValues> = TagFieldName<TFieldValues>,
>({
  initialTags,
  ...controllerProps
}: ThreadTagListFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field }) => {
        const names: string[] = Array.isArray(field.value)
          ? field.value.filter(
              (name): name is string => typeof name === "string",
            )
          : (initialTags?.map((tag) => tag.name) ?? []);

        return (
          <ThreadTagList
            editing
            disabled={field.disabled}
            initialTags={names.map(
              (name) =>
                initialTags?.find((tag) => tag.name === name) ?? {
                  name,
                  colour: deriveColour(name),
                  item_count: 0,
                },
            )}
            onChange={async (tags) => field.onChange(tags)}
          />
        );
      }}
    />
  );
}
