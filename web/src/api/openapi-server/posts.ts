/**
 * Generated by orval v6.31.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type {
  PostReactAddBody,
  PostReactAddOKResponse,
  PostSearchOKResponse,
  PostSearchParams,
  PostUpdateBody,
  PostUpdateOKResponse,
} from "../openapi-schema";
import { fetcher } from "../server";

/**
 * Publish changes to a single post.
 */
export type postUpdateResponse = {
  data: PostUpdateOKResponse;
  status: number;
};

export const getPostUpdateUrl = (postId: string) => {
  return `/posts/${postId}`;
};

export const postUpdate = async (
  postId: string,
  postUpdateBody: PostUpdateBody,
  options?: RequestInit,
): Promise<postUpdateResponse> => {
  return fetcher<Promise<postUpdateResponse>>(getPostUpdateUrl(postId), {
    ...options,
    method: "PATCH",
    body: JSON.stringify(postUpdateBody),
  });
};

/**
 * Archive a post using soft-delete.
 */
export type postDeleteResponse = {
  data: void;
  status: number;
};

export const getPostDeleteUrl = (postId: string) => {
  return `/posts/${postId}`;
};

export const postDelete = async (
  postId: string,
  options?: RequestInit,
): Promise<postDeleteResponse> => {
  return fetcher<Promise<postDeleteResponse>>(getPostDeleteUrl(postId), {
    ...options,
    method: "DELETE",
  });
};

/**
 * Search through posts using various queries and filters.
 */
export type postSearchResponse = {
  data: PostSearchOKResponse;
  status: number;
};

export const getPostSearchUrl = (params?: PostSearchParams) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value === null) {
      normalizedParams.append(key, "null");
    } else if (value !== undefined) {
      normalizedParams.append(key, value.toString());
    }
  });

  return `/posts/search?${normalizedParams.toString()}`;
};

export const postSearch = async (
  params?: PostSearchParams,
  options?: RequestInit,
): Promise<postSearchResponse> => {
  return fetcher<Promise<postSearchResponse>>(getPostSearchUrl(params), {
    ...options,
    method: "GET",
  });
};

/**
 * Add a reaction to a post.
 */
export type postReactAddResponse = {
  data: PostReactAddOKResponse;
  status: number;
};

export const getPostReactAddUrl = (postId: string) => {
  return `/posts/${postId}/reacts`;
};

export const postReactAdd = async (
  postId: string,
  postReactAddBody: PostReactAddBody,
  options?: RequestInit,
): Promise<postReactAddResponse> => {
  return fetcher<Promise<postReactAddResponse>>(getPostReactAddUrl(postId), {
    ...options,
    method: "PUT",
    body: JSON.stringify(postReactAddBody),
  });
};

/**
 * Remove a reaction from a post.
 */
export type postReactRemoveResponse = {
  data: void;
  status: number;
};

export const getPostReactRemoveUrl = (postId: string, reactId: string) => {
  return `/posts/${postId}/reacts/${reactId}`;
};

export const postReactRemove = async (
  postId: string,
  reactId: string,
  options?: RequestInit,
): Promise<postReactRemoveResponse> => {
  return fetcher<Promise<postReactRemoveResponse>>(
    getPostReactRemoveUrl(postId, reactId),
    {
      ...options,
      method: "DELETE",
    },
  );
};
