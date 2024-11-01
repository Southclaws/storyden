import { Account, Collection, Permission } from "@/api/openapi-schema";
import { hasPermissionOr } from "@/utils/permissions";

export function canDeleteCollection(col: Collection, account?: Account) {
  if (!account) return false;

  return hasPermissionOr(
    account,
    () => col.owner.id === account.id,
    Permission.MANAGE_COLLECTIONS,
  );
}

export function canEditCollection(col: Collection, account?: Account) {
  if (!account) return false;

  return hasPermissionOr(
    account,
    () => col.owner.id === account.id,
    Permission.MANAGE_COLLECTIONS,
  );
}
