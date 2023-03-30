import { useCategoryList } from "src/api/openapi/categories";

export function useSidebar() {
  const { data, error } = useCategoryList();

  return {
    categories: data?.categories,
    error,
  };
}
