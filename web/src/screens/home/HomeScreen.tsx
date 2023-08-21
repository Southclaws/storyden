"use client";

import { useInfoProvider } from "src/api/InfoProvider/useInfoProvider";

import { Content } from "./Content";
import { Onboarding } from "./Onboarding";

export function HomeScreen() {
  const info = useInfoProvider();

  console.log(info);

  switch (info?.onboarding_status) {
    // NOTE: explicit exhaustivity here because we want to default to content
    // if something went wrong with the info data. Last resort: show content!
    case "requires_first_account":
    case "requires_category":
    case "requires_more_accounts":
    case "requires_first_post":
      return <Onboarding status={info.onboarding_status} />;

    case "complete":
    default:
      return <Content />;
  }
}
