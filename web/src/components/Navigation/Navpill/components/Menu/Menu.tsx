import { Flex, VStack } from "@chakra-ui/react";
import { CategoryList } from "../CategoryList";
import { useMenu } from "./useMenu";

export function Menu() {
  const { category } = useMenu();

  return (
    <VStack width="full" p={2}>
      <Flex
        maxHeight="80vh"
        flexDir="column"
        justifyContent="center"
        alignItems="center"
        maxW="container.sm"
        width="full"
        pointerEvents="auto"
        gap={2}
      >
        <CategoryList category={category} />
      </Flex>
    </VStack>
  );
}
