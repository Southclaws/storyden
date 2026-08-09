"use client";

import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getOAuthRefreshTokenListKey,
  oAuthRefreshTokenDelete,
} from "@/api/openapi-client/auth";
import {
  Identifier,
  OAuthRefreshTokenListParams,
  OAuthRefreshTokenListResult,
} from "@/api/openapi-schema";

export function useOAuthTokenSettings(params: OAuthRefreshTokenListParams) {
  const { mutate } = useSWRConfig();

  const revokeToken = async (tokenID: Identifier) => {
    await handle(async () => {
      const cacheKey = getOAuthRefreshTokenListKey(params);

      await mutate(
        cacheKey,
        async (currentData: OAuthRefreshTokenListResult | undefined) => {
          if (!currentData) return currentData;

          const now = new Date().toISOString();
          const tokens = currentData.tokens.map((token) =>
            token.id === tokenID ? { ...token, revoked_at: now } : token,
          );

          return { ...currentData, tokens };
        },
        false,
      );

      await oAuthRefreshTokenDelete(tokenID);
      await mutate(cacheKey);
    });
  };

  return { revokeToken };
}
