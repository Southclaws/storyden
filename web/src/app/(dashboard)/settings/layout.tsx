"use client";

import { usePathname, useRouter } from "next/navigation";
import { type PropsWithChildren, useEffect } from "react";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { SectionNavigation } from "@/components/ui/section-navigation";
import { useSettings } from "@/lib/settings/settings-client";
import {
  DEFAULT_SETTINGS_SECTION,
  getSettingsSections,
} from "@/screens/settings/settingsSections";

export default function Layout({ children }: PropsWithChildren) {
  const pathname = usePathname();
  const router = useRouter();
  const account = useAccountGet();
  const settings = useSettings();
  const ready = settings.ready && !account.isLoading;
  const sections = settings.ready
    ? getSettingsSections(settings.settings, account.data)
    : [];
  const sectionAvailable =
    pathname === "/settings" || sections.some(({ href }) => href === pathname);

  useEffect(() => {
    if (ready && !sectionAvailable) {
      router.replace(`/settings/${DEFAULT_SETTINGS_SECTION}`);
    }
  }, [ready, router, sectionAvailable]);

  return (
    <div className="section-navigation-layout">
      {ready && (
        <SectionNavigation
          className="section-navigation-layout__navigation"
          groups={[{ items: sections }]}
          label="Settings sections"
        />
      )}
      <div className="section-navigation-layout__content">{children}</div>
    </div>
  );
}
