import { Flex } from "@chakra-ui/react";

import { useCategoryList } from "src/api/openapi/categories";
import { Unready } from "src/components/site/Unready";

import { CategoryListItem } from "../../Sidebar/components/CategoryList/CategoryListItem";

export function CategoryList() {
  const { data } = useCategoryList();
  if (!data) return <Unready />;

  // TODO: Handle errors somewhat well.
  // 1. data is cached, display but with an "offline-mode" warning?
  // 2. data is not cached, first-time render, show an error.
  // swr and pwa stuff probably has some tricks for this.

  return (
    <Flex
      height="full"
      width="full"
      gap={2}
      flexDir="column"
      justifyContent="space-between"
      alignItems="start"
      overflowY="scroll"
    >
      {data.categories.map((c) => (
        <CategoryListItem key={c.id} {...c} isAdmin={false} />
      ))}
    </Flex>
  );
}
