import { Divider, VStack } from "@chakra-ui/react";
import { useNavigation } from "../useNavigation";
import { CategoryList } from "./components/CategoryList";
import { Title } from "./components/Title";
import { Toolbar } from "./components/Toolbar";

export function Sidebar() {
  const { title, isAuthenticated } = useNavigation();

  return (
    <VStack as="nav" py={4} gap={2} alignItems="start">
      <Title>{title}</Title>

      <Toolbar isAuthenticated={isAuthenticated} />

      <Divider />

      <CategoryList />
    </VStack>
  );
}
