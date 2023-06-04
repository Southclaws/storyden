import {
  FormControl,
  FormErrorMessage,
  Select,
  SelectProps,
} from "@chakra-ui/react";
import { map } from "lodash/fp";
import { Controller } from "react-hook-form";

import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";

import { useCategorySelect } from "./useCategorySelect";

const mapCategories = map((c: Category) => (
  <option key={c.id} value={c.id}>
    {c.name}
  </option>
));

export function CategorySelect({ ...props }: SelectProps) {
  const { control, fieldError, categories, error } = useCategorySelect();

  if (!categories) return <Unready {...error} />;

  return (
    <FormControl isInvalid={!!fieldError}>
      <Controller
        render={() => (
          <Select w="max-content" size="xs" borderRadius="lg" {...props}>
            {mapCategories(categories)}
          </Select>
        )}
        control={control}
        name="category"
      />
      <FormErrorMessage>{fieldError?.message}</FormErrorMessage>
    </FormControl>
  );
}
