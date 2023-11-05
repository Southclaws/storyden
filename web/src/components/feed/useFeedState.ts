import { useFeed } from "./useFeed";
import { useFeedParams } from "./useFeedParams";

export function useFeedState() {
  const params = useFeedParams();
  console.log(params);
  const { data, mutate, handlers } = useFeed(params);

  return {
    ready: true as const,
    data,
    mutate,
    handlers,
  };
}
