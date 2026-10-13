import "server-only";

import { getAPIAddress } from "@/config";
import * as z from "@/lib/validation";

const CatalogEntrySchema = z
  .object({
    identifier: z.string(),
    type: z.string(),
    url: z.optional(z.url()),
    data: z.optional(z.unknown()),
  })
  .check(
    z.refine(
      (entry) => (entry.url !== undefined) !== (entry.data !== undefined),
      {
        message: "A catalog entry must contain exactly one of url or data",
      },
    ),
  );

const AICatalogSchema = z.object({
  specVersion: z.string().check(z.regex(/^1\.\d+$/)),
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
