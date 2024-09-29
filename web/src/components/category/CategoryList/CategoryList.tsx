import { DndContext, closestCenter } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { map } from "lodash/fp";

import { Category } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";

import { AddAction } from "@/components/site/Action/Add";
import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { styled } from "@/styled-system/jsx";

import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";

import { CategoryListItem } from "./CategoryListItem";
import { useCategoryList } from "./useCategoryList";

const mapCategories = (isAdmin: boolean) =>
  map((c: Category) => (
    <CategoryListItem key={c.id} {...c} isAdmin={isAdmin} />
  ));

export function CategoryList() {
  const { canManageCategories, categories, items, sensors, handleDragEnd } =
    useCategoryList();

  if (!categories) return <Unready />;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <NavigationHeader
        controls={
          canManageCategories && (
            <CategoryCreateTrigger>
              <AddAction size="xs" color="fg.subtle" title="Add a category" />
            </CategoryCreateTrigger>
          )
        }
      >
        Discussion
      </NavigationHeader>
      <styled.ul
        overflow="hidden"
        margin="0"
        display="flex"
        flexDirection="column"
        gap="1"
        width="full"
        css={{
          touchAction: "none",
        }}
      >
        <SortableContext items={items} strategy={verticalListSortingStrategy}>
          {mapCategories(canManageCategories)(categories)}
        </SortableContext>
      </styled.ul>
    </DndContext>
  );
}
