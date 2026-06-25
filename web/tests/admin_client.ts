import { buildRequest, buildResult } from "../src/api/common";
import { getAccountAddRoleMutationKey } from "../src/api/openapi-client/accounts";
import { getAdminSettingsUpdateMutationKey } from "../src/api/openapi-client/admin";
import { getCategoryCreateMutationKey } from "../src/api/openapi-client/categories";
import { getNodeCreateMutationKey } from "../src/api/openapi-client/nodes";
import { getReplyCreateMutationKey } from "../src/api/openapi-client/replies";
import {
  getRobotCreateMutationKey,
  getRobotGetKey,
  getRobotProviderUpdateMutationKey,
} from "../src/api/openapi-client/robots";
import { getThreadCreateMutationKey } from "../src/api/openapi-client/threads";
import {
  AccountUpdateOKResponse,
  AdminSettingsUpdateBody,
  AdminSettingsUpdateOKResponse,
  CategoryCreateBody,
  CategoryCreateOKResponse,
  NodeCreateBody,
  NodeCreateOKResponse,
  ReplyCreateBody,
  ReplyCreateOKResponse,
  RobotCreateBody,
  RobotCreateOKResponse,
  RobotGetOKResponse,
  RobotProviderGetOKResponse,
  RobotProviderUpdateBody,
  ThreadCreateBody,
  ThreadCreateOKResponse,
} from "../src/api/openapi-schema";

export type AccessKeyClient = {
  accountAddRole: (
    accountHandle: string,
    roleId: string,
  ) => Promise<AccountUpdateOKResponse>;
  adminSettingsUpdate: (
    adminSettingsUpdateBody: AdminSettingsUpdateBody,
  ) => Promise<AdminSettingsUpdateOKResponse>;
  categoryCreate: (
    categoryCreateBody: CategoryCreateBody,
  ) => Promise<CategoryCreateOKResponse>;
  threadCreate: (
    threadCreateBody: ThreadCreateBody,
  ) => Promise<ThreadCreateOKResponse>;
  replyCreate: (
    threadSlug: string,
    replyCreateBody: ReplyCreateBody,
  ) => Promise<ReplyCreateOKResponse>;
  nodeCreate: (
    nodeCreateBody: NodeCreateBody,
  ) => Promise<NodeCreateOKResponse>;
  robotCreate: (
    robotCreateBody: RobotCreateBody,
  ) => Promise<RobotCreateOKResponse>;
  robotGet: (robotId: string) => Promise<RobotGetOKResponse>;
  robotProviderUpdate: (
    provider: string,
    robotProviderUpdateBody: RobotProviderUpdateBody,
  ) => Promise<RobotProviderGetOKResponse>;
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
    adminSettingsUpdate: async (adminSettingsUpdateBody) => {
      return await requestWithAccessKey<AdminSettingsUpdateOKResponse>({
        accessKey,
        key: getAdminSettingsUpdateMutationKey(),
        method: "PATCH",
        data: adminSettingsUpdateBody,
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
    replyCreate: async (threadSlug, replyCreateBody) => {
      return await requestWithAccessKey<ReplyCreateOKResponse>({
        accessKey,
        key: getReplyCreateMutationKey(threadSlug),
        method: "POST",
        data: replyCreateBody,
      });
    },
    nodeCreate: async (nodeCreateBody) => {
      return await requestWithAccessKey<NodeCreateOKResponse>({
        accessKey,
        key: getNodeCreateMutationKey(),
        method: "POST",
        data: nodeCreateBody,
      });
    },
    robotCreate: async (robotCreateBody) => {
      return await requestWithAccessKey<RobotCreateOKResponse>({
        accessKey,
        key: getRobotCreateMutationKey(),
        method: "POST",
        data: robotCreateBody,
      });
    },
    robotGet: async (robotId) => {
      return await requestWithAccessKey<RobotGetOKResponse>({
        accessKey,
        key: getRobotGetKey(robotId),
        method: "GET",
      });
    },
    robotProviderUpdate: async (provider, robotProviderUpdateBody) => {
      return await requestWithAccessKey<RobotProviderGetOKResponse>({
        accessKey,
        key: getRobotProviderUpdateMutationKey(provider),
        method: "PATCH",
        data: robotProviderUpdateBody,
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
  method: "GET" | "PATCH" | "POST" | "PUT";
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
