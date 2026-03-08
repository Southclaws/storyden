import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import { usePluginGet } from "@/api/openapi-client/plugins";
import { Identifier, Plugin, PluginExternalProps } from "@/api/openapi-schema";
import { BackAction } from "@/components/site/Action/Back";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import * as Tabs from "@/components/ui/tabs";
import { CardBox, HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { PluginStatusToggle } from "../PluginStatusToggle";

import { DEFAULT_PLUGIN_TAB, PluginTab, isPluginTab } from "./model";
import { ConfigurationTab } from "./tabs/ConfigurationTab";
import { ConnectionTab } from "./tabs/ConnectionTab";
import { LogsTab } from "./tabs/LogsTab";
import { ManifestTab } from "./tabs/ManifestTab";
import { OverviewTab } from "./tabs/OverviewTab";
import { PackageTab } from "./tabs/PackageTab";

type Props = {
  plugin: Identifier;
};

export function PluginStatus({ plugin }: Props) {
  const { data, error } = usePluginGet(plugin);
  const [tab, setTab] = useQueryState("plugin-tab", {
    defaultValue: DEFAULT_PLUGIN_TAB,
  });
  const isExternal = data?.connection.mode === "external";

  useEffect(() => {
    if (!data) {
      return;
    }

    if (!isPluginTab(tab, isExternal)) {
      setTab(DEFAULT_PLUGIN_TAB);
    }
  }, [data, tab, isExternal, setTab]);

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  function handleTabChange({ value }: TabsValueChangeDetails) {
    setTab(value as PluginTab);
  }

  return (
    <CardBox className={lstack()} gap="2">
      <WStack justifyContent="space-between">
        <HStack>
          <BackAction />
          <Heading size="md" lineClamp="1">
            {data.name}
          </Heading>
        </HStack>

        <PluginStatusToggle plugin={data} />
      </WStack>

      <Tabs.Root
        lazyMount={true}
        width="full"
        variant="line"
        defaultValue={DEFAULT_PLUGIN_TAB}
        value={tab ?? DEFAULT_PLUGIN_TAB}
        onValueChange={handleTabChange}
      >
        <Tabs.List>
          <Tabs.Trigger value="overview">Overview</Tabs.Trigger>
          <Tabs.Trigger value="configuration">Configuration</Tabs.Trigger>
          <Tabs.Trigger value="manifest">Manifest</Tabs.Trigger>
          {!isExternal && <Tabs.Trigger value="package">Package</Tabs.Trigger>}
          {isExternal ? (
            <Tabs.Trigger value="connection">Connection</Tabs.Trigger>
          ) : (
            <Tabs.Trigger value="logs">Logs</Tabs.Trigger>
          )}
          <Tabs.Indicator />
        </Tabs.List>

        <Tabs.Content value="overview">
          <OverviewTab plugin={data} />
        </Tabs.Content>

        <Tabs.Content value="manifest">
          <ManifestTab
            pluginID={data.id}
            manifest={data.manifest}
            editable={isExternal}
          />
        </Tabs.Content>

        <Tabs.Content value="configuration">
          <ConfigurationTab pluginID={data.id} />
        </Tabs.Content>

        {!isExternal && (
          <Tabs.Content value="package">
            <PackageTab plugin={data} />
          </Tabs.Content>
        )}

        {!isExternal && (
          <Tabs.Content value="logs">
            <LogsTab plugin={data} />
          </Tabs.Content>
        )}

        {isExternal && (
          <Tabs.Content value="connection">
            <ConnectionTab
              plugin={data as Plugin & { connection: PluginExternalProps }}
            />
          </Tabs.Content>
        )}
      </Tabs.Root>
    </CardBox>
  );
}
