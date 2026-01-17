"use client";

import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { UnreadyBanner } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";

import { SystemSettingsForm } from "../../components/admin/SystemSettings/SystemSettings";

export function SystemSettingsScreen() {
  const { error, data } = useAdminSettingsGet();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const settings = parseAdminSettings(data);

  return <SystemSettingsForm settings={settings} />;
}
