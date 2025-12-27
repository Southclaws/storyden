"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import * as Tabs from "@/components/ui/tabs";

import { AccessKeySettingsScreen } from "./AccessKeySettingsScreen";
import { AuthenticationSettingsScreen } from "./AuthenticationSettingsScreen";
import { BrandSettingsScreen } from "./BrandSettingsScreen";
import { InterfaceSettingsScreen } from "./InterfaceSettingsScreen";
import { ModerationSettingsScreen } from "./ModerationSettingsScreen";

const DEFAULT_TAB = "brand";

export function AdminScreen() {
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
    }
  }, [tab, setTab]);

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
        <Tabs.Trigger value="brand">Brand</Tabs.Trigger>
        <Tabs.Trigger value="moderation">Moderation</Tabs.Trigger>
        <Tabs.Trigger value="interface">Interface</Tabs.Trigger>
        <Tabs.Trigger value="authentication">Authentication</Tabs.Trigger>
        <Tabs.Trigger value="access_keys">Access keys</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="brand">
        <BrandSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="moderation">
        <ModerationSettingsScreen />
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
    </Tabs.Root>
  );
}
