"use client";

import dynamic from "next/dynamic";

import { WEB_ADDRESS } from "@/config";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { useHydrated } from "@/utils/useHydrated";

const DenDoctor = dynamic(
  () => import("./DenDoctor").then((module) => module.DenDoctor),
  { ssr: false },
);

type Props = {
  settings?: Settings;
  hasAddressIssues: boolean;
};

export function FreshDenDoctor({ settings, hasAddressIssues }: Props) {
  const hydrated = useHydrated();
  const current = useSettings(settings, true);
  const status =
    current.settings?.onboarding_status ?? settings?.onboarding_status;
  const origin = hydrated ? window.location.origin : undefined;

  if (status !== "requires_first_account") {
    return null;
  }

  if (
    !hasAddressIssues &&
    (!origin || origin === new URL(WEB_ADDRESS).origin)
  ) {
    return null;
  }

  return <DenDoctor browserOrigin={origin} />;
}
