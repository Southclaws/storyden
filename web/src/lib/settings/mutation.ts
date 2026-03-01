"use client";

import { useSWRConfig } from "swr";

import { mutateTransaction } from "@/api/mutate";
import {
  adminSettingsUpdate,
  getAdminSettingsGetKey,
} from "@/api/openapi-client/admin";
import {
  getGetInfoKey,
  getGetSessionKey,
} from "@/api/openapi-client/misc";
import {
  AdminSettingsMutableProps,
  AdminSettingsProps,
  GetSessionOKResponse,
  Info,
} from "@/api/openapi-schema";

import { AdminSettings } from "./settings";

export function useSettingsMutation() {
  const { mutate } = useSWRConfig();

  const infoKey = getGetInfoKey();
  const adminSettingsKey = getAdminSettingsGetKey();
  const sessionKey = getGetSessionKey();

  async function revalidate(data?: Info) {
    if (!data) {
      await Promise.all([mutate(infoKey), mutate(sessionKey)]);
      return;
    }

    await Promise.all([
      mutate(infoKey, data, { revalidate: false }),
      mutate(
        sessionKey,
        (current: GetSessionOKResponse | undefined) => {
          if (!current) return current;

          return {
            ...current,
            info: data,
          } satisfies GetSessionOKResponse;
        },
        { revalidate: false },
      ),
    ]);
  }

  async function updateSettings(patch: AdminSettingsMutableProps) {
    await mutateTransaction(
      mutate,
      [
        {
          key: adminSettingsKey,
          optimistic: (current: AdminSettings | undefined) =>
            current
              ? {
                  ...current,
                  ...patch,
                  metadata: {
                    ...current.metadata,
                    ...patch.metadata,
                  },
                }
              : current,
          commit: (_current, result) => {
            return result;
          },
        },
        {
          key: infoKey,
          optimistic: (current: Info | undefined) => {
            if (!current) return current;
            const { services, ...publicPatch } = patch;

            return {
              ...current,
              ...publicPatch,
              metadata: {
                ...current.metadata,
                ...patch.metadata,
              },
            } satisfies Info;
          },
          commit: (current, result) => {
            const updated = result;
            if (!current) return current;

            return adminToInfo(updated);
          },
        },
        {
          key: sessionKey,
          optimistic: (current: GetSessionOKResponse | undefined) => {
            if (!current) return current;
            const { services, ...publicPatch } = patch;

            return {
              ...current,
              info: {
                ...current.info,
                ...publicPatch,
                metadata: {
                  ...current.info.metadata,
                  ...patch.metadata,
                },
              },
            } satisfies GetSessionOKResponse;
          },
          commit: (current, result) => {
            if (!current) return current;

            return {
              ...current,
              info: adminToInfo(result),
            } satisfies GetSessionOKResponse;
          },
        },
      ],
      async () => {
        return await adminSettingsUpdate(patch);
      },
    );
  }

  return {
    updateSettings,
    revalidate,
  };
}

function adminToInfo(admin: AdminSettingsProps): Info {
  return {
    ...admin,
    capabilities: admin.capabilities ?? [],
    onboarding_status: "complete",
  };
}
