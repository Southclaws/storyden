import "server-only";
import TurndownService from "turndown";

import { getAPIAddress } from "@/config";

export class PublicContentError extends Error {
  constructor(public status: number) {
    super(`Public content request failed: ${status}`);
  }
}

export async function publicJSON<T>(path: string): Promise<T> {
  let response: Response;
  try {
    response = await fetch(new URL(`/api${path}`, getAPIAddress()), {
      headers: { Accept: "application/json" },
      credentials: "omit",
      cache: "no-store",
      signal: AbortSignal.timeout(5000),
    });
  } catch {
    throw new PublicContentError(503);
  }

  if (!response.ok) {
    throw new PublicContentError(
      [401, 403, 404].includes(response.status) ? response.status : 503,
    );
  }

  try {
    return (await response.json()) as T;
  } catch {
    throw new PublicContentError(503);
  }
}

export function markdownResponse(
  body: string,
  status = 200,
  canonical?: string,
  markdown?: string,
) {
  const headers = new Headers({
    "Content-Type": "text/markdown; charset=utf-8",
    "Cache-Control": "no-store",
  });
  if (canonical && markdown) {
    headers.set(
      "Link",
      `<${canonical}>; rel="canonical"; type="text/html", <${markdown}>; rel="alternate"; type="text/markdown"`,
    );
  }
  return new Response(status === 200 ? `${body.trim()}\n` : body, {
    status,
    headers,
  });
}

export function errorResponse(error: unknown) {
  const status = error instanceof PublicContentError ? error.status : 503;
  return markdownResponse(`# ${status}\n`, status);
}

export function pageNumber(request: Request): number {
  const raw = new URL(request.url).searchParams.get("page") ?? "1";
  if (!/^[1-9][0-9]{0,6}$/.test(raw)) throw new PublicContentError(400);
  const page = Number(raw);
  if (page > 1_000_000) throw new PublicContentError(400);
  return page;
}

export function inline(value: string): string {
  return value.replace(/[\r\n\t]+/g, " ").trim();
}

export function htmlToPublicMarkdown(
  html: string,
  addresses?: { web: string; api: string; base?: string },
): string {
  const turndown = new TurndownService({
    headingStyle: "atx",
    codeBlockStyle: "fenced",
    bulletListMarker: "-",
  });
  turndown.remove(["script", "style", "iframe", "form"]);
  turndown.addRule("tables", {
    filter: "table",
    replacement(_content, node) {
      const rows = Array.from(
        (node as HTMLTableElement).querySelectorAll("tr"),
      ).map((row) =>
        Array.from(row.querySelectorAll("th, td")).map((cell) =>
          inline(cell.textContent ?? "").replaceAll("|", "\\|"),
        ),
      );
      if (!rows.length) return "";
      const width = Math.max(...rows.map((row) => row.length));
      const rowText = (cells: string[]) =>
        `| ${Array.from({ length: width }, (_, index) => cells[index] ?? "").join(" | ")} |`;
      return `\n\n${[
        rowText(rows[0]!),
        rowText(Array(width).fill("---")),
        ...rows.slice(1).map(rowText),
      ].join("\n")}\n\n`;
    },
  });
  turndown.addRule("safeLinks", {
    filter: "a",
    replacement(content, node) {
      const href = (node as HTMLAnchorElement).getAttribute("href") ?? "";
      const resolved = resolveLink(href, addresses);
      return resolved ? `[${content}](${resolved})` : content;
    },
  });
  turndown.addRule("safeImages", {
    filter: "img",
    replacement(_content, node) {
      const image = node as HTMLImageElement;
      const src = image.getAttribute("src") ?? "";
      const resolved = resolveLink(src, addresses);
      return resolved && !resolved.startsWith("#")
        ? `![${inline(image.getAttribute("alt") ?? "")}](${resolved})`
        : "";
    },
  });
  return turndown.turndown(html).trim();
}

function resolveLink(
  value: string,
  addresses?: { web: string; api: string; base?: string },
): string | null {
  if (/^https?:\/\//i.test(value) || value.startsWith("#")) return value;
  if (/^mailto:/i.test(value)) return value;
  if (!addresses) return value.startsWith("/") ? value : null;
  const native = /^sdr:([a-z]+)\/([a-z0-9]+)$/i.exec(value);
  if (native) {
    return publicURL(addresses.web, `/_/resolve/${native[1]}/${native[2]}`);
  }
  if (/^[a-z][a-z0-9+.-]*:/i.test(value)) return null;
  if (value.startsWith("/api/")) return publicURL(addresses.api, value);
  if (value.startsWith("/")) return publicURL(addresses.web, value);
  return publicURL(addresses.base ?? addresses.web, value);
}

export function publicURL(origin: string, path: string): string {
  return new URL(path, origin).href;
}
