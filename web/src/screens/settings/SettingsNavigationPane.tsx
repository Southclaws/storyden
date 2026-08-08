"use client";

import { usePathname } from "next/navigation";
import { type ReactNode } from "react";

import { useAccountGet } from "@/api/openapi-client/accounts";
import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { InboxIcon } from "@/components/ui/icons/Inbox";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import { LinkIcon } from "@/components/ui/icons/Link";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { useSettings } from "@/lib/settings/settings-client";

import {
  type SettingsSectionID,
  getSettingsSections,
} from "./settingsSections";

const sectionIcons: Record<SettingsSectionID, ReactNode> = {
  interface: <LayoutIcon />,
  authentication: <SettingsIcon />,
  email: <InboxIcon />,
  access_keys: <ToolIcon />,
  oauth: <LinkIcon />,
};

export function SettingsNavigationPane() {
  const pathname = usePathname();
  const account = useAccountGet();
  const settings = useSettings();
  const sections = settings.ready
    ? getSettingsSections(settings.settings, account.data)
    : [];

  return (
    <SidebarNavigationPane label="Settings navigation" title="Settings">
      <SidebarNavigationSection>
        {sections.map((section) => (
          <SidebarNavigationLink
            key={section.id}
            href={section.href}
            label={section.label}
            description={section.description}
            icon={sectionIcons[section.id]}
            active={section.href === pathname}
          />
        ))}
      </SidebarNavigationSection>
    </SidebarNavigationPane>
  );
}
