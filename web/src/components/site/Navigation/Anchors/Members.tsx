import { UsersIcon } from "@heroicons/react/24/outline";

import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const MembersID = "members";
export const MembersRoute = "/m";
export const MembersLabel = "Members";
export const MembersIcon = <UsersIcon />;

export function MembersAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={MembersID}
      route={MembersRoute}
      label={MembersLabel}
      icon={MembersIcon}
      {...props}
    />
  );
}

export function MembersMenuItem() {
  return (
    <MenuItem
      id={MembersID}
      route={MembersRoute}
      label={MembersLabel}
      icon={MembersIcon}
    />
  );
}
