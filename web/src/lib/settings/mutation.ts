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
            if (!current) return current;

            return mergeAdminSettingsIntoInfo(current, result);
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

export function mergeAdminSettingsIntoInfo(
  current: Info,
  admin: AdminSettingsProps,
): Info {
  return {
    ...current,
    accent_colour: admin.accent_colour,
    api_address: admin.api_address,
    authentication_mode: admin.authentication_mode,
    capabilities: admin.capabilities ?? current.capabilities,
    content: admin.content,
    description: admin.description,
    metadata: admin.metadata ?? current.metadata,
    motd: admin.motd,
    registration_mode: admin.registration_mode,
    title: admin.title,
    web_address: admin.web_address,
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
