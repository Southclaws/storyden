import { DndContext, closestCenter } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { map } from "lodash/fp";

import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/site/Unready";

import { styled } from "@/styled-system/jsx";

import { CategoryListItem } from "./CategoryListItem";
import { useCategoryList } from "./useCategoryList";

const mapCategories = (isAdmin: boolean) =>
  map((c: Category) => (
    <CategoryListItem key={c.id} {...c} isAdmin={isAdmin} />
  ));

export function CategoryList() {
  const { isAdmin, categories, items, sensors, handleDragEnd } =
    useCategoryList();

  if (!categories) return <Unready />;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <styled.ul
        overflow="hidden"
        margin="0"
        display="flex"
        flexDirection="column"
        gap="2"
        width="full"
        css={{
          touchAction: "none",
        }}
      >
        <SortableContext items={items} strategy={verticalListSortingStrategy}>
          {mapCategories(isAdmin)(categories)}
        </SortableContext>
      </styled.ul>
    </DndContext>
  );
}
