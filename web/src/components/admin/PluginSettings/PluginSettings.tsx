import { PluginList } from "@/api/openapi-schema";

import { PluginSettingsIndex } from "./PluginSettingsIndex";
import { PluginStatus } from "./PluginStatus";
import { useSelectedPlugin } from "./useSelectedPlugin";

type Props = {
  plugins: PluginList;
};

export function PluginSettings({ plugins }: Props) {
  const [pluginID] = useSelectedPlugin();

  if (pluginID) {
    return <PluginStatus plugin={pluginID} />;
  }

  return <PluginSettingsIndex plugins={plugins} />;
}
