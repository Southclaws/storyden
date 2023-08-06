import { HStack, VStack } from "@chakra-ui/react";

import { useSession } from "src/auth";
import { Logout } from "src/components/Action/Action";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";

export function Authbar() {
  const account = useSession();

  if (!account) return null;

  return (
    <HStack alignItems="center">
      <VStack alignItems="start">
        <ProfileReference profileReference={account} size="lg" />
        <HStack>
          <Logout />
          {/* NOTE: no settings yet lol */}
          {/* <SettingsMenu /> */}
        </HStack>
      </VStack>
    </HStack>
  );
}
