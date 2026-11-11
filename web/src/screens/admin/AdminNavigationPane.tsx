"use client";

import { usePathname } from "next/navigation";
import { type ReactNode } from "react";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { useGetVersion } from "@/api/openapi-client/misc";
import { Permission } from "@/api/openapi-schema";
import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { ArchiveIcon } from "@/components/ui/icons/Archive";
import { BanIcon } from "@/components/ui/icons/BanIcon";
import { BiometricIcon } from "@/components/ui/icons/Biometric";
import { CollectionIcon } from "@/components/ui/icons/Collection";
import { ColourPipetteIcon } from "@/components/ui/icons/Colour";
import { InboxIcon } from "@/components/ui/icons/Inbox";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import { LinkIcon } from "@/components/ui/icons/Link";
import { MediaIcon } from "@/components/ui/icons/Media";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { SearchIcon } from "@/components/ui/icons/Search";
import { SystemIcon } from "@/components/ui/icons/System";
import { ToolIcon } from "@/components/ui/icons/Tool";
import { Text } from "@/components/ui/text";
import { useSettings } from "@/lib/settings/settings-client";
import { hasPermission } from "@/utils/permissions";

import { type AdminSectionID, getAdminSectionGroups } from "./adminSections";

const sectionIcons: Record<AdminSectionID, ReactNode> = {
  brand: <ColourPipetteIcon />,
  interface: <LayoutIcon />,
  theme: <ColourPipetteIcon />,
  moderation: <BanIcon />,
  media: <MediaIcon />,
  system: <SystemIcon />,
  discovery: <SearchIcon />,
  audit_log: <ArchiveIcon />,
  email: <InboxIcon />,
  authentication: <BiometricIcon />,
  access_keys: <ToolIcon />,
  oauth: <LinkIcon />,
  plugins: <CollectionIcon />,
  robots: <RobotIcon />,
};

export function AdminNavigationPane() {
  const pathname = usePathname();
  const account = useAccountGet();
  const settings = useSettings();
  const authorised = hasPermission(
    account.data,
    Permission.ADMINISTRATOR,
    Permission.MANAGE_SETTINGS,
  );
  const version = useGetVersion({
    swr: { enabled: settings.ready && !account.isLoading && authorised },
  });
  let versionLabel = "Loading Storyden version…";

  if (version.data) {
    versionLabel = `Storyden ${version.data}`;
  } else if (version.error) {
    versionLabel = "Storyden version unavailable";
  }

  if (!settings.ready || account.isLoading || !authorised) {
    return null;
  }

  const groups = getAdminSectionGroups(settings.settings, account.data);

  return (
    <SidebarNavigationPane label="Admin navigation" title="Admin">
      {groups.map((group) => (
        <SidebarNavigationSection key={group.label} label={group.label}>
          {group.items.map((section) => (
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
      ))}
      <Text
        as="footer"
        variant="metadata"
        mt="auto"
        px="2"
        py="2"
        flexShrink="0"
        overflowWrap="anywhere"
      >
        {versionLabel}
      </Text>
    </SidebarNavigationPane>
  );
}
