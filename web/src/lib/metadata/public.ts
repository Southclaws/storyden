import type { Metadata } from "next";

import { RequestError } from "@/api/common";
import { WEB_ADDRESS } from "@/config";

export function isNotPublic(error: unknown): boolean {
  return (
    error instanceof RequestError && [401, 403, 404].includes(error.status)
  );
}

export type PublicContent<T> =
  | { status: "public"; data: T }
  | { status: "not-public" }
  | { status: "unavailable"; error: unknown };

// A failed SEO read cannot establish that content is private or missing.
export async function readPublicContent<T>(
  read: () => Promise<T | undefined>,
): Promise<PublicContent<T>> {
  try {
    const data = await read();
    if (data) {
      return { status: "public", data };
    }

    return { status: "not-public" };
  } catch (error) {
    if (isNotPublic(error)) {
      return { status: "not-public" };
    }

    return { status: "unavailable", error };
  }
}

export function canonicalPage(path: string, page = 1): string {
  const url = new URL(path, WEB_ADDRESS);

  if (page > 1) {
    url.searchParams.set("page", String(page));
  }

  return url.toString();
}

export function socialImage(
  path = "/opengraph-image",
  alt = "Community preview",
) {
  return {
    url: new URL(path, WEB_ADDRESS).toString(),
    width: 1200,
    height: 630,
    type: "image/png",
    alt,
  };
}

export function pageLinks(
  path: string,
  page: number,
  totalPages: number,
): Metadata["pagination"] {
  const links: Metadata["pagination"] = {};

  if (page > 1) {
    links.previous = canonicalPage(path, page - 1);
  }

  if (page < totalPages) {
    links.next = canonicalPage(path, page + 1);
  }

  return links;
}

export function jsonLD(data: object): string {
  // Escape HTML delimiters so content cannot close the script element.
  return JSON.stringify(data)
    .replace(/</g, "\\u003c")
    .replace(/>/g, "\\u003e")
    .replace(/&/g, "\\u0026")
    .replace(/\u2028/g, "\\u2028")
    .replace(/\u2029/g, "\\u2029");
}
