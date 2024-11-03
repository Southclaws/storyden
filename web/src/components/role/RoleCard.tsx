import { Role } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { css } from "@/styled-system/css";
import { CardBox, HStack } from "@/styled-system/jsx";

import { RoleEditModalTrigger } from "./RoleEdit/RoleEditModal";
import { badgeColourCSS } from "./colours";

type Props = {
  role: Role;
  editable?: boolean;
};

export function RoleCard({ role, editable }: Props) {
  const cssVars = badgeColourCSS(role.colour);

  const permissionCount = role.permissions.length;

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
      <HStack w="full" justify="space-between">
        <Heading>{role.name}</Heading>

        {editable && <RoleEditModalTrigger role={role} />}
      </HStack>

      <p>{permissionCount} permissions</p>
    </CardBox>
  );
}
