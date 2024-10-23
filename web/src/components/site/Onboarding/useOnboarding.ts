import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";

import { OnboardingStatus } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { useSettings } from "@/lib/settings/settings-client";

export function useOnboarding() {
  const { settings } = useSettings();
  const session = useSession();
  const pathName = usePathname();
  const [localStatus, setlocalStatus] = useState<string | null>(null);

  const isComposingNewThread = pathName === "/new";
  // const isAdmin = session?.admin ?? false;

  // NOTE: local onboarding status value takes priority.
  useEffect(() => {
    setlocalStatus(localStorage.getItem("onboarding-status"));
  }, [settings]);

  function onFinish() {
    localStorage.setItem("onboarding-status", "complete");
    setlocalStatus("complete");
  }

  const onboardingStatus: OnboardingStatus =
    (localStatus as OnboardingStatus) ??
    settings?.onboarding_status ??
    "complete";

  // Rules: If there's no session and the onboarding has not started, show
  // the onboarding to ANY user. But, once the first account is created, only
  // show the onboarding flow to the newly authenticated admin account. Once the
  // first stage is done there's no point showing the onboarding flow to guests.
  const isOnboardingAccount = session
    ? onboardingStatus !== "requires_first_account"
    : onboardingStatus === "requires_first_account";

  const showOnboarding =
    isOnboarding(onboardingStatus) &&
    !isComposingNewThread &&
    isOnboardingAccount;

  return {
    showOnboarding,
    onboardingStatus,
    onFinish,
  };
}

function isOnboarding(status?: OnboardingStatus) {
  switch (status) {
    // NOTE: explicit exhaustivity here because we want to default to content
    // if something went wrong with the info data. Last resort: show content!
    case "requires_first_account":
    case "requires_category":
    case "requires_more_accounts":
    case "requires_first_post":
      return true;

    case "complete":
    default:
      return false;
  }
}
