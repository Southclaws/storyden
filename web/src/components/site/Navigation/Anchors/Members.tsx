import { MembersIcon } from "@/components/ui/icons/Members";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const MembersID = "members";
export const MembersRoute = "/m";
export const MembersLabel = "Members";

export function MembersAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={MembersID}
      route={MembersRoute}
      label={MembersLabel}
      icon={<MembersIcon />}
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
      icon={<MembersIcon />}
    />
  );
}
