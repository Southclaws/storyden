"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import * as Tabs from "@/components/ui/tabs";
import { Settings } from "@/lib/settings/settings";

import { MemberAuthenticationSettingsScreen } from "./MemberAuthenticationSettingsScreen";
import { MemberEmailSettingsScreen } from "./MemberEmailSettingsScreen";

const DEFAULT_TAB = "authentication";

type Props = {
  initialSettings: Settings;
};

export function SettingsScreen({ initialSettings }: Props) {
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

  const emailEnabled = initialSettings.capabilities.includes("email_client");

  return (
    <Tabs.Root
      width="full"
      variant="enclosed"
      defaultValue={DEFAULT_TAB}
      value={tab}
      onValueChange={handleTabChange}
    >
      <Tabs.List>
        <Tabs.Trigger value="authentication">Authentication</Tabs.Trigger>
        {emailEnabled && <Tabs.Trigger value="email">Email</Tabs.Trigger>}
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="authentication">
        <MemberAuthenticationSettingsScreen />
      </Tabs.Content>

      {emailEnabled && (
        <Tabs.Content value="email">
          <MemberEmailSettingsScreen />
        </Tabs.Content>
      )}
    </Tabs.Root>
  );
}
