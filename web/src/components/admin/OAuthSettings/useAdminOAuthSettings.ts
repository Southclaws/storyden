"use client";

import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  adminOAuthRefreshTokenDelete,
  getAdminOAuthRefreshTokenListKey,
} from "@/api/openapi-client/admin";
import { Identifier, OAuthRefreshTokenList } from "@/api/openapi-schema";

export function useAdminOAuthSettings() {
  const { mutate } = useSWRConfig();

  const revokeToken = async (tokenID: Identifier) => {
    await handle(async () => {
      const cacheKey = getAdminOAuthRefreshTokenListKey();

      await mutate(
        cacheKey,
        async (currentData: { tokens: OAuthRefreshTokenList } | undefined) => {
          if (!currentData) return currentData;

          const now = new Date().toISOString();
          const tokens = currentData.tokens.map((token) =>
            token.id === tokenID ? { ...token, revoked_at: now } : token,
          );

          return { ...currentData, tokens };
        },
        false,
      );

      await adminOAuthRefreshTokenDelete(tokenID);
      await mutate(cacheKey);
    });
  };

  return { revokeToken };
}
