"use client";

import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getOAuthClientListKey,
  oAuthClientDelete,
} from "@/api/openapi-client/auth";
import { Identifier, OAuthClientList } from "@/api/openapi-schema";

export function useOAuthClientSettings() {
  const { mutate } = useSWRConfig();

  const deleteClient = async (clientID: Identifier) => {
    await handle(async () => {
      const cacheKey = getOAuthClientListKey();

      await mutate(
        cacheKey,
        async (currentData: { clients: OAuthClientList } | undefined) => {
          if (!currentData) return currentData;

          return {
            ...currentData,
            clients: currentData.clients.filter(
              (client) => client.id !== clientID,
            ),
          };
        },
        false,
      );

      await oAuthClientDelete(clientID);
      await mutate(cacheKey);
    });
  };

  return { deleteClient };
}
