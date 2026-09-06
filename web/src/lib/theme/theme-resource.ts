import { createHash } from "node:crypto";

const CACHE_CONTROL = "public, no-cache";

export function createThemeResourceResponse(
  request: Request,
  source: string,
  contentType: "application/javascript" | "text/css",
) {
  const digest = createHash("sha256").update(source).digest("base64");
  const etag = `"sha256-${digest}"`;
  const headers = {
    "Cache-Control": CACHE_CONTROL,
    "Content-Type": `${contentType}; charset=utf-8`,
    ETag: etag,
  };

  if (matchesETag(request.headers.get("If-None-Match"), etag)) {
    return new Response(null, { status: 304, headers });
  }

  return new Response(source, { headers });
}

function matchesETag(header: string | null, etag: string) {
  if (!header) return false;

  return header.split(",").some((candidate) => {
    const trimmed = candidate.trim();
    return (
      trimmed === "*" ||
      trimmed === etag ||
      (trimmed.startsWith("W/") && trimmed.slice(2) === etag)
    );
  });
}
