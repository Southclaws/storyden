import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { CategoryList } from "src/components/category/CategoryList/CategoryList";

import { Flex, VStack } from "@/styled-system/jsx";

export function Menu() {
  return (
    <VStack width="full" p="2">
      <Flex
        maxHeight="screen"
        flexDir="column"
        justifyContent="center"
        alignItems="center"
        maxW="prose"
        width="full"
        pointerEvents="auto"
        gap="2"
      >
        <CategoryList />
        <CategoryCreateTrigger />
      </Flex>
    </VStack>
  );
}
