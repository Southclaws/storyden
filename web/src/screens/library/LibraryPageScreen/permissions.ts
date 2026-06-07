import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission, hasPermissionOr } from "@/utils/permissions";

import { useWatch } from "./store";

export function useLibraryPagePermissions() {
  const account = useSession();
  const owner = useWatch((s) => s.draft.owner);

  const isLibraryManager = hasPermission(account, Permission.MANAGE_LIBRARY);
  const canSubmitNodeChanges = hasPermission(
    account,
    Permission.SUBMIT_LIBRARY_NODE_CHANGES,
  );

  const isAllowedToDirectEdit = isLibraryManager;

  const isAllowedToProposeEdit = isLibraryManager || canSubmitNodeChanges;
  const isAllowedToEdit = isAllowedToDirectEdit || isAllowedToProposeEdit;

  const isAllowedToDelete = hasPermissionOr(
    account,
    () => account?.id === owner.id,
    Permission.MANAGE_LIBRARY,
  );

  return {
    isLibraryManager,
    isAllowedToEdit,
    isAllowedToDirectEdit,
    isAllowedToProposeEdit,
    isAllowedToDelete,
  };
}
