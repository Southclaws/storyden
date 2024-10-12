import Link from "next/link";

import { Category } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { CardGrid } from "@/components/ui/rich-card";
import { categoryColourCSS } from "@/lib/category/colours";
import { CardBox, HStack, LStack, styled } from "@/styled-system/jsx";

import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";
import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

export type Props = {
  categories: Category[];
};

export function CategoryCardList({ categories }: Props) {
  return (
    <LStack>
      <HStack w="full" justify="space-between">
        <Heading>Discussion categories</Heading>

        <CategoryCreateTrigger />
      </HStack>

      <CardGrid>
        {categories.map((c) => (
          <CategoryCard key={c.id} {...c} />
        ))}
      </CardGrid>
    </LStack>
  );
}

export function CategoryCard(props: Category) {
  const cssProps = categoryColourCSS(props.colour);

  return (
    <CardBox
      style={cssProps}
      borderColor="colorPalette.muted"
      borderLeftWidth="thick"
      borderLeftStyle="solid"
      display="flex"
      justifyContent="space-between"
    >
      <LStack>
        <HStack w="full" justify="space-between" alignItems="start">
          <Link href={`/d/${props.slug}`}>
            <Heading>{props.name}</Heading>
          </Link>

          <CategoryMenu category={props} />
        </HStack>

        <styled.p color="fg.muted">{props.description}</styled.p>
      </LStack>

      <styled.p color="fg.subtle">{props.postCount} threads</styled.p>
    </CardBox>
  );
}
