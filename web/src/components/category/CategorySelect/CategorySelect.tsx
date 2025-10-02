import { FieldValues } from "react-hook-form";

import { Unready } from "@/components/site/Unready";
import {
  SelectField,
  SelectFieldProps,
} from "@/components/ui/form/SelectField";

import { useCategorySelect } from "./useCategorySelect";

export function CategorySelect<T extends FieldValues>(
  props: Omit<SelectFieldProps<T, any>, "collection" | "placeholder">,
) {
  const { ready, collection } = useCategorySelect();

  // Case 1: Categories failed to load or are still loading.
  if (!ready) {
    return <Unready />;
  }

  // Case 2: There are zero categories available, do not render.
  if (collection.size === 0) {
    return null;
  }

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
