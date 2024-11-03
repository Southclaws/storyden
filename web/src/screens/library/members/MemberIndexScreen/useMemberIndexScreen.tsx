import { useProfileList } from "src/api/openapi-client/profiles";
import { PublicProfileListResult } from "src/api/openapi-schema";

export type Props = {
  initialResult: PublicProfileListResult;
  query?: string;
  page?: number;
};

export function useMemberIndexScreen({ initialResult, query, page }: Props) {
  const { data, error } = useProfileList(
    { q: query, page: page?.toString() },
    { swr: { fallbackData: initialResult } },
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
  };
}
