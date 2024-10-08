import { PlusIcon } from "lucide-react";

import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const ComposeID = "compose";
export const ComposeRoute = "/new";
export const ComposeLabel = "Post";
export const ComposeIcon = <PlusIcon />;

export function ComposeAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={ComposeID}
      route={ComposeRoute}
      label={ComposeLabel}
      icon={ComposeIcon}
      {...props}
    />
  );
}

export function ComposeMenuItem() {
  return (
    <MenuItem
      id={ComposeID}
      route={ComposeRoute}
      label={ComposeLabel}
      icon={ComposeIcon}
    />
  );
}
