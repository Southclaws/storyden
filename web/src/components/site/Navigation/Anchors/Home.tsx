import { HomeIcon } from "@/components/ui/icons/Home";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const HomeID = "home";
export const HomeRoute = "/";
export const HomeLabel = "Home";

export function HomeAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={HomeID}
      route={HomeRoute}
      label={HomeLabel}
      icon={<HomeIcon />}
      {...props}
    />
  );
}

export function HomeMenuItem() {
  return (
    <MenuItem
      id={HomeID}
      route={HomeRoute}
      label={HomeLabel}
      icon={<HomeIcon />}
    />
  );
}
