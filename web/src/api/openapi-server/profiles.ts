/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import type {
  ProfileFollowersGetOKResponse,
  ProfileFollowersGetParams,
  ProfileFollowingGetOKResponse,
  ProfileFollowingGetParams,
  ProfileGetOKResponse,
  ProfileListOKResponse,
  ProfileListParams,
} from "../openapi-schema";
import { fetcher } from "../server";

/**
 * Query and search profiles.
 */
export type profileListResponse = {
  data: ProfileListOKResponse;
  status: number;
};

export const getProfileListUrl = (params?: ProfileListParams) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/profiles?${normalizedParams.toString()}`
    : `/profiles`;
};

export const profileList = async (
  params?: ProfileListParams,
  options?: RequestInit,
): Promise<profileListResponse> => {
  return fetcher<Promise<profileListResponse>>(getProfileListUrl(params), {
    ...options,
    method: "GET",
  });
};

/**
 * Get a public profile by ID.
 */
export type profileGetResponse = {
  data: ProfileGetOKResponse;
  status: number;
};

export const getProfileGetUrl = (accountHandle: string) => {
  return `/profiles/${accountHandle}`;
};

export const profileGet = async (
  accountHandle: string,
  options?: RequestInit,
): Promise<profileGetResponse> => {
  return fetcher<Promise<profileGetResponse>>(getProfileGetUrl(accountHandle), {
    ...options,
    method: "GET",
  });
};

/**
 * Get the followers and following details for a profile.
 */
export type profileFollowersGetResponse = {
  data: ProfileFollowersGetOKResponse;
  status: number;
};

export const getProfileFollowersGetUrl = (
  accountHandle: string,
  params?: ProfileFollowersGetParams,
) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/profiles/${accountHandle}/followers?${normalizedParams.toString()}`
    : `/profiles/${accountHandle}/followers`;
};

export const profileFollowersGet = async (
  accountHandle: string,
  params?: ProfileFollowersGetParams,
  options?: RequestInit,
): Promise<profileFollowersGetResponse> => {
  return fetcher<Promise<profileFollowersGetResponse>>(
    getProfileFollowersGetUrl(accountHandle, params),
    {
      ...options,
      method: "GET",
    },
  );
};

/**
 * Follow the specified profile as the authenticated account.
 */
export type profileFollowersAddResponse = {
  data: void;
  status: number;
};

export const getProfileFollowersAddUrl = (accountHandle: string) => {
  return `/profiles/${accountHandle}/followers`;
};

export const profileFollowersAdd = async (
  accountHandle: string,
  options?: RequestInit,
): Promise<profileFollowersAddResponse> => {
  return fetcher<Promise<profileFollowersAddResponse>>(
    getProfileFollowersAddUrl(accountHandle),
    {
      ...options,
      method: "PUT",
    },
  );
};

/**
 * Unfollow the specified profile as the authenticated account.
 */
export type profileFollowersRemoveResponse = {
  data: void;
  status: number;
};

export const getProfileFollowersRemoveUrl = (accountHandle: string) => {
  return `/profiles/${accountHandle}/followers`;
};

export const profileFollowersRemove = async (
  accountHandle: string,
  options?: RequestInit,
): Promise<profileFollowersRemoveResponse> => {
  return fetcher<Promise<profileFollowersRemoveResponse>>(
    getProfileFollowersRemoveUrl(accountHandle),
    {
      ...options,
      method: "DELETE",
    },
  );
};

/**
 * Get the profiles that this account is following.
 */
export type profileFollowingGetResponse = {
  data: ProfileFollowingGetOKResponse;
  status: number;
};

export const getProfileFollowingGetUrl = (
  accountHandle: string,
  params?: ProfileFollowingGetParams,
) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/profiles/${accountHandle}/following?${normalizedParams.toString()}`
    : `/profiles/${accountHandle}/following`;
};

export const profileFollowingGet = async (
  accountHandle: string,
  params?: ProfileFollowingGetParams,
  options?: RequestInit,
): Promise<profileFollowingGetResponse> => {
  return fetcher<Promise<profileFollowingGetResponse>>(
    getProfileFollowingGetUrl(accountHandle, params),
    {
      ...options,
      method: "GET",
    },
  );
};
