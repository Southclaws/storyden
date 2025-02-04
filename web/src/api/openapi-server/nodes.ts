/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type {
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
  NodeUpdateChildrenPropertySchemaBody,
  NodeUpdateChildrenPropertySchemaOKResponse,
  NodeUpdateOKResponse,
  NodeUpdateParams,
  NodeUpdatePropertiesBody,
  NodeUpdatePropertiesOKResponse,
  VisibilityUpdateBody,
} from "../openapi-schema";
import { fetcher } from "../server";

/**
 * Create a node for curating structured knowledge together.

 */
export type nodeCreateResponse = {
  data: NodeCreateOKResponse;
  status: number;
};

export const getNodeCreateUrl = () => {
  return `/nodes`;
};

export const nodeCreate = async (
  nodeCreateBody: NodeCreateBody,
  options?: RequestInit,
): Promise<nodeCreateResponse> => {
  return fetcher<Promise<nodeCreateResponse>>(getNodeCreateUrl(), {
    ...options,
    method: "POST",
    headers: { "Content-Type": "application/json", ...options?.headers },
    body: JSON.stringify(nodeCreateBody),
  });
};

/**
 * List nodes using the given filters. Can be used to get a full tree.

 */
export type nodeListResponse = {
  data: NodeListOKResponse;
  status: number;
};

export const getNodeListUrl = (params?: NodeListParams) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    const explodeParameters = ["visibility"];

    if (value instanceof Array && explodeParameters.includes(key)) {
      value.forEach((v) =>
        normalizedParams.append(key, v === null ? "null" : v.toString()),
      );
      return;
    }

    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/nodes?${normalizedParams.toString()}`
    : `/nodes`;
};

export const nodeList = async (
  params?: NodeListParams,
  options?: RequestInit,
): Promise<nodeListResponse> => {
  return fetcher<Promise<nodeListResponse>>(getNodeListUrl(params), {
    ...options,
    method: "GET",
  });
};

/**
 * Get a node by its URL slug.
 */
export type nodeGetResponse = {
  data: NodeGetOKResponse;
  status: number;
};

export const getNodeGetUrl = (nodeSlug: string) => {
  return `/nodes/${nodeSlug}`;
};

export const nodeGet = async (
  nodeSlug: string,
  options?: RequestInit,
): Promise<nodeGetResponse> => {
  return fetcher<Promise<nodeGetResponse>>(getNodeGetUrl(nodeSlug), {
    ...options,
    method: "GET",
  });
};

/**
 * Update a node.
 */
export type nodeUpdateResponse = {
  data: NodeUpdateOKResponse;
  status: number;
};

export const getNodeUpdateUrl = (
  nodeSlug: string,
  params?: NodeUpdateParams,
) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/nodes/${nodeSlug}?${normalizedParams.toString()}`
    : `/nodes/${nodeSlug}`;
};

export const nodeUpdate = async (
  nodeSlug: string,
  nodeUpdateBody: NodeUpdateBody,
  params?: NodeUpdateParams,
  options?: RequestInit,
): Promise<nodeUpdateResponse> => {
  return fetcher<Promise<nodeUpdateResponse>>(
    getNodeUpdateUrl(nodeSlug, params),
    {
      ...options,
      method: "PATCH",
      headers: { "Content-Type": "application/json", ...options?.headers },
      body: JSON.stringify(nodeUpdateBody),
    },
  );
};

/**
 * Delete a node and move all children to its parent or root.
 */
export type nodeDeleteResponse = {
  data: NodeDeleteOKResponse;
  status: number;
};

