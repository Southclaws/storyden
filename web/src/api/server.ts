import { cookies } from "next/headers";
import { notFound } from "next/navigation";

import { RequestError, buildRequest, buildResult } from "./common";

// Server side variant of fetcher that includes SSR cookies.

type Options = RequestInit & {
  method: string;
};

// Orval fetch generated code is a bit different to SWR fetcher for some reason.
type Result<T> = {
  data: T;
  status: number;
};

export const fetcher = async <T>(
  url: string,
  { method, ...opts }: Options,
): Promise<T> => {
  const { headers: requestHeaders, ...requestInit } = opts;
  const headers = Object.fromEntries(new Headers(requestHeaders).entries());

  const req = buildRequest({
    url,
    headers,
    method: method as any,
    data: requestInit.body,
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

  const response = await fetch(req);

  try {
    const result = await buildResult<T>(response);

    // Orval generated types are incorrect here. For some reason it generates a
    // struct with a `data` field, but the actual result type is just the data.
    // However the generated caller passes T as Promise<T> so we need to cast it.
    return { data: result, status: response.status } as T;
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
