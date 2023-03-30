import { Heading, HStack, List, VStack } from "@chakra-ui/react";
import { HomeIcon } from "@heroicons/react/24/outline";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "../Unready";
import { NavItem } from "./components/NavItem";
import { useSidebar } from "./useSidebar";

const mapCategories = map((c: Category) => (
  <NavItem key={c.id} href={`/c/${c.name}`} w="full">
    <Heading size="sm" role="navigation" variant="ghost" w="full">
      {c.name}
    </Heading>
  </NavItem>
));

export function Sidebar() {
  const { error, categories } = useSidebar();

  if (error) return <Unready {...error} />;

  return (
    <VStack as="nav" py={4} gap={4} alignItems="start">
      <HStack gap={2}>
        <NavItem href="/">
          <HomeIcon width="1.5em" />
        </NavItem>
      </HStack>
      <List
        margin={0}
        display="flex"
        flexDirection="column"
        gap={2}
        width="full"
      >
        {mapCategories(categories)}
      </List>
    </VStack>
  );
}
