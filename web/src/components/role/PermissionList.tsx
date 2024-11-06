import { PermissionList } from "@/api/openapi-schema";
import * as Popover from "@/components/ui/popover";
import { PermissionDetails } from "@/lib/permission/permission";
import { Box, LStack } from "@/styled-system/jsx";

import { Badge } from "../ui/badge";
import { Button } from "../ui/button";

type Props = {
  permissions: PermissionList;
};

export function PermissionSummary({ permissions }: Props) {
  const permissionCount = permissions.length;

  const permissionCountLabel =
    permissionCount === 1 ? "permission" : "permissions";

  const permissionLabel = `${permissionCount} ${permissionCountLabel}`;

  return (
    <Popover.Root
      positioning={{
        slide: true,
        shift: 16,
      }}
    >
      <Popover.Trigger asChild>
        <Button size="xs" variant="link">
          {permissionLabel}
        </Button>
      </Popover.Trigger>
      <Popover.Positioner>
        <Popover.Content p="2" borderRadius="2xl">
          <Popover.Arrow>
            <Popover.ArrowTip />
          </Popover.Arrow>
          <Popover.Description>
            <LStack overflowY="scroll">
              {permissions.length > 0 ? (
                permissions.map((p) => (
                  <Badge key={p}>{PermissionDetails[p].name}</Badge>
                ))
              ) : (
                <Box pl="4">
                  <p>No permissions.</p>
                </Box>
              )}
            </LStack>
          </Popover.Description>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
