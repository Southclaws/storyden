import { WEB_ADDRESS, getAPIAddress } from "@/config";

const API_PAGE_SIZE_THREADS = 50;
const API_PAGE_SIZE_NODES = 500;

export type IndexKind = "threads" | "library";
export type IndexEntry = { slug: string; updatedAt: string };
export type IndexPage = { entries: IndexEntry[]; totalPages: number };

export function canonicalURL(path: string): string {
  return new URL(path, WEB_ADDRESS).toString();
}

export function escapeXML(value: string): string {
  return value.replace(/[&<>"']/g, (char) => {
    switch (char) {
      case "&":
        return "&amp;";
      case "<":
        return "&lt;";
      case ">":
        return "&gt;";
      case '"':
        return "&quot;";
      default:
        return "&apos;";
    }
  });
}

async function anonymousJSON(path: string): Promise<unknown | null> {
  const response = await fetch(new URL(`/api${path}`, getAPIAddress()), {
    headers: { Accept: "application/json" },
    cache: "no-store",
    signal: AbortSignal.timeout(10_000),
  });
  if (response.status === 401 || response.status === 403) return null;
  if (!response.ok)
    throw new Error(`Public index API returned ${response.status}`);
  return response.json();
}

function pageCount(value: unknown): number {
  if (!Number.isSafeInteger(value) || (value as number) < 0) {
    throw new Error("Invalid public index page count");
  }
  return value as number;
}

export async function getIndexPage(
  kind: IndexKind,
  page: number,
): Promise<IndexPage | null> {
  if (!Number.isSafeInteger(page) || page < 1 || page > 1_000_000) {
    throw new Error("Invalid public index page number");
  }
  const path =
    kind === "threads"
      ? `/threads?page=${page}&visibility=published&ignore_pinned=true`
      : `/nodes/index?page=${page}`;
  const data = await anonymousJSON(path);
  if (data === null) return null;
  if (typeof data !== "object" || data === null)
    throw new Error("Invalid public index response");
  const result = data as Record<string, unknown>;
  const list = kind === "threads" ? result["threads"] : result["nodes"];
  if (!Array.isArray(list)) throw new Error("Invalid public index entries");
  const limit =
    kind === "threads" ? API_PAGE_SIZE_THREADS : API_PAGE_SIZE_NODES;
  if (list.length > limit)
    throw new Error("Public index page exceeds size limit");

  const entries = list.map((item: unknown): IndexEntry => {
    if (typeof item !== "object" || item === null)
      throw new Error("Invalid public index entry");
    const row = item as Record<string, unknown>;
    const slug = row["slug"];
    const updatedAt = kind === "threads" ? row["updatedAt"] : row["updated_at"];
    if (
      typeof slug !== "string" ||
      !slug ||
      typeof updatedAt !== "string" ||
      Number.isNaN(Date.parse(updatedAt))
    ) {
      throw new Error("Invalid public index entry fields");
    }
    if (kind === "threads" && row["visibility"] !== "published") {
      throw new Error("Thread index returned non-published content");
    }
    const lastReplyAt = kind === "threads" ? row["last_reply_at"] : undefined;
    if (
      lastReplyAt !== undefined &&
      (typeof lastReplyAt !== "string" || Number.isNaN(Date.parse(lastReplyAt)))
    ) {
      throw new Error("Invalid thread reply timestamp");
    }
    return {
      slug,
      updatedAt:
        typeof lastReplyAt === "string" &&
        Date.parse(lastReplyAt) > Date.parse(updatedAt)
          ? lastReplyAt
          : updatedAt,
    };
  });
  return { entries, totalPages: pageCount(result["total_pages"]) };
}

export function sitemapIndexXML(pages: Record<IndexKind, number>): string {
  if (
    pages.threads < 0 ||
    pages.library < 0 ||
    pages.threads + pages.library > 50_000
  ) {
    throw new Error("Sitemap index exceeds protocol limit");
  }
  const urls = (["threads", "library"] as const).flatMap((kind) =>
    Array.from({ length: pages[kind] }, (_, i) =>
      canonicalURL(`/sitemaps/${kind}/${i + 1}.xml`),
    ),
  );
  return `<?xml version="1.0" encoding="UTF-8"?>\n<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${urls.map((url) => `  <sitemap><loc>${escapeXML(url)}</loc></sitemap>`).join("\n")}\n</sitemapindex>\n`;
}

export function urlsetXML(kind: IndexKind, entries: IndexEntry[]): string {
  const prefix = kind === "threads" ? "/t/" : "/l/";
  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${entries.map((entry) => `  <url><loc>${escapeXML(canonicalURL(prefix + encodeURIComponent(entry.slug)))}</loc><lastmod>${escapeXML(new Date(entry.updatedAt).toISOString())}</lastmod></url>`).join("\n")}\n</urlset>\n`;
}

export const XML_HEADERS = {
  "Content-Type": "application/xml; charset=utf-8",
  "Cache-Control": "no-store",
};
