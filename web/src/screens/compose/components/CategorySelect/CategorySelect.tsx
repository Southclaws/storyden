import { map } from "lodash/fp";
import { Controller } from "react-hook-form";

import { Category } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";

import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { styled } from "@/styled-system/jsx";

import { useCategorySelect } from "./useCategorySelect";

const mapCategories = map((c: Category) => (
  <option key={c.id} value={c.id}>
    {c.name}
  </option>
));

export function CategorySelect() {
  const { control, fieldError, categories, error } = useCategorySelect();

  // Case 1: Categories failed to load or are still loading.
  if (categories === undefined) {
    return <Unready error={error} />;
  }

  // Case 2: There are zero categories available, do not render.
  if (categories.length === 0) {
    return null;
  }

  return (
    <FormControl w="min">
      <Controller
        render={() => (
          <styled.select
            w="max"
            backgroundColor="white.a6"
            borderColor="black.a6"
            borderRadius="lg"
            boxShadow="xs"
            py="1"
            px="2"
          >
            {mapCategories(categories)}
          </styled.select>
        )}
        control={control}
        name="category"
      />
      <FormErrorText>{fieldError?.message}</FormErrorText>
    </FormControl>
  );
}
