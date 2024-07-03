/**
 * Generated by orval v6.30.2 🍺
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
  CollectionAddNodeOKResponse,
  CollectionAddPostOKResponse,
  CollectionCreateBody,
  CollectionCreateOKResponse,
  CollectionGetOKResponse,
  CollectionListOKResponse,
  CollectionRemoveNodeOKResponse,
  CollectionRemovePostOKResponse,
  CollectionUpdateBody,
  CollectionUpdateOKResponse,
  InternalServerErrorResponse,
  NotFoundResponse,
  UnauthorisedResponse,
} from "./schemas";

/**
 * Create a collection for curating posts under the authenticated account.

 */
export const collectionCreate = (
  collectionCreateBody: CollectionCreateBody,
) => {
  return fetcher<CollectionCreateOKResponse>({
    url: `/v1/collections`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: collectionCreateBody,
  });
};

export const getCollectionCreateMutationFetcher = () => {
  return (
    _: string,
    { arg }: { arg: CollectionCreateBody },
  ): Promise<CollectionCreateOKResponse> => {
    return collectionCreate(arg);
  };
};
export const getCollectionCreateMutationKey = () => `/v1/collections` as const;

export type CollectionCreateMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionCreate>>
>;
export type CollectionCreateMutationError =
  | UnauthorisedResponse
  | InternalServerErrorResponse;

export const useCollectionCreate = <
  TError = UnauthorisedResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof collectionCreate>>,
    TError,
    string,
    CollectionCreateBody,
    Awaited<ReturnType<typeof collectionCreate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getCollectionCreateMutationKey();
  const swrFn = getCollectionCreateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * List all collections using the filtering options.
 */
export const collectionList = () => {
  return fetcher<CollectionListOKResponse>({
    url: `/v1/collections`,
    method: "GET",
  });
};

export const getCollectionListKey = () => [`/v1/collections`] as const;

export type CollectionListQueryResult = NonNullable<
  Awaited<ReturnType<typeof collectionList>>
>;
export type CollectionListQueryError =
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionList = <
  TError = NotFoundResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRConfiguration<Awaited<ReturnType<typeof collectionList>>, TError> & {
    swrKey?: Key;
    enabled?: boolean;
  };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getCollectionListKey() : null));
  const swrFn = () => collectionList();

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
 * Get a collection by its ID. Collections can be public or private so the
response will depend on which account is making the request and if the
target collection is public, private, owned or not owned by the account.

 */
export const collectionGet = (collectionId: string) => {
  return fetcher<CollectionGetOKResponse>({
    url: `/v1/collections/${collectionId}`,
    method: "GET",
  });
};

export const getCollectionGetKey = (collectionId: string) =>
  [`/v1/collections/${collectionId}`] as const;

export type CollectionGetQueryResult = NonNullable<
  Awaited<ReturnType<typeof collectionGet>>
