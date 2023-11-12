import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { CategoryList } from "src/components/category/CategoryList/CategoryList";
import { Flex, VStack } from "src/theme/components";

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
        <CategoryCreateTrigger />
      </Flex>
    </VStack>
  );
}
