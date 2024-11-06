import { Role } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { css } from "@/styled-system/css";
import { CardBox, WStack } from "@/styled-system/jsx";

import { PermissionSummary } from "./PermissionList";
import { RoleEditModalTrigger } from "./RoleEdit/RoleEditModal";
import { badgeColourCSS } from "./colours";

type Props = {
  role: Role;
  editable?: boolean;
};

export function RoleCard({ role, editable }: Props) {
  const cssVars = badgeColourCSS(role.colour);

  return (
    <CardBox
      className={css({
        borderColor: "colorPalette.text",
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

      <PermissionSummary permissions={role.permissions} />
    </CardBox>
  );
}
