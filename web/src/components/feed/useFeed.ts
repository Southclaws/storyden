import { useThreadList } from "src/api/openapi-client/threads";
import { ThreadListParams, ThreadListResult } from "src/api/openapi-schema";

export type Props = {
  params?: ThreadListParams;
  initialData?: ThreadListResult;
};

export function useFeed({ params, initialData }: Props) {
  const { data, mutate, error } = useThreadList(params, {
    swr: { fallbackData: initialData },
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
