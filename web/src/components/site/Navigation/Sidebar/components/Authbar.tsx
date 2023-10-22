import { HStack, VStack } from "@chakra-ui/react";

import { useSession } from "src/auth";
import { Admin, Logout, Settings } from "src/components/site/Action/Action";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";

export function Authbar() {
  const account = useSession();

  if (!account) return null;

  return (
    <HStack alignItems="center">
      <VStack alignItems="start">
        <ProfilePill profileReference={account} size="lg" />
        <HStack>
          <Logout />
          <Settings />
          {account.admin && <Admin />}
        </HStack>
      </VStack>
    </HStack>
  );
}
