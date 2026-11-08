import { cookies, headers as nextHeaders } from "next/headers";
import { notFound } from "next/navigation";

import { Options, RequestError, buildRequest, buildResult } from "./common";

type ResponseData<T> = T extends { data: infer Data } ? Data : never;

const SSR_NETWORK_HEADER_NAMES = ["x-forwarded-for"] as const;

export const fetcher = async <T>(
  url: string,
  { method, anonymous = false, ...opts }: Options & { anonymous?: boolean },
): Promise<T> => {
  const { headers: requestHeaders, ...requestInit } = opts;
  const headers = Object.fromEntries(new Headers(requestHeaders).entries());

  const req = buildRequest(url, {
    headers,
    method: method as any,
    body: requestInit.body,
    ...requestInit,
    // Browser mutations do not invalidate Next.js's persistent API cache.
    cache: "no-store",
  });

  if (anonymous) {
    req.headers.delete("Cookie");
    req.headers.delete("Authorization");
    // Without the visitor's address, every anonymous SSR read shares the
    // frontend server's rate limit bucket and starves under crawler traffic.
    await applyForwardedClientHeaders(req.headers).catch(() => {});
  } else {
    req.headers.set("Cookie", await getCookieHeader());
    await applyForwardedClientHeaders(req.headers);
  }

  const response = await fetch(req);

  try {
    const data = await buildResult<ResponseData<T>>(response);

    return {
      data,
      status: response.status,
      headers: response.headers,
    } as T;
  } catch (e) {
    if (e instanceof RequestError && e.status === 404 && !anonymous) {
      notFound();
    }
    throw e;
  }
};

async function getCookieHeader(): Promise<string> {
  const c = await cookies();
  return c
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");
}

async function applyForwardedClientHeaders(out: Headers): Promise<void> {
  const h = await nextHeaders();
  out.set("x-storyden-ssr", "1");
  forwardSSRNetworkHeaders(out, h);
}

export function forwardSSRNetworkHeaders(
  out: Headers,
  incoming: Headers,
): void {
  for (const headerName of SSR_NETWORK_HEADER_NAMES) {
    const normalisedHeaderName = headerName.trim().toLowerCase();
    if (!normalisedHeaderName) {
      continue;
    }

    const value = incoming.get(normalisedHeaderName)?.trim();
    if (!value) {
      continue;
    }

    // Forward canonical network headers as-is.
    out.set(normalisedHeaderName, value);
  }
}
