/**
 * Generated by orval v6.31.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type {
  AccountAuthProviderListOKResponse,
  AccountGetAvatarResponse,
  AccountGetOKResponse,
  AccountSetAvatarBody,
  AccountUpdateBody,
  AccountUpdateOKResponse,
} from "../openapi-schema";
import { fetcher } from "../server";

/**
 * Get the information for the currently authenticated account.
 */
export type accountGetResponse = {
  data: AccountGetOKResponse;
  status: number;
};

export const getAccountGetUrl = () => {
  return `/accounts`;
};

export const accountGet = async (
  options?: RequestInit,
): Promise<accountGetResponse> => {
  return fetcher<Promise<accountGetResponse>>(getAccountGetUrl(), {
    ...options,
    method: "GET",
  });
};

/**
 * Update the information for the currently authenticated account.
 */
export type accountUpdateResponse = {
  data: AccountUpdateOKResponse;
  status: number;
};

export const getAccountUpdateUrl = () => {
  return `/accounts`;
};

export const accountUpdate = async (
  accountUpdateBody: AccountUpdateBody,
  options?: RequestInit,
): Promise<accountUpdateResponse> => {
  return fetcher<Promise<accountUpdateResponse>>(getAccountUpdateUrl(), {
    ...options,
    method: "PATCH",
    body: JSON.stringify(accountUpdateBody),
  });
};

/**
 * Retrieve a list of authentication providers with a flag indicating which
ones are active for the currently authenticated account.

 */
export type accountAuthProviderListResponse = {
  data: AccountAuthProviderListOKResponse;
  status: number;
};

export const getAccountAuthProviderListUrl = () => {
  return `/accounts/self/auth-methods`;
};

export const accountAuthProviderList = async (
  options?: RequestInit,
): Promise<accountAuthProviderListResponse> => {
  return fetcher<Promise<accountAuthProviderListResponse>>(
    getAccountAuthProviderListUrl(),
    {
      ...options,
      method: "GET",
    },
  );
};

/**
 * Retrieve a list of authentication providers with a flag indicating which
ones are active for the currently authenticated account.

 */
export type accountAuthMethodDeleteResponse = {
  data: AccountAuthProviderListOKResponse;
  status: number;
};

export const getAccountAuthMethodDeleteUrl = (authMethodId: string) => {
  return `/accounts/self/auth-methods/${authMethodId}`;
};

export const accountAuthMethodDelete = async (
  authMethodId: string,
  options?: RequestInit,
): Promise<accountAuthMethodDeleteResponse> => {
  return fetcher<Promise<accountAuthMethodDeleteResponse>>(
    getAccountAuthMethodDeleteUrl(authMethodId),
    {
      ...options,
      method: "DELETE",
    },
  );
};

/**
 * Upload an avatar for the authenticated account.
 */
export type accountSetAvatarResponse = {
  data: void;
  status: number;
};

export const getAccountSetAvatarUrl = () => {
  return `/accounts/self/avatar`;
};

export const accountSetAvatar = async (
  accountSetAvatarBody: AccountSetAvatarBody,
  options?: RequestInit,
): Promise<accountSetAvatarResponse> => {
  return fetcher<Promise<accountSetAvatarResponse>>(getAccountSetAvatarUrl(), {
    ...options,
    method: "POST",
    body: JSON.stringify(accountSetAvatarBody),
  });
};

/**
 * Get an avatar for the specified account.
 */
export type accountGetAvatarResponse = {
  data: AccountGetAvatarResponse;
  status: number;
};

export const getAccountGetAvatarUrl = (accountHandle: string) => {
  return `/accounts/${accountHandle}/avatar`;
};

export const accountGetAvatar = async (
  accountHandle: string,
  options?: RequestInit,
): Promise<accountGetAvatarResponse> => {
  return fetcher<Promise<accountGetAvatarResponse>>(
    getAccountGetAvatarUrl(accountHandle),
    {
      ...options,
      method: "GET",
    },
  );
};

/**
 * Adds a role to an account. Members without the MANAGE_ROLES permission
cannot use this operation.

 */
export type accountAddRoleResponse = {
  data: AccountUpdateOKResponse;
  status: number;
};

export const getAccountAddRoleUrl = (accountHandle: string, roleId: string) => {
  return `/accounts/${accountHandle}/roles/${roleId}`;
};

export const accountAddRole = async (
  accountHandle: string,
  roleId: string,
  options?: RequestInit,
): Promise<accountAddRoleResponse> => {
  return fetcher<Promise<accountAddRoleResponse>>(
    getAccountAddRoleUrl(accountHandle, roleId),
    {
      ...options,
      method: "PUT",
    },
  );
};

/**
 * Removes a role from an account. Members without the MANAGE_ROLES cannot
use this operation. Admins cannot remove the admin role from themselves.

 */
export type accountRemoveRoleResponse = {
  data: AccountUpdateOKResponse;
  status: number;
};

export const getAccountRemoveRoleUrl = (
  accountHandle: string,
  roleId: string,
) => {
  return `/accounts/${accountHandle}/roles/${roleId}`;
};

export const accountRemoveRole = async (
  accountHandle: string,
  roleId: string,
  options?: RequestInit,
): Promise<accountRemoveRoleResponse> => {
  return fetcher<Promise<accountRemoveRoleResponse>>(
    getAccountRemoveRoleUrl(accountHandle, roleId),
    {
      ...options,
      method: "DELETE",
    },
  );
};
