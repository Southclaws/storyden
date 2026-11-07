"use client";

import { useSession } from "@/auth";
import { MemberNotificationSettings } from "@/components/settings/MemberNotificationSettings/MemberNotificationSettings";
import { Unready } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

export function MemberNotificationSettingsScreen() {
  const session = useSession();
  const settings = useSettings();

  if (!session || !settings.ready) {
    return <Unready />;
  }

  return (
    <MemberNotificationSettings
      session={session}
      settings={settings.settings}
    />
  );
}
