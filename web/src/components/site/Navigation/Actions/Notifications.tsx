import { BellIcon } from "@heroicons/react/24/outline";

import { IconButton } from "@/components/ui/icon-button";
import { LinkButtonStyleProps } from "@/components/ui/link-button";
import { Box } from "@/styled-system/jsx";

import { AnchorProps, MenuItem } from "../Anchors/Anchor";

export const NotificationsID = "notifications";
export const NotificationsRoute = "/notifications";
export const NotificationsLabel = "Notifications";
export const NotificationsIcon = <BellIcon />;

type Props = {
  unread: boolean;
};

export function NotificationAction({
  hideLabel,
  unread,
  ...props
}: AnchorProps & LinkButtonStyleProps & Props) {
  return (
    <IconButton size="sm" {...props}>
      {NotificationsIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{NotificationsLabel}</span>
        </>
      )}

      {unread && (
        <Box
          position="absolute"
          top="1"
          right="1"
          bgColor="red.8"
          borderRadius="full"
          w="2"
          h="2"
        />
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
