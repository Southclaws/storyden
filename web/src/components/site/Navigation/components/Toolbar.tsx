import { useSession } from "src/auth";
import { ComposeAction } from "src/components/site/Navigation/Anchors/Compose";
import {
  LoginAction,
  RegisterAction,
} from "src/components/site/Navigation/Anchors/Login";

import { ProfilePill } from "../../ProfilePill/ProfilePill";

import { HStack } from "@/styled-system/jsx";

export function Toolbar() {
  const account = useSession();
  return (
    <HStack gap="2" alignItems="center">
      {account ? (
        <>
          <HStack>
            <ComposeAction />
            <ProfilePill profileReference={account} />
          </HStack>
        </>
      ) : (
        <>
          <HStack>
            <RegisterAction w="full" />
            <LoginAction flexShrink={0} />
          </HStack>
        </>
      )}
    </HStack>
  );
}
