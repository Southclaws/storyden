import { NodeWithChildren, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermissionOr } from "@/utils/permissions";

export function useLibraryPagePermissions(node: NodeWithChildren) {
  const account = useSession();

  const isAllowedToEdit = hasPermissionOr(
    account,
    () => account?.id === node.owner.id,
    Permission.MANAGE_LIBRARY,
  );

  const isAllowedToDelete = hasPermissionOr(
    account,
    () => account?.id === node.owner.id,
    Permission.MANAGE_LIBRARY,
  );

  return {
    isAllowedToEdit,
    isAllowedToDelete,
  };
}
