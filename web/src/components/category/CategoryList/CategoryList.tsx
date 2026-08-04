"use client";

import { KeyedMutator } from "swr";

import { handle } from "@/api/client";
import {
  categoryUpdatePosition,
  useCategoryList as useGetCategoryList,
} from "@/api/openapi-client/categories";
import {
  Category,
  CategoryListOKResponse,
  CategoryUpdatePositionBody,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { Unready } from "@/components/site/Unready";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { DragTree } from "@/components/ui/tree-view";
import { useCategoryEvent } from "@/lib/category/events";
import { CategoryTree, buildCategoryTree } from "@/lib/category/tree";
import {
  DragItemCategory,
  DragItemCategoryDivider,
} from "@/lib/dragdrop/provider";
import { HStack, LStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { DiscussionRoute } from "../../site/Navigation/Anchors/Discussion";
import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";
import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

export type Props = {
  initialCategoryList?: CategoryListOKResponse;
  currentPath?: string;
};

const getCategoryId = (category: CategoryTree) => category.id;
const getCategoryLabel = (category: CategoryTree) => category.name;
const getCategoryHref = (category: CategoryTree) => `/d/${category.slug}`;
const getCategoryChildren = (category: CategoryTree) => category.children;

export function CategoryList({ initialCategoryList, currentPath }: Props) {
  const { data, error, mutate } = useGetCategoryList({
    swr: { fallbackData: initialCategoryList },
  });

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <CategoryListTree
      categories={data.categories}
      currentPath={currentPath}
      mutate={mutate}
    />
  );
}

export function CategoryListTree({
  categories,
  currentPath,
  mutate,
}: {
  categories: Category[];
  currentPath?: string;
  mutate: KeyedMutator<CategoryListOKResponse>;
}) {
  const session = useSession();
  const canManageCategories = hasPermission(session, "MANAGE_CATEGORIES");
  const tree = buildCategoryTree(categories);

  useCategoryEvent(
    "category:reorder-category",
    async ({ categorySlug, direction, newParent, targetCategory }) => {
      await handle(async () => {
        const params: CategoryUpdatePositionBody = (() => {
          switch (direction) {
            case "above":
              return { before: targetCategory, parent: newParent };
            case "below":
              return { after: targetCategory, parent: newParent };
            case "inside":
              return { parent: targetCategory };
          }
        })();

        const response = await categoryUpdatePosition(categorySlug, params);
        await mutate(response, { revalidate: false });
      });
    },
  );

  return (
    <LStack gap="0">
      <NavigationHeader
        href={DiscussionRoute}
        controls={canManageCategories && <CategoryCreateTrigger hideLabel />}
      >
        <HStack gap="1">
          <DiscussionIcon />
          Discussion
        </HStack>
      </NavigationHeader>

      <DragTree
        id="category-navigation"
        label="Discussion categories"
        nodes={tree}
        getId={getCategoryId}
        getLabel={getCategoryLabel}
        getHref={getCategoryHref}
        getChildren={getCategoryChildren}
        defaultExpandedIds={tree.map(getCategoryId)}
        isSelected={(category) => getCategoryHref(category) === currentPath}
        drag={{
          enabled: canManageCategories,
          getNodeData: ({ node }) =>
            ({
              type: "category",
              categoryID: node.id,
              category: node,
              hasChildren: node.children.length > 0,
            }) satisfies DragItemCategory,
          getDividerData: ({ node, parent, direction }) =>
            ({
              type: "category-divider",
              direction,
              siblingCategoryID: node.id,
              parentID: parent?.id ?? null,
            }) satisfies DragItemCategoryDivider,
        }}
        renderActions={({ node }) =>
          canManageCategories ? (
            <CategoryMenu category={node} triggerProps={{ variant: "plain" }} />
          ) : null
        }
      />
    </LStack>
  );
}
