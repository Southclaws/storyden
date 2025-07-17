"use client";

import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  adminAccessKeyDelete,
  getAdminAccessKeyListKey,
} from "@/api/openapi-client/admin";
import { Identifier, OwnedAccessKeyList } from "@/api/openapi-schema";

export function useAdminAccessKeyList() {
  const { mutate } = useSWRConfig();

  const revokeKey = async (keyId: Identifier) => {
    await handle(async () => {
      // Optimistic update: mark the key as disabled in the cache
      const cacheKey = getAdminAccessKeyListKey();

      await mutate(
        cacheKey,
        async (currentData: { keys: OwnedAccessKeyList } | undefined) => {
          if (!currentData) return currentData;

          // Update the specific key to be disabled
          const updatedKeys = currentData.keys.map((key) =>
            key.id === keyId ? { ...key, enabled: false } : key,
          );

          return { ...currentData, keys: updatedKeys };
        },
        false, // Don't revalidate yet
      );

      // Actually perform the API call
      await adminAccessKeyDelete(keyId);

      // Revalidate the cache to get the latest data
      await mutate(cacheKey);
    });
  };

  return { revokeKey };
}
