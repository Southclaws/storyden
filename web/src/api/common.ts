import { filter, flow, reduce, toPairs } from "lodash/fp";

export type Options = {
  url: string;
  method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  headers?: Record<string, string>;
  params?: Record<string, string | string[] | boolean>;
  data?: unknown;
  responseType?: string;
  cookie?: string;
};

export const buildPayload = (data: unknown) => {
  if (!data) {
    return undefined;
  }

  if (data instanceof File || data instanceof Blob) {
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
