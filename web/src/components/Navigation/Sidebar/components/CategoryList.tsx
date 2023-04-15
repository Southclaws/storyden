import { Heading, List } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { useCategoryList } from "src/api/openapi/categories";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { NavItem } from "./NavItem";

const mapCategories = map((c: Category) => (
  <NavItem key={c.id} href={`/c/${c.name}`} w="full">
    <Heading size="sm" role="navigation" variant="ghost" w="full">
      {c.name}
    </Heading>
  </NavItem>
));

export function CategoryList() {
  const { data } = useCategoryList();
  if (!data) return <Unready />;

  // TODO: Handle errors somewhat well.
  // 1. data is cached, display but with an "offline-mode" warning?
  // 2. data is not cached, first-time render, show an error.
  // swr and pwa stuff probably has some tricks for this.

  return (
    <List margin={0} display="flex" flexDirection="column" gap={2} width="full">
      {mapCategories(data.categories)}
    </List>
  );
}
