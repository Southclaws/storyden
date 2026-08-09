import { cookies, headers as nextHeaders } from "next/headers";
import { notFound } from "next/navigation";

import { Options, RequestError, buildRequest, buildResult } from "./common";

type ResponseData<T> = T extends { data: infer Data } ? Data : never;

const SSR_NETWORK_HEADER_NAMES = ["x-forwarded-for"] as const;

export const fetcher = async <T>(
  url: string,
  { method, ...opts }: Options,
): Promise<T> => {
  const { headers: requestHeaders, ...requestInit } = opts;
  const headers = Object.fromEntries(new Headers(requestHeaders).entries());

  const req = buildRequest(url, {
    headers,
    method: method as any,
    body: requestInit.body,
    // Server side requests are cached a little more aggressively than client
    // side hydration requests. The downside of this is a user may see a flash
    // of stale data as the server render loads which will be replaced by the
    // client side hydration by SWR. However, the second call will most likely
    // be a 304 if it has been loaded before by the same user already. So, in a
    // best case, we get a single database read, worst case we get two. The
    // revalidation period is set to one minute in order to cut down on the
    // flashes of stale data. However, in reality this doesn't really gain much
    // as a user landing for the first time will still trigger two DB reads,
    // and a user returning is quite likely someone who has interacted with a
    // piece of content and thus will result in a new read at least once. So,
    // it's not the most efficient approach (ignoring server-side data cache)
    // but it's the best of a not-so-great situation. This should improve a lot
    // if Next.js adds support for HTTP Conditional Requests and ETag headers.
    revalidate: 60,
    cache: "force-cache",
    ...requestInit,
  });

  req.headers.set("Cookie", await getCookieHeader());
  await applyForwardedClientHeaders(req.headers);

  const response = await fetch(req);

  try {
    const data = await buildResult<ResponseData<T>>(response);

    return {
      data,
      status: response.status,
      headers: response.headers,
    } as T;
  } catch (e) {
    if (e instanceof RequestError && e.status === 404) {
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
