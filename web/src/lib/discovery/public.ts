import "server-only";

import { getAPIAddress } from "@/config";

type AICatalog = {
  specVersion: string;
  entries: Array<{ type: string; url: string }>;
};

export async function getPublicInfo() {
  const response = await fetch(new URL("/api/info", getAPIAddress()), {
    cache: "no-store",
    signal: AbortSignal.timeout(5000),
  });
  if (!response.ok) throw new Error("Public instance information unavailable");

  const info: unknown = await response.json();
  if (
    !info ||
    typeof info !== "object" ||
    !("title" in info) ||
    typeof info.title !== "string" ||
    !("description" in info) ||
    typeof info.description !== "string" ||
    !("web_address" in info) ||
    typeof info.web_address !== "string" ||
    !("api_address" in info) ||
    typeof info.api_address !== "string"
  ) {
    throw new Error("Invalid public instance information");
  }

  return {
    title: info.title,
    description: info.description,
    webAddress: info.web_address,
    apiAddress: info.api_address,
    oauthEnabled:
      "capabilities" in info &&
      Array.isArray(info.capabilities) &&
      info.capabilities.includes("oauth"),
  };
}

export async function getAICatalog(): Promise<AICatalog | null> {
  try {
    const response = await fetch(
      new URL("/.well-known/ai-catalog.json", getAPIAddress()),
      { cache: "no-store", signal: AbortSignal.timeout(5000) },
    );
    if (!response.ok) return null;
    const catalog: unknown = await response.json();
    if (
      !catalog ||
      typeof catalog !== "object" ||
      !("specVersion" in catalog) ||
      catalog.specVersion !== "1.0" ||
      !("entries" in catalog) ||
      !Array.isArray(catalog.entries)
    )
      return null;
    return catalog as AICatalog;
  } catch {
    return null;
  }
}

export function mcpCardURL(
  catalog: AICatalog | null,
  apiAddress: string,
): string | null {
  const expected = new URL("/mcp/server-card", apiAddress).href;
  const entry = catalog?.entries.find(
    (item) =>
      item?.type === "application/mcp-server-card+json" &&
      item.url === expected,
  );
  return entry?.url ?? null;
}
