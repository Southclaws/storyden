import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const DiscussionID = "discussion";
export const DiscussionRoute = "/d";
export const DiscussionLabel = "Discussion";

export function HomeAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={DiscussionID}
      route={DiscussionRoute}
      label={DiscussionLabel}
      icon={<DiscussionIcon />}
      {...props}
    />
  );
}

export function DiscussionMenuItem() {
  return (
    <MenuItem
      id={DiscussionID}
      route={DiscussionRoute}
      label={DiscussionLabel}
      icon={<DiscussionIcon />}
    />
  );
}
