"use client";

import { useNodeList } from "@/api/openapi-client/nodes";
import { Visibility } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  label: string;
  href: string;
  currentNode: string | undefined;
  visibility: Visibility[];
};

export function useLibraryNavigationTree({ visibility }: Props) {
  const session = useSession();
  const { data, error } = useNodeList({
    visibility,
  });
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
