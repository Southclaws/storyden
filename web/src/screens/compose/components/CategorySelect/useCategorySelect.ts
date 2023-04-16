import { useCategoryList } from "src/api/openapi/categories";

export function useCategorySelect() {
  const { data, error } = useCategoryList();
  return { categories: data?.categories, error };
}
