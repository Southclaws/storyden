import { Box, List } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "../Unready";
import { NavItem } from "./components/NavItem";
import { useSidebar } from "./useSidebar";

const mapCategories = (selected?: string) =>
  map((c: Category) => (
    <NavItem
      key={c.id}
      href={`/c/${c.name}`}
      w="full"
      selected={c.name === selected}
    >
      {c.name}
    </NavItem>
  ));

export function Sidebar() {
  const { error, categories, category } = useSidebar();

  if (error) return <Unready {...error} />;

  return (
    <Box as="nav" py={4}>
      <List margin={0} display="flex" flexDirection="column" gap={2}>
        {mapCategories(category)(categories)}
      </List>
    </Box>
  );
}
