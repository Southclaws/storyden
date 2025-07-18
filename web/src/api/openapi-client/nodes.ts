/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import useSwr from "swr";
import type { Arguments, Key, SWRConfiguration } from "swr";
import useSWRMutation from "swr/mutation";
import type { SWRMutationConfiguration } from "swr/mutation";

import { fetcher } from "../client";
import type {
  BadRequestResponse,
  InternalServerErrorResponse,
  NodeAddChildOKResponse,
  NodeCreateBody,
  NodeCreateOKResponse,
  NodeDeleteOKResponse,
  NodeDeleteParams,
  NodeGenerateContentBody,
  NodeGenerateContentOKResponse,
  NodeGenerateTagsBody,
  NodeGenerateTagsOKResponse,
  NodeGenerateTitleBody,
  NodeGenerateTitleOKResponse,
  NodeGetOKResponse,
  NodeGetParams,
  NodeListChildrenParams,
  NodeListOKResponse,
  NodeListParams,
  NodeRemoveChildOKResponse,
  NodeUpdateBody,
  NodeUpdateOKResponse,
  NodeUpdatePositionBody,
  NodeUpdatePropertiesBody,
  NodeUpdatePropertiesOKResponse,
  NodeUpdatePropertySchemaBody,
  NodeUpdatePropertySchemaOKResponse,
  NotFoundResponse,
  NotImplementedResponse,
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
export const nodeGet = (nodeSlug: string, params?: NodeGetParams) => {
  return fetcher<NodeGetOKResponse>({
    url: `/nodes/${nodeSlug}`,
    method: "GET",
    params,
  });
};

export const getNodeGetKey = (nodeSlug: string, params?: NodeGetParams) =>
  [`/nodes/${nodeSlug}`, ...(params ? [params] : [])] as const;

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
  params?: NodeGetParams,
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
    swrOptions?.swrKey ??
    (() => (isEnabled ? getNodeGetKey(nodeSlug, params) : null));
  const swrFn = () => nodeGet(nodeSlug, params);

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
) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdateBody,
  });
};

export const getNodeUpdateMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdateBody },
  ): Promise<NodeUpdateOKResponse> => {
    return nodeUpdate(nodeSlug, arg);
  };
};
export const getNodeUpdateMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}`] as const;

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

  const swrKey = swrOptions?.swrKey ?? getNodeUpdateMutationKey(nodeSlug);
  const swrFn = getNodeUpdateMutationFetcher(nodeSlug);

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
 * Generate a proposed title for the specified node. Will not actually
mutate the specified node, instead will return a proposal based on the
output from a language model call.

 */
export const nodeGenerateTitle = (
  nodeSlug: string,
  nodeGenerateTitleBody: NodeGenerateTitleBody,
) => {
  return fetcher<NodeGenerateTitleOKResponse>({
    url: `/nodes/${nodeSlug}/title`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: nodeGenerateTitleBody,
  });
};

export const getNodeGenerateTitleMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeGenerateTitleBody },
  ): Promise<NodeGenerateTitleOKResponse> => {
    return nodeGenerateTitle(nodeSlug, arg);
  };
};
export const getNodeGenerateTitleMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/title`] as const;

export type NodeGenerateTitleMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeGenerateTitle>>
>;
export type NodeGenerateTitleMutationError =
  | BadRequestResponse
  | UnauthorisedResponse
  | NotFoundResponse
  | NotImplementedResponse
  | InternalServerErrorResponse;

export const useNodeGenerateTitle = <
  TError =
    | BadRequestResponse
    | UnauthorisedResponse
    | NotFoundResponse
    | NotImplementedResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeGenerateTitle>>,
      TError,
      Key,
      NodeGenerateTitleBody,
      Awaited<ReturnType<typeof nodeGenerateTitle>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeGenerateTitleMutationKey(nodeSlug);
  const swrFn = getNodeGenerateTitleMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Generate proposed tags for the specified node. Will not actually mutate
the specified node, instead will return a proposal based on the output
from a language model call.

 */
export const nodeGenerateTags = (
  nodeSlug: string,
  nodeGenerateTagsBody: NodeGenerateTagsBody,
) => {
  return fetcher<NodeGenerateTagsOKResponse>({
    url: `/nodes/${nodeSlug}/tags`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: nodeGenerateTagsBody,
  });
};

