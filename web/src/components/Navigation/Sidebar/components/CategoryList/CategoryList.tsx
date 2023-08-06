import { List } from "@chakra-ui/react";
import { DndContext, closestCenter } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { map } from "lodash/fp";

import { Category } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";

import { CategoryListItem } from "./CategoryListItem";
import { useCategoryList } from "./useCategoryList";

const mapCategories = map((c: Category) => (
  <CategoryListItem key={c.id} {...c} />
));

export function CategoryList() {
  const { categories, items, sensors, handleDragEnd } = useCategoryList();

  if (!categories) return <Unready />;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <List
        overflow="hidden"
        margin={0}
        display="flex"
        flexDirection="column"
        gap={2}
        width="full"
      >
        <SortableContext items={items} strategy={verticalListSortingStrategy}>
          {mapCategories(categories)}
        </SortableContext>
      </List>
    </DndContext>
  );
}
