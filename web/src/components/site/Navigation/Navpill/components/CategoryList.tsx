import { Button, Flex } from "@chakra-ui/react";
import { usePathname } from "next/navigation";

import { useCategoryList } from "src/api/openapi/categories";
import { Category } from "src/api/openapi/schemas";
import { Anchor } from "src/components/site/Anchor";
import { Unready } from "src/components/site/Unready";

function CategoryListItem(props: Category) {
  const pathname = usePathname();

  const href = `/c/${props.slug}`;
  const selected = href === pathname;

  return (
    <Anchor href={href} w="full">
      <Button bgColor={selected ? "blackAlpha.200" : ""} w="full">
        {props.name}
      </Button>
    </Anchor>
  );
}

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
        <CategoryListItem key={c.id} {...c} />
      ))}
    </Flex>
  );
}
