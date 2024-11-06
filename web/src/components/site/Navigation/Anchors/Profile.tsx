import { ProfileIcon } from "@/components/ui/icons/Profile";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, MenuItem } from "./Anchor";

export const ProfileID = "profile";
export const ProfileRoute = (handle: string) => `/m/${handle}`;
export const ProfileLabel = "Profile";

export function ProfileAnchor({
  handle,
  ...props
}: LinkButtonStyleProps & { handle: string }) {
  return (
    <Anchor
      id={ProfileID}
      route={ProfileRoute(handle)}
      label={ProfileLabel}
      icon={<ProfileIcon />}
      {...props}
    />
  );
}

export function ProfileMenuItem({ handle }: { handle: string }) {
  return (
    <MenuItem
      id={ProfileID}
      route={ProfileRoute(handle)}
      label={ProfileLabel}
      icon={<ProfileIcon />}
    />
  );
}
