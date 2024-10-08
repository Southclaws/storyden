import { BellIcon } from "@heroicons/react/24/outline";

import { IconButton } from "@/components/ui/icon-button";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { AnchorProps, MenuItem } from "../Anchors/Anchor";

export const NotificationsID = "notifications";
export const NotificationsRoute = "/notifications";
export const NotificationsLabel = "Notifications";
export const NotificationsIcon = <BellIcon />;

export function NotificationAction({
  hideLabel,
  ...props
}: AnchorProps & LinkButtonStyleProps) {
  return (
    <IconButton {...props}>
      {NotificationsIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{NotificationsLabel}</span>
        </>
      )}
    </IconButton>
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
