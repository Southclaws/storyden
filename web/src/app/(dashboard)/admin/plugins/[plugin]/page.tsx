"use client";

import { useParams } from "next/navigation";

import { PluginStatusScreen } from "@/screens/admin/PluginStatusScreen";

export default function Page() {
  const { plugin } = useParams<{ plugin: string }>();

  return <PluginStatusScreen pluginId={plugin} />;
}
