import { parseAsArrayOf, parseAsString, useQueryStates } from "nuqs";

import { useProfileList } from "src/api/openapi-client/profiles";
import { PublicProfileListResult } from "src/api/openapi-schema";

export type Props = {
  initialResult: PublicProfileListResult;
  query?: string;
  page?: number;
};

export function useMemberIndexScreen({ initialResult, query, page }: Props) {
  const [filters] = useQueryStates({
    roles: parseAsArrayOf(parseAsString),
    invited_by: parseAsArrayOf(parseAsString),
    joined: parseAsString,
    sort: parseAsString,
  });

  const { data, error } = useProfileList(
    {
      q: query,
      page: page?.toString(),
      roles: filters.roles || undefined,
      invited_by: filters.invited_by || undefined,
      joined: filters.joined || undefined,
      sort: filters.sort || undefined,
    },
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
