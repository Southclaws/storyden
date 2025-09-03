import { Role } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { isDefaultRole, isStoredDefaultRole } from "@/lib/role/defaults";
import { css } from "@/styled-system/css";
import { CardBox, WStack } from "@/styled-system/jsx";

import { Badge } from "../ui/badge";

import { PermissionSummary } from "./PermissionList";
import { RoleEditModalTrigger } from "./RoleEdit/RoleEditModal";
import { badgeColourCSS } from "./colours";

type Props = {
  role: Role;
  editable?: boolean;
};

export function RoleCard({ role, editable }: Props) {
  const cssVars = badgeColourCSS(role.colour);

  const isDefault = isDefaultRole(role);
  const isCustomDefault = isStoredDefaultRole(role);

  return (
    <CardBox
      className={css({
        borderColor: "colorPalette.fg",
      })}
      style={{
        ...cssVars,
        borderLeftWidth: "thick",
        borderLeftStyle: "solid",
      }}
    >
      <WStack>
        <Heading>{role.name}</Heading>

        {editable && <RoleEditModalTrigger role={role} />}
      </WStack>

      <WStack>
        <PermissionSummary permissions={role.permissions} />
        {isDefault && (
          <>
            {isCustomDefault ? (
              <Badge size="sm">Default + Custom</Badge>
            ) : (
              <Badge size="sm">Default</Badge>
            )}
          </>
        )}
      </WStack>
    </CardBox>
  );
}
