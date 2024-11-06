import { RolesIcon } from "@/components/ui/icons/Roles";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const RolesID = "roles";
export const RolesRoute = "/roles";
export const RolesLabel = "Roles";

export function RolesAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={RolesID}
      route={RolesRoute}
      label={RolesLabel}
      icon={<RolesIcon />}
      {...props}
    />
  );
}

export function RolesMenuItem() {
  return (
    <MenuItem
      id={RolesID}
      route={RolesRoute}
      label={RolesLabel}
      icon={<RolesIcon />}
    />
  );
}
