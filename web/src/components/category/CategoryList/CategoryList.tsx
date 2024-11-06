import { DndContext, closestCenter } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { map } from "lodash/fp";
import { useId } from "react";

import { Category } from "src/api/openapi-schema";
import { Unready } from "src/components/site/Unready";

import { NavigationHeader } from "@/components/site/Navigation/ContentNavigationList/NavigationHeader";
import { styled } from "@/styled-system/jsx";

import { CategoryCreateTrigger } from "../CategoryCreate/CategoryCreateTrigger";

import { CategoryListItem } from "./CategoryListItem";
import { Props, useCategoryList } from "./useCategoryList";

const mapCategories = (isAdmin: boolean) =>
  map((c: Category) => (
    <CategoryListItem key={c.id} {...c} isAdmin={isAdmin} />
  ));

export function CategoryList(props: Props) {
  const { canManageCategories, categories, items, sensors, handleDragEnd } =
    useCategoryList(props);

  const id = useId();

  if (!categories) return <Unready />;

  return (
    <DndContext
      id={id}
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <NavigationHeader
        href="/d"
        controls={canManageCategories && <CategoryCreateTrigger hideLabel />}
      >
        Discussion
      </NavigationHeader>
      <styled.ul
        flexShrink="0"
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
