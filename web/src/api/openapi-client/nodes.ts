/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import useSwr from "swr";
import type { Arguments, Key, SWRConfiguration } from "swr";
import useSWRMutation from "swr/mutation";
import type { SWRMutationConfiguration } from "swr/mutation";

import { fetcher } from "../client";
import type {
  InternalServerErrorResponse,
  NodeAddAssetParams,
  NodeAddChildOKResponse,
  NodeCreateBody,
  NodeCreateOKResponse,
  NodeDeleteOKResponse,
  NodeDeleteParams,
  NodeGetOKResponse,
  NodeListOKResponse,
  NodeListParams,
  NodeRemoveChildOKResponse,
  NodeUpdateBody,
  NodeUpdateOKResponse,
  NodeUpdateParams,
  NotFoundResponse,
  UnauthorisedResponse,
  VisibilityUpdateBody,
} from "../openapi-schema";

/**
 * Create a node for curating structured knowledge together.

 */
export const nodeCreate = (nodeCreateBody: NodeCreateBody) => {
  return fetcher<NodeCreateOKResponse>({
    url: `/nodes`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: nodeCreateBody,
  });
};

export const getNodeCreateMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: NodeCreateBody },
  ): Promise<NodeCreateOKResponse> => {
    return nodeCreate(arg);
  };
};
export const getNodeCreateMutationKey = () => [`/nodes`] as const;

export type NodeCreateMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeCreate>>
>;
export type NodeCreateMutationError =
  | UnauthorisedResponse
  | InternalServerErrorResponse;

export const useNodeCreate = <
  TError = UnauthorisedResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof nodeCreate>>,
    TError,
    Key,
    NodeCreateBody,
    Awaited<ReturnType<typeof nodeCreate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getNodeCreateMutationKey();
  const swrFn = getNodeCreateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * List nodes using the given filters. Can be used to get a full tree.

 */
export const nodeList = (params?: NodeListParams) => {
  return fetcher<NodeListOKResponse>({ url: `/nodes`, method: "GET", params });
};

export const getNodeListKey = (params?: NodeListParams) =>
  [`/nodes`, ...(params ? [params] : [])] as const;

export type NodeListQueryResult = NonNullable<
  Awaited<ReturnType<typeof nodeList>>
>;
export type NodeListQueryError = NotFoundResponse | InternalServerErrorResponse;

export const useNodeList = <
  TError = NotFoundResponse | InternalServerErrorResponse,
>(
  params?: NodeListParams,
  options?: {
    swr?: SWRConfiguration<Awaited<ReturnType<typeof nodeList>>, TError> & {
      swrKey?: Key;
      enabled?: boolean;
    };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getNodeListKey(params) : null));
  const swrFn = () => nodeList(params);

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
 * Get a node by its URL slug.
 */
export const nodeGet = (nodeSlug: string) => {
  return fetcher<NodeGetOKResponse>({
    url: `/nodes/${nodeSlug}`,
    method: "GET",
  });
};

export const getNodeGetKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}`] as const;

export type NodeGetQueryResult = NonNullable<
  Awaited<ReturnType<typeof nodeGet>>
>;
export type NodeGetQueryError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeGet = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRConfiguration<Awaited<ReturnType<typeof nodeGet>>, TError> & {
      swrKey?: Key;
      enabled?: boolean;
    };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!nodeSlug;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getNodeGetKey(nodeSlug) : null));
  const swrFn = () => nodeGet(nodeSlug);

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
 * Update a node.
 */
export const nodeUpdate = (
  nodeSlug: string,
  nodeUpdateBody: NodeUpdateBody,
  params?: NodeUpdateParams,
) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdateBody,
    params,
  });
};

export const getNodeUpdateMutationFetcher = (
  nodeSlug: string,
  params?: NodeUpdateParams,
) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdateBody },
  ): Promise<NodeUpdateOKResponse> => {
    return nodeUpdate(nodeSlug, arg, params);
  };
};
export const getNodeUpdateMutationKey = (
  nodeSlug: string,
  params?: NodeUpdateParams,
) => [`/nodes/${nodeSlug}`, ...(params ? [params] : [])] as const;

export type NodeUpdateMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdate>>
>;
export type NodeUpdateMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdate = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  params?: NodeUpdateParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdate>>,
      TError,
      Key,
      NodeUpdateBody,
      Awaited<ReturnType<typeof nodeUpdate>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeUpdateMutationKey(nodeSlug, params);
  const swrFn = getNodeUpdateMutationFetcher(nodeSlug, params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Delete a node and move all children to its parent or root.
 */
export const nodeDelete = (nodeSlug: string, params?: NodeDeleteParams) => {
  return fetcher<NodeDeleteOKResponse>({
    url: `/nodes/${nodeSlug}`,
    method: "DELETE",
    params,
  });
};

export const getNodeDeleteMutationFetcher = (
  nodeSlug: string,
  params?: NodeDeleteParams,
) => {
  return (_: Key, __: { arg: Arguments }): Promise<NodeDeleteOKResponse> => {
    return nodeDelete(nodeSlug, params);
  };
};
export const getNodeDeleteMutationKey = (
  nodeSlug: string,
  params?: NodeDeleteParams,
) => [`/nodes/${nodeSlug}`, ...(params ? [params] : [])] as const;

export type NodeDeleteMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeDelete>>
>;
export type NodeDeleteMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeDelete = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  params?: NodeDeleteParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeDelete>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof nodeDelete>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeDeleteMutationKey(nodeSlug, params);
  const swrFn = getNodeDeleteMutationFetcher(nodeSlug, params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Update the visibility of a node. When changed, this may trigger other
operations such as notifications/newsletters. Changing the visibility of
anything to "published" is often accompanied by some other side effects.

 */
export const nodeUpdateVisibility = (
  nodeSlug: string,
  visibilityUpdateBody: VisibilityUpdateBody,
) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}/visibility`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: visibilityUpdateBody,
  });
};

export const getNodeUpdateVisibilityMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: VisibilityUpdateBody },
  ): Promise<NodeUpdateOKResponse> => {
    return nodeUpdateVisibility(nodeSlug, arg);
  };
};
export const getNodeUpdateVisibilityMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/visibility`] as const;

export type NodeUpdateVisibilityMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdateVisibility>>
>;
export type NodeUpdateVisibilityMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdateVisibility = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdateVisibility>>,
      TError,
      Key,
      VisibilityUpdateBody,
      Awaited<ReturnType<typeof nodeUpdateVisibility>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeUpdateVisibilityMutationKey(nodeSlug);
  const swrFn = getNodeUpdateVisibilityMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Add an asset to a node.
 */
export const nodeAddAsset = (
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}/assets/${assetId}`,
    method: "PUT",
    params,
  });
};

export const getNodeAddAssetMutationFetcher = (
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
) => {
  return (_: Key, __: { arg: Arguments }): Promise<NodeUpdateOKResponse> => {
    return nodeAddAsset(nodeSlug, assetId, params);
  };
};
export const getNodeAddAssetMutationKey = (
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
) =>
  [
    `/nodes/${nodeSlug}/assets/${assetId}`,
    ...(params ? [params] : []),
  ] as const;

export type NodeAddAssetMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeAddAsset>>
>;
export type NodeAddAssetMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeAddAsset = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeAddAsset>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof nodeAddAsset>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeAddAssetMutationKey(nodeSlug, assetId, params);
  const swrFn = getNodeAddAssetMutationFetcher(nodeSlug, assetId, params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove an asset from a node.
 */
export const nodeRemoveAsset = (nodeSlug: string, assetId: string) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}/assets/${assetId}`,
    method: "DELETE",
  });
};

export const getNodeRemoveAssetMutationFetcher = (
  nodeSlug: string,
  assetId: string,
) => {
  return (_: Key, __: { arg: Arguments }): Promise<NodeUpdateOKResponse> => {
    return nodeRemoveAsset(nodeSlug, assetId);
  };
};
export const getNodeRemoveAssetMutationKey = (
  nodeSlug: string,
  assetId: string,
) => [`/nodes/${nodeSlug}/assets/${assetId}`] as const;

export type NodeRemoveAssetMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeRemoveAsset>>
>;
export type NodeRemoveAssetMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeRemoveAsset = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  assetId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeRemoveAsset>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof nodeRemoveAsset>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeRemoveAssetMutationKey(nodeSlug, assetId);
  const swrFn = getNodeRemoveAssetMutationFetcher(nodeSlug, assetId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Set a node's parent to the specified node
 */
export const nodeAddNode = (nodeSlug: string, nodeSlugChild: string) => {
  return fetcher<NodeAddChildOKResponse>({
    url: `/nodes/${nodeSlug}/nodes/${nodeSlugChild}`,
    method: "PUT",
  });
};

export const getNodeAddNodeMutationFetcher = (
  nodeSlug: string,
  nodeSlugChild: string,
) => {
  return (_: Key, __: { arg: Arguments }): Promise<NodeAddChildOKResponse> => {
    return nodeAddNode(nodeSlug, nodeSlugChild);
  };
};
export const getNodeAddNodeMutationKey = (
  nodeSlug: string,
  nodeSlugChild: string,
) => [`/nodes/${nodeSlug}/nodes/${nodeSlugChild}`] as const;

export type NodeAddNodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeAddNode>>
>;
export type NodeAddNodeMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeAddNode = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  nodeSlugChild: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeAddNode>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof nodeAddNode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeAddNodeMutationKey(nodeSlug, nodeSlugChild);
  const swrFn = getNodeAddNodeMutationFetcher(nodeSlug, nodeSlugChild);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove a node from its parent node and back to the top level.

 */
export const nodeRemoveNode = (nodeSlug: string, nodeSlugChild: string) => {
  return fetcher<NodeRemoveChildOKResponse>({
    url: `/nodes/${nodeSlug}/nodes/${nodeSlugChild}`,
    method: "DELETE",
  });
};

export const getNodeRemoveNodeMutationFetcher = (
  nodeSlug: string,
  nodeSlugChild: string,
) => {
  return (
    _: Key,
    __: { arg: Arguments },
  ): Promise<NodeRemoveChildOKResponse> => {
    return nodeRemoveNode(nodeSlug, nodeSlugChild);
  };
};
export const getNodeRemoveNodeMutationKey = (
  nodeSlug: string,
  nodeSlugChild: string,
) => [`/nodes/${nodeSlug}/nodes/${nodeSlugChild}`] as const;

export type NodeRemoveNodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeRemoveNode>>
>;
export type NodeRemoveNodeMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeRemoveNode = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  nodeSlugChild: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeRemoveNode>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof nodeRemoveNode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeRemoveNodeMutationKey(nodeSlug, nodeSlugChild);
  const swrFn = getNodeRemoveNodeMutationFetcher(nodeSlug, nodeSlugChild);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
