import { DraftIcon } from "@/components/ui/icons/Draft";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const DraftsID = "drafts";
export const DraftsRoute = "/drafts";
export const DraftsLabel = "Drafts";
export const DraftsIcon = <DraftIcon />;

export function DraftsAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={DraftsID}
      route={DraftsRoute}
      label={DraftsLabel}
      icon={DraftsIcon}
      {...props}
    />
  );
}

export function DraftsMenuItem() {
  return (
    <MenuItem
      id={DraftsID}
      route={DraftsRoute}
      label={DraftsLabel}
      icon={DraftsIcon}
    />
  );
}
