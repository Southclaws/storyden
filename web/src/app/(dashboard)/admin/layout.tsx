"use client";

import { usePathname, useRouter } from "next/navigation";
import { type PropsWithChildren, useEffect } from "react";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { Permission } from "@/api/openapi-schema";
import { BackAction } from "@/components/site/Action/Back";
import { UnreadyBanner } from "@/components/site/Unready";
import {
  SectionNavigation,
  isSectionActive,
} from "@/components/ui/section-navigation";
import { useSettings } from "@/lib/settings/settings-client";
import {
  DEFAULT_ADMIN_SECTION,
  getAdminSectionGroups,
} from "@/screens/admin/adminSections";
import { hasPermission } from "@/utils/permissions";

export default function Layout({ children }: PropsWithChildren) {
  const pathname = usePathname();
  const router = useRouter();
  const account = useAccountGet();
  const settings = useSettings();
  const ready = settings.ready && !account.isLoading;
  const authorised = hasPermission(
    account.data,
    Permission.ADMINISTRATOR,
    Permission.MANAGE_SETTINGS,
  );
  const groups =
    settings.ready && authorised
      ? getAdminSectionGroups(settings.settings, account.data)
      : [];
  const sectionAvailable =
    pathname === "/admin" ||
    groups.some((group) =>
      group.items.some(({ href }) => isSectionActive(href, pathname)),
    );

  useEffect(() => {
    if (ready && authorised && !sectionAvailable) {
      router.replace(`/admin/${DEFAULT_ADMIN_SECTION}`);
    }
  }, [authorised, ready, router, sectionAvailable]);

  if (!ready) {
    return <UnreadyBanner error={settings.error ?? account.error} />;
  }

  if (!authorised) {
    return (
      <UnreadyBanner error="Not authorised to view the system configuration page." />
    );
  }

  return (
    <div className="section-navigation-layout">
      <div className="section-navigation-layout__mobile-navigation">
        <BackAction flexShrink="0" />
        <SectionNavigation
          className="section-navigation-layout__navigation"
          groups={groups}
          label="Admin sections"
        />
      </div>
      <div className="section-navigation-layout__content">{children}</div>
    </div>
  );
}
