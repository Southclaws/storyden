import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermissionOr } from "@/utils/permissions";

import { useWatch } from "./store";

export function useLibraryPagePermissions() {
  const account = useSession();
  const owner = useWatch((s) => s.draft.owner);

  const isAllowedToEdit = hasPermissionOr(
    account,
    () => account?.id === owner.id,
    Permission.MANAGE_LIBRARY,
  );

  const isAllowedToDelete = hasPermissionOr(
    account,
    () => account?.id === owner.id,
    Permission.MANAGE_LIBRARY,
  );

  return {
    isAllowedToEdit,
    isAllowedToDelete,
  };
}