export const getNodeDeleteUrl = (
  nodeSlug: string,
  params?: NodeDeleteParams,
) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/nodes/${nodeSlug}?${normalizedParams.toString()}`
    : `/nodes/${nodeSlug}`;
};

export const nodeDelete = async (
  nodeSlug: string,
  params?: NodeDeleteParams,
  options?: RequestInit,
): Promise<nodeDeleteResponse> => {
  return fetcher<Promise<nodeDeleteResponse>>(
    getNodeDeleteUrl(nodeSlug, params),
    {
      ...options,
      method: "DELETE",
    },
  );
};

/**
 * Updates the property schema of the children of this node. All children
of a node use the same schema for properties resulting in a table-like
structure and behaviour. Property schemas are loosely structured and can
automatically cast their values sometimes. A failed cast will not change
data and instead just yield an empty value when reading however changing
the schema back to the original type (or a type compatible with what the
type was before changing) will retain the original data upon next read.
This permits clients to undo changes to the schema easily while allowing
quick schema changes without the need to remove or update values before.

 */
export type nodeUpdateChildrenPropertySchemaResponse = {
  data: NodeUpdateChildrenPropertySchemaOKResponse;
  status: number;
};

export const getNodeUpdateChildrenPropertySchemaUrl = (nodeSlug: string) => {
  return `/nodes/${nodeSlug}/children/property-schema`;
};

export const nodeUpdateChildrenPropertySchema = async (
  nodeSlug: string,
  nodeUpdateChildrenPropertySchemaBody: NodeUpdateChildrenPropertySchemaBody,
  options?: RequestInit,
): Promise<nodeUpdateChildrenPropertySchemaResponse> => {
  return fetcher<Promise<nodeUpdateChildrenPropertySchemaResponse>>(
    getNodeUpdateChildrenPropertySchemaUrl(nodeSlug),
    {
      ...options,
      method: "PATCH",
      headers: { "Content-Type": "application/json", ...options?.headers },
      body: JSON.stringify(nodeUpdateChildrenPropertySchemaBody),
    },
  );
};

/**
 * Update the properties of a node.
 */
export type nodeUpdatePropertiesResponse = {
  data: NodeUpdatePropertiesOKResponse;
  status: number;
};

export const getNodeUpdatePropertiesUrl = (nodeSlug: string) => {
  return `/nodes/${nodeSlug}/properties`;
};

export const nodeUpdateProperties = async (
  nodeSlug: string,
  nodeUpdatePropertiesBody: NodeUpdatePropertiesBody,
  options?: RequestInit,
): Promise<nodeUpdatePropertiesResponse> => {
  return fetcher<Promise<nodeUpdatePropertiesResponse>>(
    getNodeUpdatePropertiesUrl(nodeSlug),
    {
      ...options,
      method: "PATCH",
      headers: { "Content-Type": "application/json", ...options?.headers },
      body: JSON.stringify(nodeUpdatePropertiesBody),
    },
  );
};

/**
 * Update the visibility of a node. When changed, this may trigger other
operations such as notifications/newsletters. Changing the visibility of
anything to "published" is often accompanied by some other side effects.

 */
export type nodeUpdateVisibilityResponse = {
  data: NodeUpdateOKResponse;
  status: number;
};

export const getNodeUpdateVisibilityUrl = (nodeSlug: string) => {
  return `/nodes/${nodeSlug}/visibility`;
};

export const nodeUpdateVisibility = async (
  nodeSlug: string,
  visibilityUpdateBody: VisibilityUpdateBody,
  options?: RequestInit,
): Promise<nodeUpdateVisibilityResponse> => {
  return fetcher<Promise<nodeUpdateVisibilityResponse>>(
    getNodeUpdateVisibilityUrl(nodeSlug),
    {
      ...options,
      method: "PATCH",
      headers: { "Content-Type": "application/json", ...options?.headers },
      body: JSON.stringify(visibilityUpdateBody),
    },
  );
};

/**
 * Add an asset to a node.
 */
export type nodeAddAssetResponse = {
  data: NodeUpdateOKResponse;
  status: number;
};

export const getNodeAddAssetUrl = (
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/nodes/${nodeSlug}/assets/${assetId}?${normalizedParams.toString()}`
    : `/nodes/${nodeSlug}/assets/${assetId}`;
};

export const nodeAddAsset = async (
  nodeSlug: string,
  assetId: string,
  params?: NodeAddAssetParams,
  options?: RequestInit,
): Promise<nodeAddAssetResponse> => {
  return fetcher<Promise<nodeAddAssetResponse>>(
    getNodeAddAssetUrl(nodeSlug, assetId, params),
    {
      ...options,
      method: "PUT",
    },
  );
};

/**
 * Remove an asset from a node.
 */
export type nodeRemoveAssetResponse = {
  data: NodeUpdateOKResponse;
  status: number;
};

export const getNodeRemoveAssetUrl = (nodeSlug: string, assetId: string) => {
  return `/nodes/${nodeSlug}/assets/${assetId}`;
};

export const nodeRemoveAsset = async (
  nodeSlug: string,
  assetId: string,
  options?: RequestInit,
): Promise<nodeRemoveAssetResponse> => {
  return fetcher<Promise<nodeRemoveAssetResponse>>(
    getNodeRemoveAssetUrl(nodeSlug, assetId),
    {
      ...options,
      method: "DELETE",
    },
  );
};

/**
 * Set a node's parent to the specified node
 */
export type nodeAddNodeResponse = {
  data: NodeAddChildOKResponse;
  status: number;
};

export const getNodeAddNodeUrl = (nodeSlug: string, nodeSlugChild: string) => {
  return `/nodes/${nodeSlug}/nodes/${nodeSlugChild}`;
};

export const nodeAddNode = async (
  nodeSlug: string,
  nodeSlugChild: string,
  options?: RequestInit,
): Promise<nodeAddNodeResponse> => {
  return fetcher<Promise<nodeAddNodeResponse>>(
    getNodeAddNodeUrl(nodeSlug, nodeSlugChild),
    {
      ...options,
      method: "PUT",
    },
  );
};

/**
 * Remove a node from its parent node and back to the top level.

 */
export type nodeRemoveNodeResponse = {
  data: NodeRemoveChildOKResponse;
  status: number;
};

export const getNodeRemoveNodeUrl = (
  nodeSlug: string,
  nodeSlugChild: string,
) => {
  return `/nodes/${nodeSlug}/nodes/${nodeSlugChild}`;
};

export const nodeRemoveNode = async (
  nodeSlug: string,
  nodeSlugChild: string,
  options?: RequestInit,
): Promise<nodeRemoveNodeResponse> => {
  return fetcher<Promise<nodeRemoveNodeResponse>>(
    getNodeRemoveNodeUrl(nodeSlug, nodeSlugChild),
    {
      ...options,
      method: "DELETE",
    },
  );
};
