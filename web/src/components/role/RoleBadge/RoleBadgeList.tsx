import { AccountRoleRefList } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import * as Popover from "@/components/ui/popover";
import { HStack } from "@/styled-system/jsx";

import { RoleBadge } from "./RoleBadge";

export type Props = {
  roles: AccountRoleRefList;
  onlyBadgeRole?: boolean;
  limit?: number;
};

export function RoleBadgeList({ roles, onlyBadgeRole, limit }: Props) {
  const filtered = onlyBadgeRole
    ? filterBadgeRole(roles)
    : filterDefaults(roles);

  const isLimited = limit && filtered.length > limit;
  const rest = filtered.length - (limit ?? 0);

  const limited = filtered.slice(0, limit);

  return (
    <HStack flexWrap="wrap" py="1">
      {limited.map((r) => (
        <RoleBadge key={r.id} role={r} />
      ))}
      {isLimited && (
        <Popover.Root>
          <Popover.Trigger>
            <Badge color="fg.muted" size="sm">
              +{rest}
            </Badge>
          </Popover.Trigger>
          <Popover.Positioner>
            <Popover.Content p="2" borderRadius="2xl">
              <Popover.Arrow>
                <Popover.ArrowTip />
              </Popover.Arrow>
              <Popover.Description>
                <HStack flexWrap="wrap">
                  {filtered.map((r) => (
                    <RoleBadge key={r.id} role={r} />
                  ))}
                </HStack>
              </Popover.Description>
            </Popover.Content>
          </Popover.Positioner>
        </Popover.Root>
      )}
    </HStack>
  );
}

function filterBadgeRole(roles: AccountRoleRefList) {
  return roles.filter((r) => r.badge);
}

function filterDefaults(roles: AccountRoleRefList) {
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
