import { Permission } from "@/api/openapi-schema";
import { Tooltip } from "@/components/ui/tooltip";
import { PermissionDetails } from "@/lib/permission/permission";

import { Badge } from "../ui/badge";

type Props = {
  permission: Permission;
};

export function PermissionBadge(props: Props) {
  const p = PermissionDetails[props.permission];

  return (
    <Tooltip
      content={p.description}
      contentProps={{ p: "2", borderRadius: "2xl" }}
      openDelay={0}
      positioning={{
        shift: 16,
      }}
    >
      <Badge cursor="pointer">{p.name}</Badge>
    </Tooltip>
  );
}
