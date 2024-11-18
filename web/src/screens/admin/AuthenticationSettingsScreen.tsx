"use client";

import { AuthenticationSettingsForm } from "@/components/admin/AuthenticationSettings/AuthenticationSettingsForm";
import { UnreadyBanner } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

export function AuthenticationSettingsScreen() {
  const { ready, error, settings } = useSettings();
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  return <AuthenticationSettingsForm settings={settings} />;
}
