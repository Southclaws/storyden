"use client";

import { Arguments, useSWRConfig } from "swr";

import {
  adminSettingsUpdate,
  getAdminSettingsGetKey,
} from "@/api/openapi-client/admin";
import { getGetInfoKey } from "@/api/openapi-client/misc";
import { AdminSettingsMutableProps, Info } from "@/api/openapi-schema";

import { AdminSettings } from "./settings";

export function useSettingsMutation() {
  const { mutate } = useSWRConfig();

  const infoKey = getGetInfoKey()[0];
  const adminSettingsKey = getAdminSettingsGetKey();

  function keyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(infoKey);
  }

  async function revalidate(data?: Info) {
    await mutate(keyFilterFn, data);
  }

  async function updateSettings(data: AdminSettingsMutableProps) {
    await mutate(
      adminSettingsKey,
      (current: AdminSettings | undefined) => {
        if (!current) return undefined;

        return {
          ...current,
          ...data,
          metadata: {
            ...current.metadata,
            ...data.metadata,
          },
        } satisfies AdminSettings;
      },
      { revalidate: false },
    );

    await mutate(
      keyFilterFn,
      (currentInfo: Info | undefined) => {
        if (!currentInfo) return undefined;

        // NOTE: This is kinda hacky, we need to rework this in future. It's
        // because the very old /info endpoint returns a subset of settings that
        // are public, then the admin settings includes additional fields that
        // are not public. This doesn't expose any data (as it's just a mutation
        // for the current session) but it's still not ideal.
        const { services, ...infoData } = data;

        return {
          ...currentInfo,
          ...infoData,
          metadata: {
            ...currentInfo.metadata,
            ...data.metadata,
          },
        } satisfies Info;
      },
      { revalidate: false },
    );

    await adminSettingsUpdate(data);
  }

  return {
    updateSettings,
    revalidate,
  };
}
