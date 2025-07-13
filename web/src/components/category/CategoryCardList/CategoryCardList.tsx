import Link from "next/link";

import { Category } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { categoryColourCSS } from "@/lib/category/colours";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";

import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";
import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

export type Props = {
  categories: Category[];
};

export function CategoryCardGrid({ categories }: Props) {
  return (
    <LStack>
      <WStack>
        <Heading>Discussion categories</Heading>

        <CategoryCreateTrigger />
      </WStack>

      <CardGrid>
        {categories.map((c) => (
          <CategoryCard key={c.id} {...c} />
        ))}
      </CardGrid>
    </LStack>
  );
}

export function CategoryCardList({ categories }: Props) {
  return (
    <LStack>
      <WStack>
        <Heading>Discussion categories</Heading>

        <CategoryCreateTrigger />
      </WStack>

      <CardRows>
        {categories.map((c) => (
          <CategoryCard key={c.id} {...c} />
        ))}
      </CardRows>
    </LStack>
  );
}

export function CategoryCard(props: Category) {
  const cssProps = categoryColourCSS(props.colour);

  return (
    <CardBox
      position="relative"
      style={cssProps}
      borderColor="colorPalette.muted"
      borderLeftWidth="thick"
      borderLeftStyle="solid"
      display="flex"
      justifyContent="space-between"
    >
      <WStack alignItems="start">
        <Link className={linkOverlay()} href={`/d/${props.slug}`}>
          <Heading>{props.name}</Heading>
        </Link>

        <CategoryMenu category={props} />
      </WStack>

      <styled.p color="fg.muted">{props.description}</styled.p>

      <HStack gap="1" color="fg.subtle">
        <DiscussionIcon w="4" />
        <styled.p>{props.postCount} threads</styled.p>
      </HStack>
    </CardBox>
  );
}
