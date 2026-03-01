"use client";

import { getSession } from "@/api/openapi-client/misc";

import { DefaultSettings, Settings, parseSettings } from "./settings";
import { useSessionData } from "./session-client";

export async function getSettings(): Promise<Settings> {
  try {
    const data = await getSession();
    return parseSettings(data.info);
  } catch (e) {
    return DefaultSettings;
  }
}

export function useSettings(
  fallbackData?: Settings,
  revalidateOnMount = false,
) {
  const { settings, error } = useSessionData(
    undefined,
    fallbackData,
    revalidateOnMount,
  );

  if (!settings) {
    return {
      ready: false as const,
      error,
    };
  }

  return {
    ready: true as const,
    settings,
  };
}
