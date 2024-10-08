import { HomeIcon as HomeHeroIcon } from "@heroicons/react/24/outline";

import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const HomeID = "home";
export const HomeRoute = "/";
export const HomeLabel = "Home";
export const HomeIcon = <HomeHeroIcon />;

export function HomeAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={HomeID}
      route={HomeRoute}
      label={HomeLabel}
      icon={HomeIcon}
      {...props}
    />
  );
}

export function HomeMenuItem() {
  return (
    <MenuItem id={HomeID} route={HomeRoute} label={HomeLabel} icon={HomeIcon} />
  );
}
