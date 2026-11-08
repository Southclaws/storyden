"use client";

import { useNotificationSettingsGet } from "@/api/openapi-client/notifications";
import { useSession } from "@/auth";
import { MemberNotificationSettings } from "@/components/settings/MemberNotificationSettings/MemberNotificationSettings";
import { Unready, UnreadyBanner } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

export function MemberNotificationSettingsScreen() {
  const session = useSession();
  const settings = useSettings();
  const notifications = useNotificationSettingsGet();

  if (!notifications.data) {
    return <UnreadyBanner error={notifications.error} />;
  }

  if (!session || !settings.ready) {
    return <Unready />;
  }

  return (
    <MemberNotificationSettings
      session={session}
      settings={settings.settings}
      notificationSettings={notifications.data}
    />
  );
}
