import { useThreadList } from "src/api/openapi/threads";
import { useQueryParameters } from "./utils";

export function useHomeScreen() {
  const { category } = useQueryParameters();

  const threads = useThreadList({
    categories: category ? [category] : undefined,
  });

  return threads;
}
