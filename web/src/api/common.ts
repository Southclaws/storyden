import { filter, flow, reduce, toPairs } from "lodash/fp";
import { z } from "zod";

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

export const ProblemDetailsSchema = z.object({
  trace_id: z.string().min(1),
  type: z.string().optional(),
  title: z.string().optional(),
  detail: z.string().optional(),
  metadata: z.unknown().optional(),
});
export type ProblemDetails = z.infer<typeof ProblemDetailsSchema>;

const OAuthErrorSchema = z.object({
  error: z.string().min(1),
  error_description: z.string().optional(),
});

type OAuthError = z.infer<typeof OAuthErrorSchema>;

const genericErrorMessage = `An unexpected error occurred.`;

export class RequestError extends Error {
  public status: number;
  public problem?: ProblemDetails;

  constructor(message: string, status: number, problem?: ProblemDetails) {
    super(message);
    this.status = status;
    this.problem = problem;
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
    const data = (await response.json().catch(() => undefined)) as unknown;

    const err = normaliseRequestError(data);

    throw new RequestError(
      err.title ?? genericErrorMessage,
      response.status,
      err,
    );
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

function normaliseRequestError(e: unknown): ProblemDetails {
  const problem = parseProblemDetails(e);
  if (problem) {
    return problem;
  }

  const oauthError = parseOAuthError(e);
  if (oauthError) {
    return {
      trace_id: "unknown",
      type: oauthProblemType(oauthError.error),
      title: oauthError.error_description ?? genericErrorMessage,
      detail: `OAuth error: ${oauthError.error}`,
    };
  }

  return {
    trace_id: "unknown",
    title: "An unexpected error occurred.",
    detail: "unknown error",
  };
}

export function parseProblemDetails(data: unknown): ProblemDetails | undefined {
  const result = ProblemDetailsSchema.safeParse(data);

  return result.success ? result.data : undefined;
}

function parseOAuthError(data: unknown): OAuthError | undefined {
  const result = OAuthErrorSchema.safeParse(data);

  return result.success ? result.data : undefined;
}

function oauthProblemType(code: string): string {
  return `urn:storyden:problem:oauth:${slugProblemCode(code)}`;
}

function slugProblemCode(code: string): string {
  return code.toLowerCase().replaceAll("_", "-").replaceAll(" ", "-");
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
