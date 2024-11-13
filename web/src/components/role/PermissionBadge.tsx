import { Permission } from "@/api/openapi-schema";
import * as Tooltip from "@/components/ui/tooltip";
import { PermissionDetails } from "@/lib/permission/permission";

import { Badge } from "../ui/badge";

type Props = {
  permission: Permission;
};

export function PermissionBadge(props: Props) {
  const p = PermissionDetails[props.permission];

  return (
    <Tooltip.Root
      openDelay={0}
      positioning={{
        slide: true,
        shift: 16,
      }}
    >
      <Tooltip.Trigger asChild>
        <Badge cursor="pointer">{p.name}</Badge>
      </Tooltip.Trigger>
      <Tooltip.Positioner>
        <Tooltip.Arrow>
          <Tooltip.ArrowTip />
        </Tooltip.Arrow>

        <Tooltip.Content p="2" borderRadius="2xl">
          {p.description}
        </Tooltip.Content>
      </Tooltip.Positioner>
    </Tooltip.Root>
  );
}
