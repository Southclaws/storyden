"use client";

import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { DiscoverySettingsForm } from "@/components/admin/DiscoverySettings/DiscoverySettings";
import { UnreadyBanner } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";

export function DiscoverySettingsScreen() {
  const { error, data } = useAdminSettingsGet();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <DiscoverySettingsForm settings={parseAdminSettings(data)} />;
}
