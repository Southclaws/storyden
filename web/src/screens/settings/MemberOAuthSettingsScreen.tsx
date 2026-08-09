"use client";

import { parseAsInteger, useQueryStates } from "nuqs";

import {
  useOAuthClientList,
  useOAuthRefreshTokenList,
} from "@/api/openapi-client/auth";
import { OAuthTokenSettings } from "@/components/settings/OAuthTokenSettings/OAuthTokenSettings";
import { Unready } from "@/components/site/Unready";

export function MemberOAuthSettingsScreen() {
  const [filters] = useQueryStates({
    page: parseAsInteger.withDefault(1),
  });

  const tokens = useOAuthRefreshTokenList({ page: filters.page.toString() });
  const clients = useOAuthClientList();

  if (!tokens.data || !clients.data) {
    return <Unready error={tokens.error ?? clients.error} />;
  }

  return (
    <OAuthTokenSettings
      tokens={tokens.data.tokens}
      clients={clients.data.clients}
      tokenPage={{
        currentPage: filters.page,
        totalPages: tokens.data.total_pages,
        pageSize: tokens.data.page_size,
      }}
    />
  );
}
