import { PluginActiveState, PluginList } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { LStack, WStack, styled } from "@/styled-system/jsx";

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
        <Heading size="md">Plugins</Heading>

        <PluginAddTrigger />
      </WStack>

      <styled.p color="text.subtle">
        {hasInactive ? (
          <span>
            {totalPlugins} plugins, {activePlugins} active.
          </span>
        ) : (
          <span>{plugins.length} plugins.</span>
        )}
      </styled.p>

      <PluginItemList plugins={plugins} />
    </LStack>
  );
}
