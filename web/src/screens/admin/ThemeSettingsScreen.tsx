"use client";

import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { UnreadyBanner } from "@/components/site/Unready";
import { ThemeSettingsForm } from "@/components/admin/ThemeSettings/ThemeSettings";
import { parseAdminSettings } from "@/lib/settings/settings";

export function ThemeSettingsScreen() {
  const { error, data } = useAdminSettingsGet();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <ThemeSettingsForm settings={parseAdminSettings(data)} />;
}
