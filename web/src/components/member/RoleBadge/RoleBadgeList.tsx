import { AccountRoleList } from "@/api/openapi-schema";
import { HStack } from "@/styled-system/jsx";

import { RoleBadge } from "./RoleBadge";

export type Props = {
  roles: AccountRoleList;
  onlyBadgeRole?: boolean;
};

export function RoleBadgeList({ roles, onlyBadgeRole }: Props) {
  const filtered = onlyBadgeRole
    ? filterBadgeRole(roles)
    : filterDefaults(roles);

  return (
    <HStack flexWrap="wrap">
      {filtered.map((r) => (
        <RoleBadge key={r.id} role={r} />
      ))}
    </HStack>
  );
}

function filterBadgeRole(roles: AccountRoleList) {
  return roles.filter((r) => r.badge);
}

function filterDefaults(roles: AccountRoleList) {
  return roles.filter((r) => {
    if (r.default) {
      if (r.name === "Admin") {
        return true;
      }

      return false;
    }

    return true;
  });
}
