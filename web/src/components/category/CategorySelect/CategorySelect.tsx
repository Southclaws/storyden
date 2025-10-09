import { FieldValues } from "react-hook-form";

import { ErrorTooltip } from "@/components/ui/ErrorTooltip";
import {
  SelectField,
  SelectFieldProps,
} from "@/components/ui/form/SelectField";
import { HStack } from "@/styled-system/jsx";

import { useCategorySelect } from "./useCategorySelect";

export function CategorySelect<T extends FieldValues>(
  props: Omit<SelectFieldProps<T, any>, "collection" | "placeholder">,
) {
  const result = useCategorySelect();
  const { ready, collection, error } = result;

  return (
    <HStack gap="2" alignItems="center">
      <SelectField
        control={props.control}
        name={props.name}
        disabled={!ready}
        placeholder={ready ? "Category" : "Loading categories..."}
        collection={collection}
      />
      <ErrorTooltip error={error} />
    </HStack>
  );
}
