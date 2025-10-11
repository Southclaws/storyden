import { filter, flow, reduce, toPairs } from "lodash/fp";

import { getAPIAddress } from "@/config";

export type Options = {
  url: string;
  method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  headers?: Record<string, string>;
  params?: Record<string, string | string[] | boolean>;
  data?: unknown;
  responseType?: string;
  cookie?: string;
  revalidate?: number;
  cache?: RequestCache;
};

export class RequestError extends Error {
  public status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export function buildRequest({
  url,
  method = "GET",
  headers,
  params,
  data,
  revalidate,
  cache,
}: Options): Request {
  const apiAddress = getAPIAddress();
  const address = `${apiAddress}/api${url}${cleanQuery(params)}`;
  const _method = method.toUpperCase();

  const tags = buildNextTagsFromURL(address);

  return new Request(address, {
    method: _method,
    mode: "cors",
    credentials: "include",
    headers,
    body: buildPayload(data),
    cache,
    next: {
      tags,
      revalidate,
    },
  });
}

export async function buildResult<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const data = await response.json().catch(() => undefined);
    const fallback = `${response.status} ${response.statusText}`;

    throw new RequestError(data?.message ?? fallback, response.status);
  }

  // NOTE: The API code generator returns empty responses where there is no
  // response type specified with a content type so this is the easy way to
  // escape that code path and exit easily.
  if (response.headers.get("content-length") === "0") {
    return undefined as T;
  }

  if (response.headers.get("content-type")?.includes("json")) {
    return response.json();
  }

  return response.blob() as T;
}

export const buildPayload = (data: unknown) => {
  if (!data) {
    return undefined;
  }

  if (data instanceof File || data instanceof Blob) {
    return data;
  }

  if (typeof data === "string") {
    return data;
  }

  return JSON.stringify(data);
};

type ParameterValue = string | string[] | boolean;

type QueryParameters = Record<string, ParameterValue>;

const removeEmpty = filter<[string, ParameterValue]>(
  ([, v]) => v !== null && v !== undefined && v !== "",
);

// NOTE: This is for correctly formatting arrays as multiple instances of the
// same key-value pair (the correct way to send arrays in query parameters).
const objectToParams = (init: URLSearchParams) =>
  reduce((acc: URLSearchParams, [k, v]: [string, ParameterValue]) => {
    if (Array.isArray(v)) {
      v.forEach((vv) => acc.append(k, vv));
    } else if (typeof v === "boolean") {
      acc.append(k, v ? "true" : "false");
    } else {
      acc.append(k, v);
    }

    return acc;
  }, init);

const processQueryParams = (init: URLSearchParams) =>
  flow(toPairs, removeEmpty, objectToParams(init));

export const cleanQuery = (params?: QueryParameters): string => {
  if (!params) return "";

  const usp = processQueryParams(new URLSearchParams())(params);

  const format = usp.toString();

  if (!format) return "";

  return `?${format}`;
};

export function shouldLog(status: number) {
  if (status < 400) {
    return false;
  }

  if (status === 404) {
    return false;
  }

  return true;
}

function buildNextTagsFromURL(url: string) {
  const u = new URL(url);

  const segments = u.pathname.split("/").filter(Boolean);

  return segments;
}
