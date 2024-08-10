import { useThreadList } from "src/api/openapi-client/threads";
import { ThreadListResult } from "src/api/openapi-schema";

export type Props = {
  query?: string;
  page?: number;
  threads: ThreadListResult;
};

export function useThreadIndexScreen(props: Props) {
  const { data, mutate, error } = useThreadList(
    {
      q: props.query,
      page: props.page?.toString(),
    },
    {
      swr: {
        fallbackData: props.threads,
      },
    },
  );

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
