import Link from "next/link";

import { Heading } from "@/components/ui/heading";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { CardRows } from "@/components/ui/rich-card";
import { categoryColourCSS } from "@/lib/category/colours";
import { CategoryTree } from "@/lib/category/tree";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

type CategoryCardProps = {
  category: CategoryTree;
  showChildren: boolean;
};

export function CategoryCard({ category, showChildren }: CategoryCardProps) {
  const cssProps = categoryColourCSS(category.colour);

  const hasSubcategories = category.children.length > 0;

  const plural =
    category.children.length === 1 ? "subcategory" : "subcategories";

  return (
    <CardBox
      position="relative"
      style={cssProps}
      borderColor="colorPalette.border"
      borderLeftWidth="thick"
      borderLeftStyle="solid"
      display="flex"
      justifyContent="space-between"
      gap="0"
      p="0"
      overflow="hidden"
    >
      <img
        src={getAssetURL(category.cover_image?.path)}
        alt="" // No alt image, decorative
        aria-hidden="true"
      />

      <LStack p="1">
        <LStack h="full" gap="1" justifyContent="space-between">
          <LStack h="full" gap="1">
            <WStack alignItems="start">
              <Link className={linkOverlay()} href={`/d/${category.slug}`}>
                <Heading>{category.name}</Heading>
              </Link>

              <CategoryMenu category={category} />
            </WStack>

            <styled.p color="fg.muted" fontSize="sm">
              {category.description}
            </styled.p>
          </LStack>

          <WStack>
            <HStack gap="1" color="fg.subtle" fontSize="sm">
              <DiscussionIcon w="4" />
              <styled.p>{category.postCount} {category.postCount === 1 ? "thread" : "threads"}</styled.p>
              {hasSubcategories && (
                <HStack gap="1" color="fg.subtle" fontSize="sm">
                  <CategoryIcon w="4" />
                  <styled.p>
                    {category.children.length} {plural}
                  </styled.p>
                </HStack>
              )}
            </HStack>
          </WStack>
        </LStack>

        {hasSubcategories && showChildren && (
          <CardRows>
            {category.children.map((c) => {
              const cssProps = categoryColourCSS(c.colour);

              return (
                <CardBox
                  key={c.id}
                  position="relative"
                  style={cssProps}
                  borderColor="bg.subtle"
                  borderWidth="hairline"
                  borderStyle="solid"
                  borderLeftColor="colorPalette.border"
                  borderLeftWidth="thick"
                  borderLeftStyle="solid"
                  display="flex"
                  justifyContent="space-between"
                  gap="4"
                  boxShadow="[none]"
                  px="2"
                  py="1"
                >
                  <WStack alignItems="start">
                    <HStack gap="1">
                      <Link className={linkOverlay()} href={`/d/${c.slug}`}>
                        <Heading textWrap="nowrap" fontSize="sm">
                          {c.name}
                        </Heading>
                      </Link>
                      <BulletIcon />
                      <styled.p lineClamp={1} color="fg.muted" fontSize="sm">
                        {c.description}
                      </styled.p>
                    </HStack>
                    <CategoryMenu category={c} />
                  </WStack>
                </CardBox>
              );
            })}
          </CardRows>
        )}
      </LStack>
    </CardBox>
  );
}
