import { useSession } from "src/auth";
import { ComposeAction } from "src/components/site/Navigation/Anchors/Compose";
import { HomeAction } from "src/components/site/Navigation/Anchors/Home";
import { LoginAction } from "src/components/site/Navigation/Anchors/Login";
import { NotificationsAction } from "src/components/site/Navigation/Anchors/Notifications";

import { HStack } from "@/styled-system/jsx";

export function Toolbar() {
  const account = useSession();
  return (
    <HStack gap="2" pb="2">
      {account ? (
        <>
          <HomeAction />
          <NotificationsAction />
          <ComposeAction />
        </>
      ) : (
        <>
          <HomeAction />
          <LoginAction />
        </>
      )}
    </HStack>
  );
}
