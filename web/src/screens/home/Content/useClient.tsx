import { useThreadList } from "src/api/openapi/threads";

import { useQueryParameters } from "../utils";

export function useClient() {
  const { category } = useQueryParameters();

  const threads = useThreadList({
    categories: category ? [category] : undefined,
  });

  return threads;
}
