import { PropsWithChildren, ReactNode } from "react";

import { MotdBanner } from "@/components/site/MotdBanner/MotdBanner";
import { Navigation } from "@/components/site/Navigation/Navigation";
import { getSettings } from "@/lib/settings/settings-server";
import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  sidebar: ReactNode;
};

export async function Default({ sidebar, children }: PropsWithChildren<Props>) {
  const settings = await getSettings();

  return (
    <Flex
      minHeight="dvh"
      width="full"
      flexDirection="row"
      backgroundColor="background.canvas"
      vaul-drawer-wrapper=""
    >
      <Navigation sidebar={sidebar}>
        <styled.main
          containerType="inline-size"
          width="full"
          height="full"
          minW="0"
        >
          <MotdBanner motd={settings.motd} />
          {children}
          {/* <Box height="24"></Box> */}
        </styled.main>
      </Navigation>
    </Flex>
  );
}
