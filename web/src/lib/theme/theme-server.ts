import "server-only";

import { getAPIAddress } from "@/config";
import { EMPTY_THEME_MANIFEST, parseThemeManifest } from "@/lib/theme/manifest";

export async function getServerThemeManifest() {
  try {
    const response = await fetch(`${getAPIAddress()}/api/info/theme`, {
      cache: "no-store",
    });
    if (!response.ok) {
      return EMPTY_THEME_MANIFEST;
    }

    const data = await response.json();
    return parseThemeManifest(data);
  } catch {
    return EMPTY_THEME_MANIFEST;
  }
}
