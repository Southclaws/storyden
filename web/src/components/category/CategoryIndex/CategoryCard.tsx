import Link from "next/link";

import { CardBox } from "@/components/ui/card-box";
import { BulletIcon } from "@/components/ui/icons/Bullet";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { CardRows } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { categoryColourCSS } from "@/lib/category/colours";
import { CategoryTree } from "@/lib/category/tree";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

type CategoryCardProps = {
  category: CategoryTree;
  showChildren: boolean;
};

export function CategoryCard({ category, showChildren }: CategoryCardProps) {
  const coverImage = getAssetURL(category.cover_image?.path);

  const hasSubcategories = category.children.length > 0;

  const plural =
    category.children.length === 1 ? "subcategory" : "subcategories";

  return (
    <CardBox
      position="relative"
      borderWidth="none"
      display="flex"
      height="full"
      justifyContent="space-between"
      gap="0"
      p="0"
      overflow="hidden"
    >
      {coverImage && (
        <img
          src={coverImage}
          alt="" // No alt image, decorative
          aria-hidden="true"
        />
      )}

      <LStack
        flex="1"
        p="1"
        borderWidth="thin"
        borderTopWidth={coverImage ? "none" : undefined}
        borderRadius="lg"
        borderTopRadius={coverImage ? "none" : undefined}
      >
        <LStack h="full" gap="1" justifyContent="space-between">
          <LStack h="full" gap="1">
            <WStack alignItems="start">
              <Link className={linkOverlay()} href={`/d/${category.slug}`}>
                <styled.h2
                  color="text.default"
                  fontWeight="semibold"
                  fontSize="sm"
                  lineHeight="normal"
                >
                  {category.name}
                </styled.h2>
              </Link>

              <CategoryMenu category={category} />
            </WStack>

            <Text variant="supporting">{category.description}</Text>
          </LStack>

          <WStack>
            <HStack gap="1" color="text.muted" fontSize="sm">
              <DiscussionIcon w="4" />
              <Text as="span" variant="metadata">
                {category.postCount}{" "}
                {category.postCount === 1 ? "thread" : "threads"}
              </Text>
              {hasSubcategories && (
                <HStack gap="1" color="text.muted" fontSize="sm">
                  <CategoryIcon w="4" />
                  <Text as="span" variant="metadata">
                    {category.children.length} {plural}
                  </Text>
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
                  background="background.inset"
                  borderColor="border.muted"
                  borderWidth="hairline"
                  borderStyle="solid"
                  borderRadius="sm"
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
                        <styled.h3
                          color="text.default"
                          fontWeight="semibold"
                          fontSize="sm"
                          lineHeight="normal"
                          textWrap="nowrap"
                        >
                          {c.name}
                        </styled.h3>
                      </Link>
                      <BulletIcon />
                      <Text variant="supporting" lineClamp={1}>
                        {c.description}
                      </Text>
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
