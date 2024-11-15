"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";

import * as Tabs from "@/components/ui/tabs";

import { AuthenticationSettingsScreen } from "./AuthenticationSettingsScreen";
import { BrandSettingsScreen } from "./BrandSettingsScreen";

const DEFAULT_TAB = "brand";

export function AdminScreen() {
  const [tab, setTab] = useQueryState("tab", {
    defaultValue: DEFAULT_TAB,
  });

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
      lazyMount
    >
      <Tabs.List>
        <Tabs.Trigger value="brand">Brand</Tabs.Trigger>
        {/* <Tabs.Trigger value="authentication">Authentication</Tabs.Trigger> */}
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="brand">
        <BrandSettingsScreen />
      </Tabs.Content>

      {/* <Tabs.Content value="authentication">
        <AuthenticationSettingsScreen />
      </Tabs.Content> */}
    </Tabs.Root>
  );
}
