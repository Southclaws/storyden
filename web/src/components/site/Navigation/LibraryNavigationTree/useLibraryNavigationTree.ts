"use client";

import { useNodeList } from "@/api/openapi-client/nodes";
import { NodeListResult, Visibility } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  initialNodeList?: NodeListResult;
  currentNode: string | undefined;
  visibility: Visibility[];
};

export function useLibraryNavigationTree({
  visibility,
  initialNodeList,
}: Props) {
  const session = useSession();
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
