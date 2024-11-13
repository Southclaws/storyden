import { LinkIcon } from "@/components/ui/icons/Link";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const LinksID = "links";
export const LinksRoute = "/links";
export const LinksLabel = "Links";

export function LinksAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={LinksID}
      route={LinksRoute}
      label={LinksLabel}
      icon={<LinkIcon />}
      {...props}
    />
  );
}

export function LinksMenuItem() {
  return (
    <MenuItem
      id={LinksID}
      route={LinksRoute}
      label={LinksLabel}
      icon={<LinkIcon />}
    />
  );
}
