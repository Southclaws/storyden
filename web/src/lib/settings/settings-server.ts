import { cache } from "react";
import "server-only";

import { getInfo } from "@/api/openapi-server/misc";

import { DefaultSettings, Settings, parseSettings } from "./settings";

const getSettingsCached = cache(async () => {
  const { data } = await getInfo({
    cache: "no-store",
  });

  return parseSettings(data);
});

export type SettingsResult =
  | { ready: true; settings: Settings }
  | { ready: false; error: unknown };

export async function getSettingsResult(): Promise<SettingsResult> {
  try {
    const settings = await getSettingsCached();

    return { ready: true, settings };
  } catch (error) {
    return { ready: false, error };
  }
}

export async function getSettings(): Promise<Settings> {
  const result = await getSettingsResult();

  return result.ready ? result.settings : DefaultSettings;
}
