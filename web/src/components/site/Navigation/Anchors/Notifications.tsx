import { NotificationIcon } from "@/components/ui/icons/Notification";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const NotificationsID = "notifications";
export const NotificationsRoute = "/notifications";
export const NotificationsLabel = "Notifications";
export const NotificationsIcon = <NotificationIcon />;

export function NotificationsAnchor(props: AnchorProps & LinkButtonStyleProps) {
  return (
    <Anchor
      id={NotificationsID}
      route={NotificationsRoute}
      label={NotificationsLabel}
      icon={NotificationsIcon}
      {...props}
    />
  );
}

export function NotificationsMenuItem() {
  return (
    <MenuItem
      id={NotificationsID}
      route={NotificationsRoute}
      label={NotificationsLabel}
      icon={NotificationsIcon}
    />
  );
}
