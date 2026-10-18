"use client";

import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { MediaSettingsForm } from "@/components/admin/MediaSettings/MediaSettings";
import { UnreadyBanner } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";

export function MediaSettingsScreen() {
  const { error, data } = useAdminSettingsGet();
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const settings = parseAdminSettings(data);

  return <MediaSettingsForm settings={settings} />;
}
