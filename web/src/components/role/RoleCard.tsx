import { ReactNode } from "react";

import { Role } from "@/api/openapi-schema";
import { CardBox } from "@/components/ui/card-box";
import { Text } from "@/components/ui/text";
import { isDefaultRole, isStoredDefaultRole } from "@/lib/role/defaults";
import { css } from "@/styled-system/css";
import { HStack, WStack } from "@/styled-system/jsx";

import { Badge } from "../ui/badge";

import { PermissionSummary } from "./PermissionList";
import { RoleEditModalTrigger } from "./RoleEdit/RoleEditModal";
import { badgeColourCSS } from "./colours";

type Props = {
  role: Role;
  editable?: boolean;
  dragHandle?: ReactNode;
};

export function RoleCard({ role, editable, dragHandle }: Props) {
  const cssVars = badgeColourCSS(role.colour);

  const isDefault = isDefaultRole(role);
  const isCustomDefault = isStoredDefaultRole(role);

  return (
    <CardBox
      className={css({
        borderColor: "colorPalette.6",
        display: "flex",
        gap: "2",
      })}
      style={{
        ...cssVars,
        borderLeftWidth: "thick",
        borderLeftStyle: "solid",
      }}
    >
      <WStack alignItems="flex-start">
        <Text variant="supporting" color="text.default" fontWeight="semibold">
          {role.name}
        </Text>

        <HStack>
          {isDefault && (
            <>
              {isCustomDefault ? (
                <Badge size="sm">Default + Custom</Badge>
              ) : (
                <Badge size="sm">Default</Badge>
              )}
            </>
          )}

          {dragHandle}
        </HStack>
      </WStack>

      <WStack alignItems="flex-end">
        <PermissionSummary permissions={role.permissions} />

        {editable && <RoleEditModalTrigger role={role} />}
      </WStack>
    </CardBox>
  );
}
