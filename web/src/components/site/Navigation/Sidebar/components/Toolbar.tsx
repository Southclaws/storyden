import { useSession } from "src/auth";
import { ComposeAction } from "src/components/site/Navigation/Anchors/Compose";
import { HomeAction } from "src/components/site/Navigation/Anchors/Home";
import {
  LoginAction,
  RegisterAction,
} from "src/components/site/Navigation/Anchors/Login";
import { NotificationsAction } from "src/components/site/Navigation/Anchors/Notifications";

import { HStack } from "@/styled-system/jsx";

export function Toolbar() {
  const account = useSession();
  return (
    <HStack gap="2" pb="2" w="full">
      {account ? (
        <>
          <HomeAction />
          <NotificationsAction />
          <ComposeAction />
        </>
      ) : (
        <>
          <HStack w="full">
            <RegisterAction w="full" />
            <LoginAction flexShrink={0} />
          </HStack>
        </>
      )}
    </HStack>
  );
}
