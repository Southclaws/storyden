import * as Tabs from "@/components/ui/tabs";
import { useI18n } from "@/i18n/provider";
import { UseDisclosureProps } from "@/utils/useDisclosure";

import { PluginAddExternal } from "./PluginAddExternal";
import { PluginAddUpload } from "./PluginAddUpload";

type Props = UseDisclosureProps;

export function PluginAddScreen({ onClose }: Props) {
  const { t } = useI18n();
  return (
    <Tabs.Root width="full" variant="line" defaultValue="supervised">
      <Tabs.List>
        <Tabs.Trigger value="supervised">{t("Upload")}</Tabs.Trigger>
        <Tabs.Trigger value="external">{t("External")}</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="supervised">
        <PluginAddUpload onClose={onClose} />
      </Tabs.Content>

      <Tabs.Content value="external">
        <PluginAddExternal onClose={onClose} />
      </Tabs.Content>
    </Tabs.Root>
  );
}
