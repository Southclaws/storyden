import "server-only";
import { z } from "zod";

import { getAPIAddress } from "@/config";

const CatalogEntrySchema = z
  .object({
    identifier: z.string(),
    type: z.string(),
    url: z.url().optional(),
    data: z.unknown().optional(),
  })
  .refine((entry) => (entry.url !== undefined) !== (entry.data !== undefined), {
    message: "A catalog entry must contain exactly one of url or data",
  });

const AICatalogSchema = z.object({
  specVersion: z.string().regex(/^1\.\d+$/),
  entries: z.array(CatalogEntrySchema),
});

type AICatalog = z.infer<typeof AICatalogSchema>;

export async function getAICatalog(): Promise<AICatalog | null> {
  try {
    const response = await fetch(
      new URL("/.well-known/ai-catalog.json", getAPIAddress()),
      { cache: "no-store", signal: AbortSignal.timeout(5000) },
    );

    if (!response.ok) {
      return null;
    }

    const catalog = AICatalogSchema.safeParse(await response.json());

    if (!catalog.success) {
      return null;
    }

    return catalog.data;
  } catch {
    return null;
  }
}

export function mcpCardURL(catalog: AICatalog | null): string | null {
  const entry = catalog?.entries.find(
    (item) =>
      item.type === "application/mcp-server-card+json" &&
      item.url &&
      ["http:", "https:"].includes(new URL(item.url).protocol),
  );
  return entry?.url ?? null;
}
