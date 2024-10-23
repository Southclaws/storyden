"use client";

import { getInfo, useGetInfo } from "@/api/openapi-client/misc";

import { DefaultSettings, Settings, parseSettings } from "./settings";

export async function getSettings(): Promise<Settings> {
  try {
    const data = await getInfo();
    return parseSettings(data);
  } catch (e) {
    return DefaultSettings;
  }
}

export function useSettings(fallbackData?: Settings) {
  const { data, error } = useGetInfo({ swr: { fallbackData } });
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const settings = parseSettings(data);
  return {
    ready: true as const,
    settings,
  };
}
