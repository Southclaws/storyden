"use client";

import { usePluginList } from "@/api/openapi-client/plugins";
import { PluginSettingsIndex } from "@/components/admin/PluginSettings/PluginSettingsIndex";
import { Unready } from "@/components/site/Unready";

export function PluginSettingsScreen() {
  const { data, error } = usePluginList();
  if (!data) {
    return <Unready error={error} />;
  }

  return <PluginSettingsIndex plugins={data.plugins} />;
}
