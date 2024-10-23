"use client";

import { Arguments, useSWRConfig } from "swr";

import { adminSettingsUpdate } from "@/api/openapi-client/admin";
import { getGetInfoKey } from "@/api/openapi-client/misc";
import { AdminSettingsMutableProps, Info } from "@/api/openapi-schema";

export function useSettingsMutation(initialValue: Info) {
  const { mutate } = useSWRConfig();

  const infoKey = getGetInfoKey()[0];

  function keyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(infoKey);
  }

  async function revalidate(data?: Info) {
    await mutate(keyFilterFn, data);
  }

  async function updateSettings(data: AdminSettingsMutableProps) {
    const newData = { ...initialValue, ...data } satisfies Info;

    await mutate(keyFilterFn, newData, { revalidate: false });

    await adminSettingsUpdate(data);
  }

  return {
    updateSettings,
    revalidate,
  };
}
