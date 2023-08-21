import { useThreadList } from "src/api/openapi/threads";

import { useQueryParameters } from "./utils";

export function useContent() {
  const { category } = useQueryParameters();

  const threads = useThreadList({
    categories: category ? [category] : undefined,
  });

  return threads;
}
