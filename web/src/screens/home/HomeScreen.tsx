"use client";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { Onboarding } from "src/components/site/Onboarding/Onboarding";

import { Content } from "./Content";
import { useHomeScreen } from "./useHomeScreen";

export function HomeScreen() {
  const { onboardingStatus, onFinish } = useHomeScreen();

  const showOnboarding = isOnboarding(onboardingStatus) && onboardingStatus;

  return (
    <>
      <Content showEmptyState={!showOnboarding} />

      {showOnboarding && (
        <Onboarding status={onboardingStatus} onFinish={onFinish} />
      )}
    </>
  );
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
