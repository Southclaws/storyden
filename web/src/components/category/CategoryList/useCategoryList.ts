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
import { Category, CategoryListOKResponse } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { hasPermission } from "@/utils/permissions";

export type Props = {
  initialCategoryList?: CategoryListOKResponse;
};

export function useCategoryList({ initialCategoryList }: Props) {
  const session = useSession();
  const categoryListResponse = useGetCategoryList({
    swr: { fallbackData: initialCategoryList },
  });
  const [items, setItems] = useState<string[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 4,
      },
    }),
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

  const canManageCategories = hasPermission(session, "MANAGE_CATEGORIES");

  // Categories are empty on server render, but we still want to show something.
  const orderedCategories =
    categories.length > 0
      ? reorder(items)(categories)
      : (categoryListResponse.data?.categories ?? []);

  return {
    canManageCategories,
    // always use the items array as the source of truth for ordering.
    categories: orderedCategories,
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
