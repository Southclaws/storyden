"use client";

import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import { InstanceCapability } from "@/api/openapi-schema";
import { UnreadyBanner } from "@/components/site/Unready";
import { PageLayout } from "@/components/ui/page-layout";
import { hasCapability } from "@/lib/settings/capabilities";
import { useSettings } from "@/lib/settings/settings-client";
import { RobotsTabs } from "@/screens/robots/RobotsTabs";

export default function Layout({ children }: PropsWithChildren) {
  const settings = useSettings();

  if (!settings.ready) {
    return (
      <PageLayout width="wide">
        <UnreadyBanner error={settings.error} />
      </PageLayout>
    );
  }

  if (
    !hasCapability(InstanceCapability.robots, settings.settings.capabilities)
  ) {
    notFound();
  }

  return (
    <PageLayout width="wide">
      <RobotsTabs>{children}</RobotsTabs>
    </PageLayout>
  );
}
