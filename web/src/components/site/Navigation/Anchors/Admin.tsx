import { AdminIcon } from "@/components/ui/icons/Admin";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const AdminID = "admin";
export const AdminRoute = "/admin";
export const AdminLabel = "Admin";

type Props = AnchorProps & LinkButtonStyleProps;

export function AdminAnchor(props: Props) {
  return (
    <Anchor
      id={AdminID}
      route={AdminRoute}
      label={AdminLabel}
      icon={<AdminIcon />}
      {...props}
    />
  );
}

export function AdminMenuItem() {
  return (
    <MenuItem
      id={AdminID}
      route={AdminRoute}
      label={AdminLabel}
      icon={<AdminIcon />}
    />
  );
}
