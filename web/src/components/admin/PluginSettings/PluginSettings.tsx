import { PluginActiveState, PluginList } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { PluginAddTrigger } from "./PluginAddModal";
import { PluginItemList } from "./PluginItemList";
import { getPluginActiveState } from "./utils";

type Props = {
  plugins: PluginList;
};

export function PluginSettings({ plugins }: Props) {
  const totalPlugins = plugins.length;
  const activePlugins = plugins.filter(
    (plugin) => getPluginActiveState(plugin) === PluginActiveState.active,
  ).length;
  const hasInactive = totalPlugins !== activePlugins;

  return (
    <CardBox className={lstack()}>
      <WStack justifyContent="space-between">
        <Heading size="md">Plugins</Heading>

        <PluginAddTrigger />
      </WStack>

      <styled.p color="fg.muted">
        {hasInactive ? (
          <span>
            {totalPlugins} plugins, {activePlugins} active.
          </span>
        ) : (
          <span>{plugins.length} plugins.</span>
        )}
      </styled.p>

      <PluginItemList plugins={plugins} />
    </CardBox>
  );
}
