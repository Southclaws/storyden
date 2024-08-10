"use client";

import {
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { arrayMove, sortableKeyboardCoordinates } from "@dnd-kit/sortable";
import { keyBy } from "lodash";
import { useEffect, useState } from "react";

import {
  categoryUpdateOrder,
  useCategoryList as useGetCategoryList,
} from "src/api/openapi-client/categories";
import { Category } from "src/api/openapi-schema";
import { useSession } from "src/auth";

export function useCategoryList() {
  const session = useSession();
  const categoryListResponse = useGetCategoryList();
  const [items, setItems] = useState<string[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  useEffect(() => {
    // Whenever we get new data from the server, update the category list but
    // also update the ordering array to keep the drag state in sync.

    const cats = categoryListResponse.data?.categories ?? [];
    const ids = cats.map((c) => c.id);

    setCategories(cats);
    setItems(ids);
  }, [categoryListResponse.data]);

  async function handleDragEnd(event: any) {
    const { active, over } = event;

    if (active.id === over.id) {
      return;
    }

    const oldIndex = items.indexOf(active.id);
    const newIndex = items.indexOf(over.id);

    // First, update the items ordering.
    const newOrder = arrayMove(items, oldIndex, newIndex);
    setItems(newOrder);

    // Next, re-order the categories according to the new ordering array.
    const newCategories = reorder(newOrder)(categories);
    setCategories(newCategories);

    // Finally, send the new order to the server and set the categories list one
    // last time just in case it changed on the server side since the last call.
    const newCategoriesFromServer = await categoryUpdateOrder(newOrder);
    setCategories(newCategoriesFromServer.categories);
  }

  return {
    isAdmin: session?.admin ?? false,
    // always use the items array as the source of truth for ordering.
    categories: reorder(items)(categories),
    items,
    sensors,
    handleDragEnd,
  };
}

const reorder =
  (ids: string[]) =>
  (x: Category[]): Category[] => {
    const categoryTable = keyBy(x, "id");

    return ids.map((id) => categoryTable[id] as Category);
  };
