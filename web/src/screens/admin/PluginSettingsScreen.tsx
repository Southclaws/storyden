import { usePluginList } from "@/api/openapi-client/plugins";
import { PluginSettings } from "@/components/admin/PluginSettings/PluginSettings";
import { Unready } from "@/components/site/Unready";

export function PluginSettingsScreen() {
  const { data, error } = usePluginList();
  if (!data) {
    return <Unready error={error} />;
  }

  return <PluginSettings plugins={data.plugins} />;
}
