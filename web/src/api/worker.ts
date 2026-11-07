type Options = RequestInit & { apiAddress?: string };

export const fetcher = async <T>(
  url: string,
  { apiAddress, ...options }: Options
): Promise<T> => {
  if (!apiAddress) {
    throw new Error("Worker API requests require the public API address");
  }

  const request = new Request(`${apiAddress.replace(/\/$/, "")}/api${url}`, {
    ...options,
    credentials: "include",
    cache: "no-store",
    redirect: "error",
    signal: options.signal ?? AbortSignal.timeout(10000),
  });

  const response = await fetch(request);

  const body = await response.text();
  const data = body ? JSON.parse(body) : undefined;

  return { data, status: response.status, headers: response.headers } as T;
};
