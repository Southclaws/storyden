import type { Metadata } from "next";

import { WEB_ADDRESS } from "@/config";

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
