import { PluginActiveState, PluginList } from "@/api/openapi-schema";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { LStack, WStack } from "@/styled-system/jsx";

import { PluginAddTrigger } from "./PluginAddModal/PluginAddModal";
import { PluginItemList } from "./PluginItemList";
import { getPluginActiveState } from "./utils";

type Props = {
  plugins: PluginList;
};

export function PluginSettingsIndex({ plugins }: Props) {
  const totalPlugins = plugins.length;
  const activePlugins = plugins.filter(
    (plugin) => getPluginActiveState(plugin) === PluginActiveState.active,
  ).length;
  const hasInactive = totalPlugins !== activePlugins;

  return (
    <LStack>
      <WStack justifyContent="space-between">
        <PageHeading>Plugins</PageHeading>

        <PluginAddTrigger />
      </WStack>

      <Text variant="supporting">
        {hasInactive ? (
          <span>
            {totalPlugins} plugins, {activePlugins} active.
          </span>
        ) : (
          <span>{plugins.length} plugins.</span>
        )}
      </Text>

      <PluginItemList plugins={plugins} />
    </LStack>
  );
}
