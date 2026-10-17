import { PropsWithChildren } from "react";

import { FreshDenDoctor } from "@/components/site/DenDoctor/FreshDenDoctor";
import { LinkButton } from "@/components/ui/link-button";
import { API_ADDRESS, WEB_ADDRESS } from "@/config";
import { diagnoseInstallation } from "@/lib/den-doctor/diagnostics";
import { getSettingsResult } from "@/lib/settings/settings-server";
import { VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  const result = await getSettingsResult();
  const settings = result.ready ? result.settings : undefined;

  return (
    <VStack w="full">
      {(!settings ||
        settings.onboarding_status === "requires_first_account") && (
        <FreshDenDoctor
          settings={settings}
          hasAddressIssues={
            diagnoseInstallation(API_ADDRESS, WEB_ADDRESS).length > 0
          }
        />
      )}
      {children}

      <LinkButton variant="subtle" href="/login">
        Sign in
      </LinkButton>
    </VStack>
  );
}
