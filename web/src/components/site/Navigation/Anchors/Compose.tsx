import { CreateIcon } from "@/components/ui/icons/Create";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const ComposeID = "compose";
export const ComposeRoute = "/new";
export const ComposeLabel = "Post";
export const ComposeIcon = <CreateIcon />;

export function ComposeAnchor({
  label = ComposeLabel,
  ...props
}: AnchorProps & LinkButtonStyleProps & { label?: string }) {
  return (
    <Anchor
      id={ComposeID}
      route={ComposeRoute}
      label={label}
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
