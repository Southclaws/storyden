"use client";

import { useSWRConfig } from "swr";

import { mutateTransaction } from "@/api/mutate";
import {
  adminSettingsUpdate,
  getAdminSettingsGetKey,
} from "@/api/openapi-client/admin";
import { getGetInfoKey } from "@/api/openapi-client/misc";
import {
  AdminSettingsMutableProps,
  AdminSettingsProps,
  Info,
  MessageOfTheDay,
  MessageOfTheDayMutableProps,
} from "@/api/openapi-schema";

import { AdminSettings } from "./settings";

export function useSettingsMutation() {
  const { mutate } = useSWRConfig();

  const infoKey = getGetInfoKey();
  const adminSettingsKey = getAdminSettingsGetKey();

  async function revalidate(data?: Info) {
    await mutate(infoKey, data);
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
            const { services, motd, ...publicPatch } = patch;

            return {
              ...current,
              ...publicPatch,
              motd: mergeMotd(current.motd, motd),
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

function mergeMotd(
  current: MessageOfTheDay | undefined,
  patch: MessageOfTheDayMutableProps | undefined,
): MessageOfTheDay | undefined {
  if (patch === undefined) {
    return current;
  }

  // Sending an empty object is used by admin settings as an explicit clear.
  if (Object.keys(patch).length === 0) {
    return undefined;
  }

  const content = patch.content ?? current?.content;
  if (content === undefined) {
    return current;
  }

  return {
    ...current,
    ...patch,
    content,
    metadata:
      patch.metadata || current?.metadata
        ? {
            ...(current?.metadata ?? {}),
            ...(patch.metadata ?? {}),
          }
        : undefined,
  };
}
