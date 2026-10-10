import "server-only";

import { getAPIAddress } from "@/config";

const responseHeaders = [
  "Content-Type",
  "Cache-Control",
  "ETag",
  "Access-Control-Allow-Origin",
  "Access-Control-Allow-Methods",
  "Access-Control-Allow-Headers",
  "Access-Control-Expose-Headers",
  "Link",
];

export async function proxyDiscovery(request: Request, path: string) {
  const headers = new Headers();
  const etag = request.headers.get("If-None-Match");
  if (etag) headers.set("If-None-Match", etag);

  try {
    const upstream = await fetch(new URL(path, getAPIAddress()), {
      method: request.method,
      headers,
      cache: "no-store",
      signal: AbortSignal.timeout(5000),
    });
    const publicHeaders = new Headers();
    for (const header of responseHeaders) {
      const value = upstream.headers.get(header);
      if (value) publicHeaders.set(header, value);
    }

    return new Response(
      upstream.status === 204 ||
        upstream.status === 304 ||
        request.method === "HEAD"
        ? null
        : upstream.body,
      { status: upstream.status, headers: publicHeaders },
    );
  } catch {
    return new Response("Discovery service unavailable", {
      status: 502,
      headers: { "Cache-Control": "no-store" },
    });
  }
}
