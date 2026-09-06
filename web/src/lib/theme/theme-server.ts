import { unstable_cache } from "next/cache";
import "server-only";

import { getAPIAddress } from "@/config";

import { EMPTY_THEME_MANIFEST, parseThemeManifest } from "./manifest";

export const THEME_CACHE_TAG = "theme";

type ThemeBundle = {
  stylesheet: string;
  script: string;
};

const EMPTY_THEME_BUNDLE: ThemeBundle = {
  stylesheet: "",
  script: "",
};

const getCachedThemeBundle = unstable_cache(
  async () => fetchThemeBundle(),
  ["theme-bundle"],
  { revalidate: 300, tags: [THEME_CACHE_TAG] },
);

export async function getServerThemeBundle() {
  return getCachedThemeBundle();
}

async function fetchThemeBundle(): Promise<ThemeBundle> {
  try {
    const apiAddress = getAPIAddress();
    const response = await fetch(`${apiAddress}/api/info/theme`, {
      cache: "no-store",
      signal: AbortSignal.timeout(3_000),
    });
    if (!response.ok) {
      return EMPTY_THEME_BUNDLE;
    }

    const manifest = parseThemeManifest(await response.json(), apiAddress);
    if (manifest === EMPTY_THEME_MANIFEST) return EMPTY_THEME_BUNDLE;

    const [stylesheets, scripts] = await Promise.all([
      Promise.all(manifest.stylesheets.map(fetchThemeAsset)),
      Promise.all(manifest.scripts.map(fetchThemeAsset)),
    ]);

    return {
      stylesheet: stylesheets.join("\n"),
      script: scripts.join("\n;\n"),
    };
  } catch {
    return EMPTY_THEME_BUNDLE;
  }
}

async function fetchThemeAsset(asset: { href: string }) {
  const response = await fetch(asset.href, {
    cache: "force-cache",
    signal: AbortSignal.timeout(3_000),
  });
  if (!response.ok) {
    throw new Error(`Theme asset returned ${response.status}.`);
  }
  return response.text();
}