export const getNodeGenerateTagsMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeGenerateTagsBody },
  ): Promise<NodeGenerateTagsOKResponse> => {
    return nodeGenerateTags(nodeSlug, arg);
  };
};
export const getNodeGenerateTagsMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/tags`] as const;

export type NodeGenerateTagsMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeGenerateTags>>
>;
export type NodeGenerateTagsMutationError =
  | BadRequestResponse
  | UnauthorisedResponse
  | NotFoundResponse
  | NotImplementedResponse
  | InternalServerErrorResponse;

export const useNodeGenerateTags = <
  TError =
    | BadRequestResponse
    | UnauthorisedResponse
    | NotFoundResponse
    | NotImplementedResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeGenerateTags>>,
      TError,
      Key,
      NodeGenerateTagsBody,
      Awaited<ReturnType<typeof nodeGenerateTags>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getNodeGenerateTagsMutationKey(nodeSlug);
  const swrFn = getNodeGenerateTagsMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Generate proposed content for the specified node. Will not actually
mutate the specified node, instead will return a proposal based on the
output from a language model call.

 */
export const nodeGenerateContent = (
  nodeSlug: string,
  nodeGenerateContentBody: NodeGenerateContentBody,
) => {
  return fetcher<NodeGenerateContentOKResponse>({
    url: `/nodes/${nodeSlug}/content`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: nodeGenerateContentBody,
  });
};

export const getNodeGenerateContentMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeGenerateContentBody },
  ): Promise<NodeGenerateContentOKResponse> => {
    return nodeGenerateContent(nodeSlug, arg);
  };
};
export const getNodeGenerateContentMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/content`] as const;

export type NodeGenerateContentMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeGenerateContent>>
>;
export type NodeGenerateContentMutationError =
  | BadRequestResponse
  | UnauthorisedResponse
  | NotFoundResponse
  | NotImplementedResponse
  | InternalServerErrorResponse;

export const useNodeGenerateContent = <
  TError =
    | BadRequestResponse
    | UnauthorisedResponse
    | NotFoundResponse
    | NotImplementedResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeGenerateContent>>,
      TError,
      Key,
      NodeGenerateContentBody,
      Awaited<ReturnType<typeof nodeGenerateContent>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeGenerateContentMutationKey(nodeSlug);
  const swrFn = getNodeGenerateContentMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Get all the children of a given node using the provided filters and page
parameters. This can be used for rendering the child nodes of the given
node as an interactive table where properties can be used as columns.

 */
export const nodeListChildren = (
  nodeSlug: string,
  params?: NodeListChildrenParams,
) => {
  return fetcher<NodeListOKResponse>({
    url: `/nodes/${nodeSlug}/children`,
    method: "GET",
    params,
  });
};

export const getNodeListChildrenKey = (
  nodeSlug: string,
  params?: NodeListChildrenParams,
) => [`/nodes/${nodeSlug}/children`, ...(params ? [params] : [])] as const;

export type NodeListChildrenQueryResult = NonNullable<
  Awaited<ReturnType<typeof nodeListChildren>>
>;
export type NodeListChildrenQueryError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeListChildren = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  params?: NodeListChildrenParams,
  options?: {
    swr?: SWRConfiguration<
      Awaited<ReturnType<typeof nodeListChildren>>,
      TError
    > & { swrKey?: Key; enabled?: boolean };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!nodeSlug;
  const swrKey =
    swrOptions?.swrKey ??
    (() => (isEnabled ? getNodeListChildrenKey(nodeSlug, params) : null));
  const swrFn = () => nodeListChildren(nodeSlug, params);

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
 * Updates the property schema of the children of this node. All children
of a node use the same schema for properties resulting in a table-like
structure and behaviour. See also: NodeUpdatePropertySchema

 */
export const nodeUpdateChildrenPropertySchema = (
  nodeSlug: string,
  nodeUpdatePropertySchemaBody: NodeUpdatePropertySchemaBody,
) => {
  return fetcher<NodeUpdatePropertySchemaOKResponse>({
    url: `/nodes/${nodeSlug}/children/property-schema`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdatePropertySchemaBody,
  });
};

export const getNodeUpdateChildrenPropertySchemaMutationFetcher = (
  nodeSlug: string,
) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdatePropertySchemaBody },
  ): Promise<NodeUpdatePropertySchemaOKResponse> => {
    return nodeUpdateChildrenPropertySchema(nodeSlug, arg);
  };
};
export const getNodeUpdateChildrenPropertySchemaMutationKey = (
  nodeSlug: string,
) => [`/nodes/${nodeSlug}/children/property-schema`] as const;

export type NodeUpdateChildrenPropertySchemaMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdateChildrenPropertySchema>>
>;
export type NodeUpdateChildrenPropertySchemaMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdateChildrenPropertySchema = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdateChildrenPropertySchema>>,
      TError,
      Key,
      NodeUpdatePropertySchemaBody,
      Awaited<ReturnType<typeof nodeUpdateChildrenPropertySchema>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ??
    getNodeUpdateChildrenPropertySchemaMutationKey(nodeSlug);
  const swrFn = getNodeUpdateChildrenPropertySchemaMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Updates the property schema of this node and its siblings. All children
of a node use the same schema for properties resulting in a table-like
structure and behaviour. Property schemas are loosely structured and can
automatically cast their values sometimes. A failed cast will not change
data and instead just yield an empty value when reading however changing
the schema back to the original type (or a type compatible with what the
type was before changing) will retain the original data upon next read.
This permits clients to undo changes to the schema easily while allowing
quick schema changes without the need to remove or update values before.

 */
