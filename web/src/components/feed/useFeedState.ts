import { useFeed } from "./useFeed";
import { useFeedParams } from "./useFeedParams";

export function useFeedState() {
  const params = useFeedParams();
  const { data, mutate, handlers } = useFeed({
    params: {
      threads: params,
    },
  });

  return {
    ready: true as const,
    data,
    mutate,
    handlers,
  };
}
