import { useDatagraphSearch } from "src/api/openapi-client/datagraph";
import {
  DatagraphSearchParams,
  DatagraphSearchResult,
} from "src/api/openapi-schema";

export function useSearch(
  params?: DatagraphSearchParams,
  initialResults?: DatagraphSearchResult,
) {
  const { data, error, mutate } = useDatagraphSearch(params, {
    swr: {
      fallbackData: initialResults,
    },
  });

  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    data,
    mutate,
  };
}
