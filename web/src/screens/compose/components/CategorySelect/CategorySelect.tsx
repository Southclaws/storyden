import { Select, SelectProps } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { useCategorySelect } from "./useCategorySelect";

const mapCategories = map((c: Category) => (
  <option key={c.id} value={c.id}>
    {c.name}
  </option>
));

export function CategorySelect(props: SelectProps) {
  const { categories, error } = useCategorySelect();
  if (!categories) return <Unready {...error} />;

  return (
    <Select {...props} defaultValue={categories[0]?.id} w="max-content">
      {mapCategories(categories)}
    </Select>
  );
}
