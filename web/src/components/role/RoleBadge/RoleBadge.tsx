import { Role } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";

import { badgeColourCSS } from "../colours";

export type Props = {
  role: Role;
};

export function RoleBadge({ role }: Props) {
  const cssVars = badgeColourCSS(role.colour);

  return (
    <Badge
      size="sm"
      style={cssVars}
      bgColor="colorPalette.bg"
      borderColor="colorPalette.border"
      color="colorPalette.fg"
    >
      {role.name}
    </Badge>
  );
}
