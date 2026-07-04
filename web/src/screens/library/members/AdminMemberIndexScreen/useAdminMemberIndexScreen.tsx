import {
  parseAsArrayOf,
  parseAsBoolean,
  parseAsInteger,
  parseAsString,
  useQueryStates,
} from "nuqs";

import { useAccountList } from "@/api/openapi-client/accounts";
import { AccountKind, AccountListResult } from "@/api/openapi-schema";

export type Props = {
  initialResult: AccountListResult;
  query?: string;
  page?: number;
  initialSort?: string;
};

export function useAdminMemberIndexScreen({
  initialResult,
  query,
  page,
  initialSort,
}: Props) {
  const [filters, setFilters] = useQueryStates({
    page: parseAsInteger.withDefault(page ?? 1),
    q: parseAsString.withDefault(query ?? ""),
    sort: parseAsString.withDefault(initialSort ?? "-created_at"),
    roles: parseAsArrayOf(parseAsString),
    invited_by: parseAsArrayOf(parseAsString),
    joined: parseAsString,
    admin: parseAsBoolean,
    suspended: parseAsBoolean,
    kind: parseAsString,
  });

  const kind = parseAccountKind(filters.kind);

  const { data, error } = useAccountList(
    {
      q: filters.q || undefined,
      page: filters.page?.toString(),
      sort: filters.sort || undefined,
      roles: filters.roles?.length ? filters.roles : undefined,
      invited_by: filters.invited_by?.length ? filters.invited_by : undefined,
      joined: filters.joined || undefined,
      admin: filters.admin ?? undefined,
      suspended: filters.suspended ?? undefined,
      kind,
    },
    { swr: { fallbackData: initialResult } },
  );

  return {
    data,
    error,
    filters: {
      ...filters,
      admin: filters.admin ?? undefined,
      suspended: filters.suspended ?? undefined,
      kind,
      roles: filters.roles ?? [],
      invited_by: filters.invited_by ?? [],
    },
    setFilters,
  };
}

function parseAccountKind(value: string | null): AccountKind | undefined {
  if (value === AccountKind.human || value === AccountKind.bot) {
    return value;
  }

  return undefined;
}
