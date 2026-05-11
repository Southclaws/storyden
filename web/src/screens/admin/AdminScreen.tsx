"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import * as Tabs from "@/components/ui/tabs";
import { useI18n } from "@/i18n/provider";
import { useCapability } from "@/lib/settings/capabilities";

import { AccessKeySettingsScreen } from "./AccessKeySettingsScreen";
import { AuditLogSettingsScreen } from "./AuditLogSettingsScreen/AuditLogSettingsScreen";
import { AuthenticationSettingsScreen } from "./AuthenticationSettingsScreen";
import { BrandSettingsScreen } from "./BrandSettingsScreen";
import { InterfaceSettingsScreen } from "./InterfaceSettingsScreen";
import { ModerationSettingsScreen } from "./ModerationSettingsScreen";
import { PluginSettingsScreen } from "./PluginSettingsScreen";
import { SystemSettingsScreen } from "./SystemSettingsScreen";

const DEFAULT_TAB = "brand";

export function AdminScreen() {
  const { t } = useI18n();
  const pluginsEnabled = useCapability("plugins");
  const [tab, setTab] = useQueryState("tab", {
    defaultValue: DEFAULT_TAB,
  });

  // NOTE: A hack because for some reason, the tab component renders twice and
  // the associated hook gets lost and results in `ready` always being false,
  // despite the useSettings hook returning the correct data. Not sure if this
  // is a Next.js bug, a React bug or a Ark, Park or something else bug...
  useEffect(() => {
    if (!tab) {
      setTab(DEFAULT_TAB);
      return;
    }

    if (!pluginsEnabled && tab === "plugins") {
      setTab(DEFAULT_TAB);
    }
  }, [pluginsEnabled, tab, setTab]);

  function handleTabChange({ value }: TabsValueChangeDetails) {
    setTab(value);
  }

  return (
    <Tabs.Root
      width="full"
      variant="enclosed"
      // variant="line"
      // variant="outline"
      defaultValue={DEFAULT_TAB}
      value={tab}
      onValueChange={handleTabChange}
    >
      <Tabs.List>
        <Tabs.Trigger value="brand">{t("Brand")}</Tabs.Trigger>
        <Tabs.Trigger value="moderation">{t("Moderation")}</Tabs.Trigger>
        <Tabs.Trigger value="system">{t("System")}</Tabs.Trigger>
        <Tabs.Trigger value="audit">{t("Audit Log")}</Tabs.Trigger>
        <Tabs.Trigger value="interface">{t("Interface")}</Tabs.Trigger>
        <Tabs.Trigger value="authentication">{t("Authentication")}</Tabs.Trigger>
        <Tabs.Trigger value="access_keys">{t("Access keys")}</Tabs.Trigger>
        {pluginsEnabled && (
          <Tabs.Trigger value="plugins">{t("Plugins")}</Tabs.Trigger>
        )}
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="brand">
        <BrandSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="moderation">
        <ModerationSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="system">
        <SystemSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="audit">
        <AuditLogSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="interface">
        <InterfaceSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="authentication">
        <AuthenticationSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="access_keys">
        <AccessKeySettingsScreen />
      </Tabs.Content>

      {pluginsEnabled && (
        <Tabs.Content value="plugins">
          <PluginSettingsScreen />
        </Tabs.Content>
      )}
    </Tabs.Root>
  );
}
