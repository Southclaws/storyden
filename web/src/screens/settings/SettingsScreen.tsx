"use client";

import { TabsValueChangeDetails } from "@ark-ui/react";
import { useQueryState } from "nuqs";
import { useEffect } from "react";

import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import * as Tabs from "@/components/ui/tabs";
import { Settings } from "@/lib/settings/settings";
import { hasPermission } from "@/utils/permissions";

import { MemberAccessKeysSettingsScreen } from "./MemberAccessKeysSettingsScreen";
import { MemberAuthenticationSettingsScreen } from "./MemberAuthenticationSettingsScreen";
import { MemberEmailSettingsScreen } from "./MemberEmailSettingsScreen";

const DEFAULT_TAB = "authentication";

type Props = {
  initialSettings: Settings;
};

export function SettingsScreen({ initialSettings }: Props) {
  const session = useSession();
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

  const accessKeysEnabled = hasPermission(
    session,
    Permission.USE_PERSONAL_ACCESS_KEYS,
  );

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
        {accessKeysEnabled && (
          <Tabs.Trigger value="access_keys">Access keys</Tabs.Trigger>
        )}
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

      {accessKeysEnabled && (
        <Tabs.Content value="access_keys">
          <MemberAccessKeysSettingsScreen />
        </Tabs.Content>
      )}
    </Tabs.Root>
  );
}
