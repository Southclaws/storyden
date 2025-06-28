"use client";

import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  accessKeyDelete,
  getAccessKeyListKey,
} from "@/api/openapi-client/auth";
import { AccessKeyList, Identifier } from "@/api/openapi-schema";

export function useAccessKeySettings() {
  const { mutate } = useSWRConfig();

  const revokeKey = async (keyId: Identifier) => {
    await handle(async () => {
      const cacheKey = getAccessKeyListKey();

      await mutate(
        cacheKey,
        async (currentData: { keys: AccessKeyList } | undefined) => {
          if (!currentData) return currentData;

          const updatedKeys = currentData.keys.map((key) =>
            key.id === keyId ? { ...key, enabled: false } : key,
          );

          return { ...currentData, keys: updatedKeys };
        },
        false,
      );

      await accessKeyDelete(keyId);

      await mutate(cacheKey);
    });
  };

  return { revokeKey };
}