export const nodeUpdatePropertySchema = (
  nodeSlug: string,
  nodeUpdatePropertySchemaBody: NodeUpdatePropertySchemaBody,
) => {
  return fetcher<NodeUpdatePropertySchemaOKResponse>({
    url: `/nodes/${nodeSlug}/property-schema`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdatePropertySchemaBody,
  });
};

export const getNodeUpdatePropertySchemaMutationFetcher = (
  nodeSlug: string,
) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdatePropertySchemaBody },
  ): Promise<NodeUpdatePropertySchemaOKResponse> => {
    return nodeUpdatePropertySchema(nodeSlug, arg);
  };
};
export const getNodeUpdatePropertySchemaMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/property-schema`] as const;

export type NodeUpdatePropertySchemaMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdatePropertySchema>>
>;
export type NodeUpdatePropertySchemaMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdatePropertySchema = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdatePropertySchema>>,
      TError,
      Key,
      NodeUpdatePropertySchemaBody,
      Awaited<ReturnType<typeof nodeUpdatePropertySchema>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeUpdatePropertySchemaMutationKey(nodeSlug);
  const swrFn = getNodeUpdatePropertySchemaMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Update the properties of a node. New schema fields will result in the
schema of the node being updated before values are assigned. This will
also propagate to all sibling nodes as they all share the same schema.

 */
export const nodeUpdateProperties = (
  nodeSlug: string,
  nodeUpdatePropertiesBody: NodeUpdatePropertiesBody,
) => {
  return fetcher<NodeUpdatePropertiesOKResponse>({
    url: `/nodes/${nodeSlug}/properties`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdatePropertiesBody,
  });
};

export const getNodeUpdatePropertiesMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdatePropertiesBody },
  ): Promise<NodeUpdatePropertiesOKResponse> => {
    return nodeUpdateProperties(nodeSlug, arg);
  };
};
export const getNodeUpdatePropertiesMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/properties`] as const;

export type NodeUpdatePropertiesMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdateProperties>>
>;
export type NodeUpdatePropertiesMutationError =
  | BadRequestResponse
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdateProperties = <
  TError =
    | BadRequestResponse
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdateProperties>>,
      TError,
      Key,
      NodeUpdatePropertiesBody,
      Awaited<ReturnType<typeof nodeUpdateProperties>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeUpdatePropertiesMutationKey(nodeSlug);
  const swrFn = getNodeUpdatePropertiesMutationFetcher(nodeSlug);

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
export const nodeAddAsset = (nodeSlug: string, assetId: string) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}/assets/${assetId}`,
    method: "PUT",
  });
};

export const getNodeAddAssetMutationFetcher = (
  nodeSlug: string,
  assetId: string,
) => {
  return (_: Key, __: { arg: Arguments }): Promise<NodeUpdateOKResponse> => {
    return nodeAddAsset(nodeSlug, assetId);
  };
};
export const getNodeAddAssetMutationKey = (nodeSlug: string, assetId: string) =>
  [`/nodes/${nodeSlug}/assets/${assetId}`] as const;

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
    swrOptions?.swrKey ?? getNodeAddAssetMutationKey(nodeSlug, assetId);
  const swrFn = getNodeAddAssetMutationFetcher(nodeSlug, assetId);

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
/**
 * Update the node's position in the tree, which optionally allows for 
changing the node's parent either to another node or to `null` which
severs the parent and moves the node to the root. This endpoint also
allows for moving the node's sort position within either its current
parent, or when moving it to a new parent. Use this operation for a
draggable tree interface or a table interface.

 */
export const nodeUpdatePosition = (
  nodeSlug: string,
  nodeUpdatePositionBody: NodeUpdatePositionBody,
) => {
  return fetcher<NodeUpdateOKResponse>({
    url: `/nodes/${nodeSlug}/position`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: nodeUpdatePositionBody,
  });
};

export const getNodeUpdatePositionMutationFetcher = (nodeSlug: string) => {
  return (
    _: Key,
    { arg }: { arg: NodeUpdatePositionBody },
  ): Promise<NodeUpdateOKResponse> => {
    return nodeUpdatePosition(nodeSlug, arg);
  };
};
export const getNodeUpdatePositionMutationKey = (nodeSlug: string) =>
  [`/nodes/${nodeSlug}/position`] as const;

export type NodeUpdatePositionMutationResult = NonNullable<
  Awaited<ReturnType<typeof nodeUpdatePosition>>
>;
export type NodeUpdatePositionMutationError =
  | BadRequestResponse
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useNodeUpdatePosition = <
  TError =
    | BadRequestResponse
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  nodeSlug: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof nodeUpdatePosition>>,
      TError,
      Key,
      NodeUpdatePositionBody,
      Awaited<ReturnType<typeof nodeUpdatePosition>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getNodeUpdatePositionMutationKey(nodeSlug);
  const swrFn = getNodeUpdatePositionMutationFetcher(nodeSlug);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
