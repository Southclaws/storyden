import "server-only";

import { getInfo } from "@/api/openapi-server/misc";

import { DefaultSettings, Settings, parseSettings } from "./settings";

export async function getSettings(): Promise<Settings> {
  try {
    const { data } = await getInfo({
      nocookies: true,
    } as any);
    return parseSettings(data);
  } catch (e) {
    return DefaultSettings;
  }
}
