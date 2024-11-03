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
      bgColor="colorPalette"
      borderColor="colorPalette.muted"
      color="colorPalette.text"
    >
      {role.name}
    </Badge>
  );
}
