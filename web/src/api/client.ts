import { isNil, omitBy } from "lodash/fp";

import { API_ADDRESS } from "src/config";

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
    // NOTE: this is forced uppercase due to a bug somewhere in another part of
    // the code. It might be Orval as it generates all lowercase methods for all
    // requests, however this seems to work fine with every operation except the
    // PATCH calls. Not really sure if that's a browser issue or not though...
    method: method.toUpperCase(),
    mode: "cors",
    credentials: "include",
    ...(headers ? { headers } : {}),
    body: buildPayload(data),
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

  // NOTE: The API code generator returns empty responses where there is no
  // response type specified with a content type so this is the easy way to
  // escape that code path and exit easily.
  if (response.headers.get("content-length") === "0") {
    return undefined as T;
  }

  return response.json();
};

export const server = async <T>(url: string, options?: Options): Promise<T> => {
  const req = new Request(
    `${API_ADDRESS}/api/${url}${cleanQuery(options?.params)}`,
    {
      method: "GET",
      mode: "cors",
      credentials: "include",
    }
  );

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

  // NOTE: The API code generator returns empty responses where there is no
  // response type specified with a content type so this is the easy way to
  // escape that code path and exit easily.
  if (response.headers.get("content-length") === "0") {
    return undefined as T;
  }

  return response.json();
};

const buildPayload = (data: unknown) => {
  if (!data) {
    return undefined;
  }

  if (data instanceof File) {
    return data;
  }

  return JSON.stringify(data);
};

export default fetcher;

const removeEmpty = omitBy(isNil);

const cleanQuery = (params?: Record<string, string | string[]>): string => {
  if (!params) return "";

  const clean = removeEmpty(params);

  const format = new URLSearchParams(clean).toString();

  return `?${format}`;
};
