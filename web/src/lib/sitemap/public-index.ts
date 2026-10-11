import { escape } from "lodash";

import { RequestError } from "@/api/common";
import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";
import { WEB_ADDRESS } from "@/config";

export type IndexKind = "threads" | "library";

export type IndexEntry = {
  slug: string;
  updatedAt: string;
};

export type IndexPage = {
  entries: IndexEntry[];
  totalPages: number;
};

const MAX_SITEMAP_ENTRIES = 50_000;

export const XML_HEADERS = {
  "Content-Type": "application/xml; charset=utf-8",
  "Cache-Control": "no-store",
};

export async function getIndexPage(
  kind: IndexKind,
  page: number,
): Promise<IndexPage | null> {
  if (!Number.isSafeInteger(page) || page < 1) {
    throw new Error("Invalid sitemap page number");
  }

  try {
    switch (kind) {
      case "threads":
        return await getThreadIndexPage(page);

      case "library":
        return await getLibraryIndexPage(page);
    }
  } catch (error) {
    if (
      error instanceof RequestError &&
      (error.status === 401 || error.status === 403)
    ) {
      return null;
    }

    throw error;
  }
}

async function getThreadIndexPage(page: number): Promise<IndexPage> {
  const { data } = await threadList(
    { page: String(page), visibility: ["published"], ignore_pinned: true },
    { anonymous: true, cache: "no-store", signal: AbortSignal.timeout(10_000) },
  );

  if (data.threads.some((thread) => thread.visibility !== "published")) {
    throw new Error("Thread list returned non-published content");
  }

  const entries = data.threads.map((thread) => {
    let updatedAt = thread.updatedAt;

    if (
      thread.last_reply_at &&
      Date.parse(thread.last_reply_at) > Date.parse(updatedAt)
    ) {
      updatedAt = thread.last_reply_at;
    }

    return { slug: thread.slug, updatedAt };
  });

  return { entries, totalPages: data.total_pages };
}

async function getLibraryIndexPage(page: number): Promise<IndexPage> {
  const { data } = await nodeList(
    { page: String(page), visibility: ["published"], format: "flat" },
    { anonymous: true, cache: "no-store", signal: AbortSignal.timeout(10_000) },
  );

  if (data.nodes.some((node) => node.visibility !== "published")) {
    throw new Error("Node list returned non-published content");
  }

  const entries = data.nodes.map(({ slug, updatedAt }) => ({
    slug,
    updatedAt,
  }));

  return { entries, totalPages: data.total_pages };
}

export function canonicalURL(path: string): string {
  return new URL(path, WEB_ADDRESS).toString();
}

export function sitemapIndexXML(pages: Record<IndexKind, number>): string {
  const pageCounts = Object.values(pages);
  const validCounts = pageCounts.every(
    (count) => Number.isSafeInteger(count) && count >= 0,
  );
  const totalPages = pages.threads + pages.library;

  if (!validCounts || totalPages > MAX_SITEMAP_ENTRIES) {
    throw new Error("Sitemap index exceeds protocol limit");
  }

  const sitemaps: string[] = [];

  for (const kind of ["threads", "library"] as const) {
    for (let page = 1; page <= pages[kind]; page++) {
      const url = canonicalURL(`/sitemaps/${kind}/${page}.xml`);

      sitemaps.push(`  <sitemap>
    <loc>${escape(url)}</loc>
  </sitemap>`);
    }
  }

  return `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${sitemaps.join("\n")}
</sitemapindex>
`;
}

export function urlsetXML(kind: IndexKind, entries: IndexEntry[]): string {
  if (entries.length > MAX_SITEMAP_ENTRIES) {
    throw new Error("Sitemap exceeds protocol limit");
  }

  const prefix = kind === "threads" ? "/t/" : "/l/";
  const urls = entries.map((entry) => {
    const url = canonicalURL(prefix + encodeURIComponent(entry.slug));
    const lastModified = new Date(entry.updatedAt).toISOString();

    return `  <url>
    <loc>${escape(url)}</loc>
    <lastmod>${lastModified}</lastmod>
  </url>`;
  });

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.join("\n")}
</urlset>
`;
}
