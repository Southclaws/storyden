import { PluginActiveState, PluginList } from "@/api/openapi-schema";
import { Heading } from "@/components/ui/heading";
import { useI18n } from "@/i18n/provider";
import { CardBox, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { PluginAddTrigger } from "./PluginAddModal/PluginAddModal";
import { PluginItemList } from "./PluginItemList";
import { getPluginActiveState } from "./utils";

type Props = {
  plugins: PluginList;
};

export function PluginSettingsIndex({ plugins }: Props) {
  const { t } = useI18n();
  const totalPlugins = plugins.length;
  const activePlugins = plugins.filter(
    (plugin) => getPluginActiveState(plugin) === PluginActiveState.active,
  ).length;
  const hasInactive = totalPlugins !== activePlugins;

  return (
    <CardBox className={lstack()}>
      <WStack justifyContent="space-between">
        <Heading size="md">{t("Plugins")}</Heading>

        <PluginAddTrigger />
      </WStack>

      <styled.p color="fg.muted">
        {hasInactive ? (
          <span>
            {totalPlugins} {t("plugins")}, {activePlugins} {t("active")}.
          </span>
        ) : (
          <span>
            {plugins.length} {t("plugins")}.
          </span>
        )}
      </styled.p>

      <PluginItemList plugins={plugins} />
    </CardBox>
  );
}
