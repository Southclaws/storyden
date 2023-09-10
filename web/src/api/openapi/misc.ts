/**
 * Generated by orval v6.17.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import useSwr from "swr";
import type { Key, SWRConfiguration } from "swr";

import { fetcher } from "../client";

import type {
  AssetGetOKResponse,
  AssetUploadBody,
  GetInfoOKResponse,
  InternalServerErrorResponse,
} from "./schemas";

type AwaitedInput<T> = PromiseLike<T> | T;

type Awaited<O> = O extends AwaitedInput<infer T> ? T : never;

/**
 * The version number includes the date and time of the release build as
well as a short representation of the Git commit hash.

 * @summary Get the software version string.
 */
export const getVersion = () => {
  return fetcher<string>({ url: `/version`, method: "get" });
};

export const getGetVersionKey = () => [`/version`] as const;

export type GetVersionQueryResult = NonNullable<
  Awaited<ReturnType<typeof getVersion>>
>;
export type GetVersionQueryError = unknown;

/**
 * @summary Get the software version string.
 */
export const useGetVersion = <TError = unknown>(options?: {
  swr?: SWRConfiguration<Awaited<ReturnType<typeof getVersion>>, TError> & {
    swrKey?: Key;
    enabled?: boolean;
  };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getGetVersionKey() : null));
  const swrFn = () => getVersion();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};

/**
 * Note: the generator creates a `map[string]interface{}` if this is set to
`application/json`... so I'm just using plain text for now.

 * @summary Get the OpenAPI 3.0 specification as JSON.
 */
export const getSpec = () => {
  return fetcher<string>({ url: `/openapi.json`, method: "get" });
};

export const getGetSpecKey = () => [`/openapi.json`] as const;

export type GetSpecQueryResult = NonNullable<
  Awaited<ReturnType<typeof getSpec>>
>;
export type GetSpecQueryError = unknown;

/**
 * @summary Get the OpenAPI 3.0 specification as JSON.
 */
export const useGetSpec = <TError = unknown>(options?: {
  swr?: SWRConfiguration<Awaited<ReturnType<typeof getSpec>>, TError> & {
    swrKey?: Key;
    enabled?: boolean;
  };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getGetSpecKey() : null));
  const swrFn = () => getSpec();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};

/**
 * Get the basic forum installation info such as title, description, etc.

 */
export const getInfo = () => {
  return fetcher<GetInfoOKResponse>({ url: `/v1/info`, method: "get" });
};

export const getGetInfoKey = () => [`/v1/info`] as const;

export type GetInfoQueryResult = NonNullable<
  Awaited<ReturnType<typeof getInfo>>
>;
export type GetInfoQueryError = InternalServerErrorResponse;

export const useGetInfo = <TError = InternalServerErrorResponse>(options?: {
  swr?: SWRConfiguration<Awaited<ReturnType<typeof getInfo>>, TError> & {
    swrKey?: Key;
    enabled?: boolean;
  };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getGetInfoKey() : null));
  const swrFn = () => getInfo();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};

/**
 * Get the logo icon image.
 */
export const iconGet = (
  iconSize: "512x512" | "32x32" | "180x180" | "120x120" | "167x167" | "152x152",
) => {
  return fetcher<AssetGetOKResponse>({
    url: `/v1/info/icon/${iconSize}`,
    method: "get",
  });
};

export const getIconGetKey = (
  iconSize: "512x512" | "32x32" | "180x180" | "120x120" | "167x167" | "152x152",
) => [`/v1/info/icon/${iconSize}`] as const;

export type IconGetQueryResult = NonNullable<
  Awaited<ReturnType<typeof iconGet>>
>;
export type IconGetQueryError = InternalServerErrorResponse;

export const useIconGet = <TError = InternalServerErrorResponse>(
  iconSize: "512x512" | "32x32" | "180x180" | "120x120" | "167x167" | "152x152",
  options?: {
    swr?: SWRConfiguration<Awaited<ReturnType<typeof iconGet>>, TError> & {
      swrKey?: Key;
      enabled?: boolean;
    };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!iconSize;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getIconGetKey(iconSize) : null));
  const swrFn = () => iconGet(iconSize);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};

/**
 * Upload and process the installation's logo image.
 */
export const iconUpload = (assetUploadBody: AssetUploadBody) => {
  return fetcher<void>({
    url: `/v1/info/icon`,
    method: "post",
    headers: { "Content-Type": "application/octet-stream" },
    data: assetUploadBody,
  });
};
