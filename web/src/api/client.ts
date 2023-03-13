type Options = {
  url: string;
  method: "get" | "post" | "put" | "delete" | "patch";
  headers?: Record<string, string>;
  params?: any;
  data?: unknown;
  responseType?: string;
};

export const fetcher = async <T>({
  url,
  method,
  headers,
  params,
  data,
}: Options): Promise<T> => {
  const query = params ? `?${new URLSearchParams(params)}` : "";

  const req = new Request(`/api${url}${query}`, {
    method,
    mode: "cors",
    credentials: "include",
    ...(headers ? { headers } : {}),
    ...(data ? { body: JSON.stringify(data) } : {}),
  });

  const response = await fetch(req);

  if (!response.ok) {
    const data = await response
      .json()
      .catch(() => ({ error: "Failed to parse API response" }));
    console.warn(data);
    throw new Error(
      data.message ?? `An unexpected error occurred: ${response.statusText}`
    );
  }

  return response.json();
};

export default fetcher;
