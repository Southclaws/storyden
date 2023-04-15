import { Divider, Heading, List, VStack } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { useNavigation } from "../useNavigation";
import { NavItem } from "./components/NavItem";
import { Title } from "./components/Title";
import { Toolbar } from "./components/Toolbar";

const mapCategories = map((c: Category) => (
  <NavItem key={c.id} href={`/c/${c.name}`} w="full">
    <Heading size="sm" role="navigation" variant="ghost" w="full">
      {c.name}
    </Heading>
  </NavItem>
));

export function Sidebar() {
  const { error, categories, title } = useNavigation();

  if (error) return <Unready {...error} />;

  return (
    <VStack as="nav" py={4} gap={2} alignItems="start">
      <Title>{title}</Title>

      <Toolbar />

      <Divider />

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
