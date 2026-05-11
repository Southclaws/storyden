import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { NotificationIcon } from "@/components/ui/icons/Notification";
import { useI18n } from "@/i18n/provider";
import { Box } from "@/styled-system/jsx";

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
  const { t } = useI18n();

  return (
    <IconButton size="sm" aria-label={t("Notifications")} {...props}>
      {NotificationsIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{t(NotificationsLabel)}</span>
        </>
      )}

      {unread && (
        <Box
          position="absolute"
          top="1"
          right="1"
          bgColor="fg.destructive/60"
          borderRadius="full"
          w="2"
          h="2"
        />
      )}
    </IconButton>
  );
}
