import { Flex, VStack } from "@chakra-ui/react";

import { useNavigation } from "src/components/site/Navigation/useNavigation";

import { CategoryList } from "./CategoryList";

export function Menu() {
  const { category } = useNavigation();

  return (
    <VStack width="full" p={2}>
      <Flex
        maxHeight="80vh"
        flexDir="column"
        justifyContent="center"
        alignItems="center"
        maxW="768px"
        width="full"
        pointerEvents="auto"
        gap={2}
      >
        <CategoryList category={category} />
      </Flex>
    </VStack>
  );
}
