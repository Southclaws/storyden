import { buildRequest, buildResult } from "../src/api/common";
import { getAccountAddRoleMutationKey } from "../src/api/openapi-client/accounts";
import { getCategoryCreateMutationKey } from "../src/api/openapi-client/categories";
import { getThreadCreateMutationKey } from "../src/api/openapi-client/threads";
import {
  AccountUpdateOKResponse,
  CategoryCreateBody,
  CategoryCreateOKResponse,
  ThreadCreateBody,
  ThreadCreateOKResponse,
} from "../src/api/openapi-schema";

export type AccessKeyClient = {
  accountAddRole: (
    accountHandle: string,
    roleId: string,
  ) => Promise<AccountUpdateOKResponse>;
  categoryCreate: (
    categoryCreateBody: CategoryCreateBody,
  ) => Promise<CategoryCreateOKResponse>;
  threadCreate: (
    threadCreateBody: ThreadCreateBody,
  ) => Promise<ThreadCreateOKResponse>;
};

export function createAccessKeyClient(accessKey: string): AccessKeyClient {
  return {
    accountAddRole: async (accountHandle, roleId) => {
      return await requestWithAccessKey<AccountUpdateOKResponse>({
        accessKey,
        key: getAccountAddRoleMutationKey(accountHandle, roleId),
        method: "PUT",
      });
    },
    categoryCreate: async (categoryCreateBody) => {
      return await requestWithAccessKey<CategoryCreateOKResponse>({
        accessKey,
        key: getCategoryCreateMutationKey(),
        method: "POST",
        data: categoryCreateBody,
      });
    },
    threadCreate: async (threadCreateBody) => {
      return await requestWithAccessKey<ThreadCreateOKResponse>({
        accessKey,
        key: getThreadCreateMutationKey(),
        method: "POST",
        data: threadCreateBody,
      });
    },
  };
}

async function requestWithAccessKey<T>({
  accessKey,
  key,
  method,
  data,
}: {
  accessKey: string;
  key: readonly [string];
  method: "POST" | "PUT";
  data?: unknown;
}): Promise<T> {
  const headers: Record<string, string> = {
    Authorization: `Bearer ${accessKey}`,
  };
  if (data !== undefined) {
    headers["Content-Type"] = "application/json";
  }

  const request = buildRequest({
    url: key[0],
    method,
    data,
    headers,
  });

  const response = await fetch(request);
  return await buildResult<T>(response);
}
