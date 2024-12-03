import { FieldValues } from "react-hook-form";

import {
  SelectField,
  SelectFieldProps,
} from "@/components/ui/form/SelectField";

import { useCategorySelect } from "./useCategorySelect";

export function CategorySelect<T extends FieldValues>(
  props: Omit<SelectFieldProps<T, any>, "collection" | "placeholder">,
) {
  const { ready, collection } = useCategorySelect();

  return (
    <SelectField
      control={props.control}
      name={props.name}
      disabled={!ready}
      placeholder={ready ? "Select a category" : "Loading categories..."}
      collection={collection}
    />
  );
}
