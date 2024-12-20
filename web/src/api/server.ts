import { cookies } from "next/headers";

import { buildRequest, buildResult } from "./common";

// Server side variant of fetcher that includes SSR cookies.

type Options = RequestInit & {
  method: string;
};

// Orval fetch generated code is a bit different to SWR fetcher for some reason.
type Result<T> = {
  data: T;
  status: number;
};

export const fetcher = async <T>(url: string, opts: Options): Promise<T> => {
  const headers = Object.fromEntries(new Headers(opts.headers).entries());

  const req = buildRequest({
    url,
    headers,
    method: opts.method as any,
    data: opts.body,
    revalidate: 5,
  });

  req.headers.set("Cookie", await getCookieHeader());

  const response = await fetch(req);
  const result = await buildResult<T>(response);

  // Orval generated types are incorrect here. For some reason it generates a
  // struct with a `data` field, but the actual result type is just the data.
  // However the generated caller passes T as Promise<T> so we need to cast it.
  return { data: result, status: response.status } as T;
};

async function getCookieHeader(): Promise<string> {
  const c = await cookies();
  return c
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");
}
