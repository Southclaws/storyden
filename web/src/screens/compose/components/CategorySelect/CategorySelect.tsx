import { map } from "lodash/fp";
import { Controller } from "react-hook-form";

import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/site/Unready";
import { FormControl } from "src/theme/components/FormControl";
import { FormErrorText } from "src/theme/components/FormErrorText";

import { styled } from "@/styled-system/jsx";

import { useCategorySelect } from "./useCategorySelect";

const mapCategories = map((c: Category) => (
  <option key={c.id} value={c.id}>
    {c.name}
  </option>
));

export function CategorySelect() {
  const { control, fieldError, categories, error } = useCategorySelect();

  if (!categories) return <Unready {...error} />;

  return (
    <FormControl>
      <Controller
        render={() => (
          <styled.select
            w="max"
            backgroundColor="whiteAlpha.600"
            borderColor="blackAlpha.50"
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
