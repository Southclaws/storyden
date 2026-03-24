import {
  parseAsArrayOf,
  parseAsBoolean,
  parseAsInteger,
  parseAsString,
  useQueryStates,
} from "nuqs";

import { useAccountList } from "@/api/openapi-client/accounts";
import { AccountListResult } from "@/api/openapi-schema";

export type Props = {
  initialResult: AccountListResult;
  query?: string;
  page?: number;
  initialAdmin?: boolean;
  initialSuspended?: boolean;
  initialAuthServices?: string[];
  initialSort?: string;
};

export function useAdminMemberIndexScreen({
  initialResult,
  query,
  page,
  initialAdmin,
  initialSuspended,
  initialAuthServices,
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
    auth_service: parseAsArrayOf(parseAsString),
  });

  const { data, error } = useAccountList(
    {
      q: filters.q || undefined,
      page: filters.page?.toString(),
      sort: filters.sort || undefined,
      roles: filters.roles?.length ? filters.roles : undefined,
      invited_by: filters.invited_by?.length ? filters.invited_by : undefined,
      joined: filters.joined || undefined,
      admin: filters.admin ?? initialAdmin,
      suspended: filters.suspended ?? initialSuspended,
      auth_service: filters.auth_service?.length
        ? filters.auth_service
        : initialAuthServices,
    },
    { swr: { fallbackData: initialResult } },
  );

  return {
    data,
    error,
    filters: {
      ...filters,
      admin: filters.admin ?? initialAdmin,
      suspended: filters.suspended ?? initialSuspended,
      roles: filters.roles ?? [],
      invited_by: filters.invited_by ?? [],
      auth_service: filters.auth_service ?? initialAuthServices ?? [],
    },
    setFilters,
  };
}
