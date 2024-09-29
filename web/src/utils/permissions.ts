import { flatten } from "lodash";

import { Account, Permission } from "@/api/openapi-schema";

export function hasPermission(account?: Account, ...permissions: Permission[]) {
  if (!account) return false;

  // extract each permission from each role
  const accountPermissions = new Set(
    flatten(account.roles.map((role) => role.permissions)),
  );

  if (accountPermissions.has("ADMINISTRATOR")) {
    return true;
  }

  return permissions.every((permission) => accountPermissions.has(permission));
}
