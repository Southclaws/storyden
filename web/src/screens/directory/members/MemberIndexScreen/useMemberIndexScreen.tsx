import { useProfileList } from "src/api/openapi/profiles";
import { PublicProfileListResult } from "src/api/openapi/schemas";

export type Props = {
  query?: string;
  page?: number;
  profiles: PublicProfileListResult;
};

export function useMemberIndexScreen(props: Props) {
  const { data, mutate, error } = useProfileList(
    {
      q: props.query,
      page: props.page?.toString(),
    },
    {
      swr: {
        fallbackData: props.profiles,
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
