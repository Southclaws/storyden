/**
 * Generated by orval v6.28.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import useSwr from "swr";
import type { Arguments, Key, SWRConfiguration } from "swr";
import useSWRMutation from "swr/mutation";
import type { SWRMutationConfiguration } from "swr/mutation";

import { fetcher } from "../client";

import type {
  InternalServerErrorResponse,
  ItemCreateBody,
  ItemCreateOKResponse,
  ItemGetOKResponse,
  ItemListOKResponse,
  ItemListParams,
  ItemUpdateBody,
  ItemUpdateOKResponse,
  NotFoundResponse,
  UnauthorisedResponse,
  VisibilityUpdateBody,
} from "./schemas";

/**
 * Create a item to represent a piece of structured data such as an item in
a video game, an article of clothing, a product in a store, etc.

 */
export const itemCreate = (itemCreateBody: ItemCreateBody) => {
  return fetcher<ItemCreateOKResponse>({
    url: `/v1/items`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: itemCreateBody,
  });
};

export const getItemCreateMutationFetcher = () => {
  return (
    _: string,
    { arg }: { arg: ItemCreateBody },
  ): Promise<ItemCreateOKResponse> => {
    return itemCreate(arg);
  };
};
export const getItemCreateMutationKey = () => `/v1/items` as const;

export type ItemCreateMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemCreate>>
>;
export type ItemCreateMutationError =
  | UnauthorisedResponse
  | InternalServerErrorResponse;

export const useItemCreate = <
  TError = UnauthorisedResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof itemCreate>>,
    TError,
    string,
    ItemCreateBody,
    Awaited<ReturnType<typeof itemCreate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getItemCreateMutationKey();
  const swrFn = getItemCreateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * List all items using the filtering options.
 */
export const itemList = (params?: ItemListParams) => {
  return fetcher<ItemListOKResponse>({
    url: `/v1/items`,
    method: "GET",
    params,
  });
};

export const getItemListKey = (params?: ItemListParams) =>
  [`/v1/items`, ...(params ? [params] : [])] as const;

export type ItemListQueryResult = NonNullable<
  Awaited<ReturnType<typeof itemList>>
>;
export type ItemListQueryError = NotFoundResponse | InternalServerErrorResponse;

export const useItemList = <
  TError = NotFoundResponse | InternalServerErrorResponse,
>(
  params?: ItemListParams,
  options?: {
    swr?: SWRConfiguration<Awaited<ReturnType<typeof itemList>>, TError> & {
      swrKey?: Key;
      enabled?: boolean;
    };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getItemListKey(params) : null));
  const swrFn = () => itemList(params);

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
 * Get a item by its URL slug.
 */
export const itemGet = (itemSlug: string) => {
  return fetcher<ItemGetOKResponse>({
    url: `/v1/items/${itemSlug}`,
    method: "GET",
  });
};

export const getItemGetKey = (itemSlug: string) =>
  [`/v1/items/${itemSlug}`] as const;

export type ItemGetQueryResult = NonNullable<
  Awaited<ReturnType<typeof itemGet>>
>;
export type ItemGetQueryError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemGet = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  options?: {
    swr?: SWRConfiguration<Awaited<ReturnType<typeof itemGet>>, TError> & {
      swrKey?: Key;
      enabled?: boolean;
    };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!itemSlug;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getItemGetKey(itemSlug) : null));
  const swrFn = () => itemGet(itemSlug);

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
 * Update a item.
 */
export const itemUpdate = (
  itemSlug: string,
  itemUpdateBody: ItemUpdateBody,
) => {
  return fetcher<ItemUpdateOKResponse>({
    url: `/v1/items/${itemSlug}`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: itemUpdateBody,
  });
};

export const getItemUpdateMutationFetcher = (itemSlug: string) => {
  return (
    _: string,
    { arg }: { arg: ItemUpdateBody },
  ): Promise<ItemUpdateOKResponse> => {
    return itemUpdate(itemSlug, arg);
  };
};
export const getItemUpdateMutationKey = (itemSlug: string) =>
  `/v1/items/${itemSlug}` as const;

