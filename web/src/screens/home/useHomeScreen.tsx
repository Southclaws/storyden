import { useEffect, useState } from "react";

import { useInfoProvider } from "src/api/InfoProvider/useInfoProvider";
import { OnboardingStatus } from "src/api/openapi/schemas";

export function useHomeScreen() {
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

  return {
    onboardingStatus,
    onFinish,
  };
}
