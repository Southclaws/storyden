import { useEffect, useState } from "react";

import { useInfoProvider } from "src/api/InfoProvider/useInfoProvider";
import { OnboardingStatus } from "src/api/openapi/schemas";

export function useOnboarding() {
  const info = useInfoProvider();
  const [localStatus, setlocalStatus] = useState<string | null>(null);

  // NOTE: local onboarding status value takes priority.
  useEffect(() => {
    setlocalStatus(localStorage.getItem("onboarding-status"));
  }, [info]);

  function onFinish() {
    localStorage.setItem("onboarding-status", "complete");
    setlocalStatus("complete");
  }

  const onboardingStatus: OnboardingStatus =
    (localStatus as OnboardingStatus) ?? info?.onboarding_status ?? "complete";

  const showOnboarding = isOnboarding(onboardingStatus) && onboardingStatus;

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