export type ItemUpdateMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemUpdate>>
>;
export type ItemUpdateMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemUpdate = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof itemUpdate>>,
      TError,
      string,
      ItemUpdateBody,
      Awaited<ReturnType<typeof itemUpdate>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getItemUpdateMutationKey(itemSlug);
  const swrFn = getItemUpdateMutationFetcher(itemSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Delete an item.
 */
export const itemDelete = (itemSlug: string) => {
  return fetcher<void>({ url: `/v1/items/${itemSlug}`, method: "DELETE" });
};

export const getItemDeleteMutationFetcher = (itemSlug: string) => {
  return (_: string, __: { arg: Arguments }): Promise<void> => {
    return itemDelete(itemSlug);
  };
};
export const getItemDeleteMutationKey = (itemSlug: string) =>
  `/v1/items/${itemSlug}` as const;

export type ItemDeleteMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemDelete>>
>;
export type ItemDeleteMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemDelete = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof itemDelete>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof itemDelete>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getItemDeleteMutationKey(itemSlug);
  const swrFn = getItemDeleteMutationFetcher(itemSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Update the visibility of an item.
 */
export const itemUpdateVisibility = (
  itemSlug: string,
  visibilityUpdateBody: VisibilityUpdateBody,
) => {
  return fetcher<ItemUpdateOKResponse>({
    url: `/v1/items/${itemSlug}/visibility`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: visibilityUpdateBody,
  });
};

export const getItemUpdateVisibilityMutationFetcher = (itemSlug: string) => {
  return (
    _: string,
    { arg }: { arg: VisibilityUpdateBody },
  ): Promise<ItemUpdateOKResponse> => {
    return itemUpdateVisibility(itemSlug, arg);
  };
};
export const getItemUpdateVisibilityMutationKey = (itemSlug: string) =>
  `/v1/items/${itemSlug}/visibility` as const;

export type ItemUpdateVisibilityMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemUpdateVisibility>>
>;
export type ItemUpdateVisibilityMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemUpdateVisibility = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof itemUpdateVisibility>>,
      TError,
      string,
      VisibilityUpdateBody,
      Awaited<ReturnType<typeof itemUpdateVisibility>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getItemUpdateVisibilityMutationKey(itemSlug);
  const swrFn = getItemUpdateVisibilityMutationFetcher(itemSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Add an asset to an item.
 */
export const itemAddAsset = (itemSlug: string, assetId: string) => {
  return fetcher<ItemUpdateOKResponse>({
    url: `/v1/items/${itemSlug}/assets/${assetId}`,
    method: "PUT",
  });
};

export const getItemAddAssetMutationFetcher = (
  itemSlug: string,
  assetId: string,
) => {
  return (_: string, __: { arg: Arguments }): Promise<ItemUpdateOKResponse> => {
    return itemAddAsset(itemSlug, assetId);
  };
};
export const getItemAddAssetMutationKey = (itemSlug: string, assetId: string) =>
  `/v1/items/${itemSlug}/assets/${assetId}` as const;

export type ItemAddAssetMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemAddAsset>>
>;
export type ItemAddAssetMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemAddAsset = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  assetId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof itemAddAsset>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof itemAddAsset>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getItemAddAssetMutationKey(itemSlug, assetId);
  const swrFn = getItemAddAssetMutationFetcher(itemSlug, assetId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove an asset from an item.
 */
export const itemRemoveAsset = (itemSlug: string, assetId: string) => {
  return fetcher<ItemUpdateOKResponse>({
    url: `/v1/items/${itemSlug}/assets/${assetId}`,
    method: "DELETE",
  });
};

export const getItemRemoveAssetMutationFetcher = (
  itemSlug: string,
  assetId: string,
) => {
  return (_: string, __: { arg: Arguments }): Promise<ItemUpdateOKResponse> => {
    return itemRemoveAsset(itemSlug, assetId);
  };
};
export const getItemRemoveAssetMutationKey = (
  itemSlug: string,
  assetId: string,
) => `/v1/items/${itemSlug}/assets/${assetId}` as const;

export type ItemRemoveAssetMutationResult = NonNullable<
  Awaited<ReturnType<typeof itemRemoveAsset>>
>;
export type ItemRemoveAssetMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useItemRemoveAsset = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  itemSlug: string,
  assetId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof itemRemoveAsset>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof itemRemoveAsset>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getItemRemoveAssetMutationKey(itemSlug, assetId);
  const swrFn = getItemRemoveAssetMutationFetcher(itemSlug, assetId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
