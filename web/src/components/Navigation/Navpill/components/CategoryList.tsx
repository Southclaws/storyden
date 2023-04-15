import { Button, Flex } from "@chakra-ui/react";
import { map } from "lodash/fp";
import { useCategoryList } from "src/api/openapi/categories";
import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { Anchor } from "src/components/site/Anchor";

const mapCategories = (selected?: string) =>
  map((c: Category) => (
    <Anchor key={c.id} href={`/c/${c.name}`} w="full">
      <Button bgColor={c.name === selected ? "blackAlpha.200" : ""} w="full">
        {c.name}
      </Button>
    </Anchor>
  ));

type Props = { category?: string };

export function CategoryList({ category }: Props) {
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
      {mapCategories(category)(data.categories)}
    </Flex>
  );
}
