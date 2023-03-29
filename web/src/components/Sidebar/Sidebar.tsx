import { Box, Button, List } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Anchor } from "../site/Anchor";
import { Unready } from "../Unready";
import { useSidebar } from "./useSidebar";

const mapCategories = (selected?: string) =>
  map((c: Category) => (
    <Anchor key={c.id} href={`/c/${c.name}`} w="full">
      <Button bgColor={c.name === selected ? "blackAlpha.200" : ""} w="full">
        {c.name}
      </Button>
    </Anchor>
  ));

export function Sidebar() {
  const { error, categories, category } = useSidebar();

  if (error) return <Unready {...error} />;

  return (
    <Box as="nav">
      <List margin={0}>{mapCategories(category)(categories)}</List>
    </Box>
  );
}
