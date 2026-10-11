import "server-only";

import { getInfo } from "@/api/openapi-server/misc";

import { resolveDiscoverySettings } from "./discovery";

export async function getDiscoverySettings() {
  const { data } = await getInfo({
    cache: "no-store",
  });

  return resolveDiscoverySettings(data.discovery);
}
