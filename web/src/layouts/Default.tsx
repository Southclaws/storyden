import { PropsWithChildren, ReactNode } from "react";

import { Navigation } from "src/components/site/Navigation/Navigation";

import { MotdBanner } from "@/components/site/MotdBanner/MotdBanner";
import { getSettings } from "@/lib/settings/settings-server";
import { Box, Flex, styled } from "@/styled-system/jsx";

type Props = {
  contextpane: ReactNode;
};

export async function Default({
  contextpane,
  children,
}: PropsWithChildren<Props>) {
  const settings = await getSettings();

  return (
    <Flex
      className="sd-layout sd-layout--dashboard"
      minHeight="dvh"
      width="full"
      flexDirection="row"
      backgroundColor="bg.site"
      vaul-drawer-wrapper=""
    >
      <Navigation contextpane={contextpane}>
        <styled.main
          className="sd-screen sd-screen--dashboard"
          containerType="inline-size"
          width="full"
          height="full"
          minW="0"
        >
          <MotdBanner motd={settings.motd} />
          {children}
          <Box height="24"></Box>
        </styled.main>
      </Navigation>
    </Flex>
  );
}
