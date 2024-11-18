"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import * as Tabs from "@/components/ui/tabs";

import { AuthenticationSettingsScreen } from "./AuthenticationSettingsScreen";
import { BrandSettingsScreen } from "./BrandSettingsScreen";

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
        <Tabs.Trigger value="authentication">Authentication</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="brand">
        <BrandSettingsScreen />
      </Tabs.Content>

      <Tabs.Content value="authentication">
        <AuthenticationSettingsScreen />
      </Tabs.Content>
    </Tabs.Root>
  );
}
