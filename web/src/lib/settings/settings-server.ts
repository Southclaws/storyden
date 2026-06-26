import { cacheLife } from "next/cache";
import "server-only";

import { getInfo } from "@/api/openapi-server/misc";

import { DefaultSettings, Settings, parseSettings } from "./settings";

async function getSettingsCached(): Promise<Settings> {
  "use cache";
  cacheLife("hours");
  const { data } = await getInfo({
    cache: "no-store",
  });
  return parseSettings(data);
}

export async function getSettings(): Promise<Settings> {
  try {
    return await getSettingsCached();
  } catch (e) {
    return DefaultSettings;
  }
}
