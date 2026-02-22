import { Badge } from "@/components/ui/badge";

import { badgeColourCSS } from "../colours";

type RoleRef = {
  name: string;
  colour: string;
};

export type Props = {
  role: RoleRef;
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