>;
export type CollectionGetQueryError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionGet = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  options?: {
    swr?: SWRConfiguration<
      Awaited<ReturnType<typeof collectionGet>>,
      TError
    > & { swrKey?: Key; enabled?: boolean };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!collectionId;
  const swrKey =
    swrOptions?.swrKey ??
    (() => (isEnabled ? getCollectionGetKey(collectionId) : null));
  const swrFn = () => collectionGet(collectionId);

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
 * Update a collection owned by the authenticated account.
 */
export const collectionUpdate = (
  collectionId: string,
  collectionUpdateBody: CollectionUpdateBody,
) => {
  return fetcher<CollectionUpdateOKResponse>({
    url: `/v1/collections/${collectionId}`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: collectionUpdateBody,
  });
};

export const getCollectionUpdateMutationFetcher = (collectionId: string) => {
  return (
    _: string,
    { arg }: { arg: CollectionUpdateBody },
  ): Promise<CollectionUpdateOKResponse> => {
    return collectionUpdate(collectionId, arg);
  };
};
export const getCollectionUpdateMutationKey = (collectionId: string) =>
  `/v1/collections/${collectionId}` as const;

export type CollectionUpdateMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionUpdate>>
>;
export type CollectionUpdateMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionUpdate = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionUpdate>>,
      TError,
      string,
      CollectionUpdateBody,
      Awaited<ReturnType<typeof collectionUpdate>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getCollectionUpdateMutationKey(collectionId);
  const swrFn = getCollectionUpdateMutationFetcher(collectionId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Delete a collection owned by the authenticated account.
 */
export const collectionDelete = (collectionId: string) => {
  return fetcher<void>({
    url: `/v1/collections/${collectionId}`,
    method: "DELETE",
  });
};

export const getCollectionDeleteMutationFetcher = (collectionId: string) => {
  return (_: string, __: { arg: Arguments }): Promise<void> => {
    return collectionDelete(collectionId);
  };
};
export const getCollectionDeleteMutationKey = (collectionId: string) =>
  `/v1/collections/${collectionId}` as const;

export type CollectionDeleteMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionDelete>>
>;
export type CollectionDeleteMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionDelete = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionDelete>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof collectionDelete>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getCollectionDeleteMutationKey(collectionId);
  const swrFn = getCollectionDeleteMutationFetcher(collectionId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Add a post to a collection. The collection must be owned by the account
making the request. The post can be any published post of any kind.

 */
export const collectionAddPost = (collectionId: string, postId: string) => {
  return fetcher<CollectionAddPostOKResponse>({
    url: `/v1/collections/${collectionId}/posts/${postId}`,
    method: "PUT",
  });
};

export const getCollectionAddPostMutationFetcher = (
  collectionId: string,
  postId: string,
) => {
  return (
    _: string,
    __: { arg: Arguments },
  ): Promise<CollectionAddPostOKResponse> => {
    return collectionAddPost(collectionId, postId);
  };
};
export const getCollectionAddPostMutationKey = (
  collectionId: string,
  postId: string,
) => `/v1/collections/${collectionId}/posts/${postId}` as const;

export type CollectionAddPostMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionAddPost>>
>;
export type CollectionAddPostMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionAddPost = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  postId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionAddPost>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof collectionAddPost>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getCollectionAddPostMutationKey(collectionId, postId);
  const swrFn = getCollectionAddPostMutationFetcher(collectionId, postId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove a post from a collection. The collection must be owned by the
account making the request.

 */
export const collectionRemovePost = (collectionId: string, postId: string) => {
  return fetcher<CollectionRemovePostOKResponse>({
    url: `/v1/collections/${collectionId}/posts/${postId}`,
    method: "DELETE",
  });
};

export const getCollectionRemovePostMutationFetcher = (
  collectionId: string,
  postId: string,
) => {
  return (
    _: string,
    __: { arg: Arguments },
  ): Promise<CollectionRemovePostOKResponse> => {
    return collectionRemovePost(collectionId, postId);
  };
};
export const getCollectionRemovePostMutationKey = (
  collectionId: string,
  postId: string,
) => `/v1/collections/${collectionId}/posts/${postId}` as const;

export type CollectionRemovePostMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionRemovePost>>
>;
export type CollectionRemovePostMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionRemovePost = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  postId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionRemovePost>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof collectionRemovePost>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ??
    getCollectionRemovePostMutationKey(collectionId, postId);
  const swrFn = getCollectionRemovePostMutationFetcher(collectionId, postId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Add a node to a collection. The collection must be owned by the account
making the request. The node can be any published node or any node
not published but owned by the collection owner.

 */
export const collectionAddNode = (collectionId: string, nodeId: string) => {
  return fetcher<CollectionAddNodeOKResponse>({
    url: `/v1/collections/${collectionId}/nodes/${nodeId}`,
    method: "PUT",
  });
};

export const getCollectionAddNodeMutationFetcher = (
  collectionId: string,
  nodeId: string,
) => {
  return (
    _: string,
    __: { arg: Arguments },
  ): Promise<CollectionAddNodeOKResponse> => {
    return collectionAddNode(collectionId, nodeId);
  };
};
export const getCollectionAddNodeMutationKey = (
  collectionId: string,
  nodeId: string,
) => `/v1/collections/${collectionId}/nodes/${nodeId}` as const;

export type CollectionAddNodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionAddNode>>
>;
export type CollectionAddNodeMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionAddNode = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  nodeId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionAddNode>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof collectionAddNode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getCollectionAddNodeMutationKey(collectionId, nodeId);
  const swrFn = getCollectionAddNodeMutationFetcher(collectionId, nodeId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove a node from a collection. The collection must be owned by the
account making the request.

 */
export const collectionRemoveNode = (collectionId: string, nodeId: string) => {
  return fetcher<CollectionRemoveNodeOKResponse>({
    url: `/v1/collections/${collectionId}/nodes/${nodeId}`,
    method: "DELETE",
  });
};

export const getCollectionRemoveNodeMutationFetcher = (
  collectionId: string,
  nodeId: string,
) => {
  return (
    _: string,
    __: { arg: Arguments },
  ): Promise<CollectionRemoveNodeOKResponse> => {
    return collectionRemoveNode(collectionId, nodeId);
  };
};
export const getCollectionRemoveNodeMutationKey = (
  collectionId: string,
  nodeId: string,
) => `/v1/collections/${collectionId}/nodes/${nodeId}` as const;

export type CollectionRemoveNodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof collectionRemoveNode>>
>;
export type CollectionRemoveNodeMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useCollectionRemoveNode = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  collectionId: string,
  nodeId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof collectionRemoveNode>>,
      TError,
      string,
      Arguments,
      Awaited<ReturnType<typeof collectionRemoveNode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ??
    getCollectionRemoveNodeMutationKey(collectionId, nodeId);
  const swrFn = getCollectionRemoveNodeMutationFetcher(collectionId, nodeId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
