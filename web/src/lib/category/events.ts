import mitt from "mitt";
import { useEffect } from "react";

import { Identifier } from "@/api/openapi-schema";

export type CategoryEvents = {
  "category:reorder-category": {
    categoryID: Identifier;
    targetCategory: Identifier;
    direction: "above" | "below" | "inside";
    newParent: Identifier | null;
  };
};

export const categoryBus = mitt<CategoryEvents>();

export function useCategoryEvent<K extends keyof CategoryEvents>(
  type: K,
  handler: (event: CategoryEvents[K]) => void,
) {
  useEffect(() => {
    categoryBus.on(type, handler);
    return () => {
      categoryBus.off(type, handler);
    };
  }, [type, handler]);
}

export function useEmitCategoryEvent() {
  return <K extends keyof CategoryEvents>(
    type: K,
    payload: CategoryEvents[K],
  ) => {
    categoryBus.emit(type, payload);
  };
}
