"use client";

import { UnreadyBanner } from "@/components/site/Unready";
import { useSettings } from "@/lib/settings/settings-client";

import { BrandSettingsForm } from "../../components/admin/BrandSettings/BrandSettings";

export function BrandSettingsScreen() {
  const { ready, error, settings } = useSettings();
  if (!ready) {
    return <UnreadyBanner error={error} />;
  }

  return <BrandSettingsForm settings={settings} />;
}
