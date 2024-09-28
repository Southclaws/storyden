import { useFeed } from "./useFeed";
import { useFeedParams } from "./useFeedParams";

export function useFeedState() {
  const params = useFeedParams();
  const { data, mutate } = useFeed({
    params,
  });

  return {
    ready: true as const,
    data,
    mutate,
  };
}
