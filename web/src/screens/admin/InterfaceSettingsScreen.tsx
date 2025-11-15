"use client";

import { UnreadyBanner } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

import { InterfaceSettingsForm } from "../../components/admin/InterfaceSettings/InterfaceSettings";

export function InterfaceSettingsScreen() {
  const { ready, error, settings } = useSettings();
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  return <InterfaceSettingsForm settings={settings} />;
}
