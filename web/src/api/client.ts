import { isNil, omitBy } from "lodash/fp";

type Options = {
  url: string;
  method: "get" | "post" | "put" | "delete" | "patch";
  headers?: Record<string, string>;
  params?: Record<string, string | string[]>;
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
  const req = new Request(`/api${url}${cleanQuery(params)}`, {
    method,
    // mode: "cors",
    // credentials: "include",
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

const removeEmpty = omitBy(isNil);

const cleanQuery = (params?: Record<string, string | string[]>): string => {
  if (!params) return "";

  const clean = removeEmpty(params);

  const format = new URLSearchParams(clean).toString();

  return `?${format}`;
};
