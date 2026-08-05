"use client";

import { useNodeList } from "@/api/openapi-client/nodes";
import { type Account, NodeListResult, Visibility } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  initialSession?: Account;
  initialNodeList?: NodeListResult;
  currentPath?: string;
  visibility: Visibility[];
};

export function useLibraryNavigationTree({
  visibility,
  initialSession,
  initialNodeList,
}: Props) {
  const session = useSession(initialSession);
  const { data, error } = useNodeList(
    {
      visibility,
    },
    {
      swr: {
        fallbackData: initialNodeList,
      },
    },
  );
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const canManageLibrary = hasPermission(session, "MANAGE_LIBRARY");

  return {
    ready: true as const,
    data,
    canManageLibrary,
  };
}
