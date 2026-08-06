import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { NotificationIcon } from "@/components/ui/icons/Notification";

import { AnchorProps, MenuItem } from "../site/Navigation/Anchors/Anchor";

export const NotificationsID = "notifications";
export const NotificationsRoute = "/notifications";
export const NotificationsLabel = "Notifications";
export const NotificationsIcon = <NotificationIcon />;

type Props = {
  unread?: boolean;
};

export function NotificationsTrigger({
  hideLabel,
  unread,
  ...props
}: AnchorProps & ButtonProps & Props) {
  const accessibleLabel = unread
    ? `${NotificationsLabel}, unread`
    : NotificationsLabel;

  return (
    <IconButton size="sm" aria-label={accessibleLabel} {...props}>
      {NotificationsIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{NotificationsLabel}</span>
        </>
      )}

      {unread && (
        <span
          aria-hidden="true"
          className="notifications-trigger__unread-indicator"
        />
      )}
    </IconButton>
  );
}
