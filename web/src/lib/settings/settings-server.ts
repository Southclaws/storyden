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

export async function getSettings(): Promise<Settings> {
  try {
    const settings = await getSettingsCached();
    return settings;
  } catch (e) {
    return DefaultSettings;
  }
}
