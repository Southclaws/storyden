"use client";

import { useAdminSettingsGet } from "@/api/openapi-client/admin";
import { UnreadyBanner } from "@/components/site/Unready";
import { parseAdminSettings } from "@/lib/settings/settings";

import { InterfaceSettingsForm } from "../../components/admin/InterfaceSettings/InterfaceSettings";

export function InterfaceSettingsScreen() {
  const { error, data } = useAdminSettingsGet();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return <InterfaceSettingsForm settings={parseAdminSettings(data)} />;
}
