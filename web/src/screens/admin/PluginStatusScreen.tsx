"use client";

import { PluginStatus } from "@/components/admin/PluginSettings/PluginStatus";

type Props = {
  pluginId: string;
};

export function PluginStatusScreen({ pluginId }: Props) {
  return <PluginStatus plugin={pluginId} />;
}
