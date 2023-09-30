import { Flex, VStack } from "@chakra-ui/react";

import { CategoryList } from "./CategoryList";

export function Menu() {
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
        <CategoryList />
      </Flex>
    </VStack>
  );
}
