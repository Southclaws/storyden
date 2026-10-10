import type { Metadata } from "next";

import type { Node, Thread } from "@/api/openapi-schema";
import { WEB_ADDRESS, getAPIAddress } from "@/config";

export function publicURL(path: string): string {
  return new URL(path, `${WEB_ADDRESS.replace(/\/$/, "")}/`).toString();
}

export function canonicalPage(path: string, rawPage: unknown): string {
  const page = pageNumber(rawPage);
  const suffix = page !== undefined && page > 1 ? `?page=${page}` : "";
  return publicURL(`${path}${suffix}`);
}

export function pageNumber(rawPage: unknown): number | undefined {
  if (typeof rawPage !== "string" || !/^[1-9]\d*$/.test(rawPage)) return;
  const page = Number(rawPage);
  return Number.isSafeInteger(page) ? page : undefined;
}

export function pageLinks(
  path: string,
  page: number,
  totalPages: number,
): Metadata["pagination"] {
  return {
    ...(page > 1 ? { previous: canonicalPage(path, String(page - 1)) } : {}),
    ...(page < totalPages
      ? { next: canonicalPage(path, String(page + 1)) }
      : {}),
  };
}

export const hiddenMetadata: Metadata = {
  robots: { index: false, follow: false },
};

/** A metadata read must never inherit a visitor's cookies or authorization. */
export async function getAnonymousPublicResource<
  T extends { visibility: string },
>(path: string): Promise<T | undefined> {
  const resource = await getAnonymousJSON<T>(path);
  return resource?.visibility === "published" ? resource : undefined;
}

export async function getAnonymousPageCount(
  path: string,
): Promise<number | undefined> {
  const data = await getAnonymousJSON<{ total_pages?: unknown }>(path);
  return typeof data?.total_pages === "number" &&
    Number.isSafeInteger(data.total_pages) &&
    data.total_pages > 0
    ? data.total_pages
    : undefined;
}

async function getAnonymousJSON<T>(path: string): Promise<T | undefined> {
  try {
    const url = new URL(
      `/api${path}`,
      `${getAPIAddress().replace(/\/$/, "")}/`,
    );
    const response = await fetch(url, {
      cache: "no-store",
      credentials: "omit",
      signal: AbortSignal.timeout(5000),
    });
    if (!response.ok) return;
    return (await response.json()) as T;
  } catch {
    return;
  }
}

export function getPublicThread(slug: string): Promise<Thread | undefined> {
  return getAnonymousPublicResource<Thread>(
    `/threads/${encodeURIComponent(slug)}`,
  );
}

export function getPublicNode(slug: string): Promise<Node | undefined> {
  return getAnonymousPublicResource<Node>(`/nodes/${encodeURIComponent(slug)}`);
}

export function jsonLD(data: object): string {
  // JSON inside a script tag must not be able to terminate the script element.
  return JSON.stringify(data)
    .replace(/</g, "\\u003c")
    .replace(/>/g, "\\u003e")
    .replace(/&/g, "\\u0026")
    .replace(/\u2028/g, "\\u2028")
    .replace(/\u2029/g, "\\u2029");
}

export function websiteSchema(name: string, description: string) {
  return {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name,
    description,
    url: publicURL("/"),
  };
}

export function threadSchema(thread: Thread) {
  const url = publicURL(`/t/${encodeURIComponent(thread.slug)}`);
  return {
    "@context": "https://schema.org",
    "@type": "DiscussionForumPosting",
    headline: thread.title,
    ...(thread.description ? { description: thread.description } : {}),
    url,
    mainEntityOfPage: url,
    datePublished: thread.createdAt,
    dateModified: thread.updatedAt,
    author: { "@type": "Person", name: thread.author.name },
  };
}

export function libraryPageSchema(node: Node) {
  const url = publicURL(`/l/${encodeURIComponent(node.slug)}`);
  return {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: node.name,
    ...(node.description ? { description: node.description } : {}),
    url,
    datePublished: node.createdAt,
    dateModified: node.updatedAt,
  };
}
