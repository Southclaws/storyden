import { useSession } from "src/auth";
import { ComposeAction } from "src/components/site/Action/Compose";
import { HomeAction } from "src/components/site/Action/Home";
import { LoginAction } from "src/components/site/Action/Login";
import { NotificationsAction } from "src/components/site/Action/Notifications";

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
