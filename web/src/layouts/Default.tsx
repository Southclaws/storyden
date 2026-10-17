import { PropsWithChildren, ReactNode } from "react";

import { FreshDenDoctor } from "@/components/site/DenDoctor/FreshDenDoctor";
import { MotdBanner } from "@/components/site/MotdBanner/MotdBanner";
import { Navigation } from "@/components/site/Navigation/Navigation";
import { API_ADDRESS, WEB_ADDRESS } from "@/config";
import { diagnoseInstallation } from "@/lib/den-doctor/diagnostics";
import { getSettingsResult } from "@/lib/settings/settings-server";
import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  sidebar: ReactNode;
};

export async function Default({ sidebar, children }: PropsWithChildren<Props>) {
  const result = await getSettingsResult();
  const settings = result.ready ? result.settings : undefined;

  return (
    <Flex
      data-sd-layout="default"
      minHeight="dvh"
      width="full"
      flexDirection="row"
      backgroundColor="background.canvas"
      vaul-drawer-wrapper=""
    >
      <Navigation sidebar={sidebar}>
        <styled.main
          data-sd-region="main"
          containerType="inline-size"
          width="full"
          height="full"
          minW="0"
        >
          <MotdBanner motd={settings?.motd} />
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
          {/* <Box height="24"></Box> */}
        </styled.main>
      </Navigation>
    </Flex>
  );
}
