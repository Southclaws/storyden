import { HStack, VStack } from "@chakra-ui/react";
import { useSession } from "src/auth";
import { Logout, Settings } from "src/components/Action/Action";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";

export function Authbar() {
  const account = useSession();

  if (!account) return null;

  return (
    <HStack alignItems="center">
      <VStack alignItems="start">
        <HStack>
          <Logout />
          <Settings />
        </HStack>
        <ProfileReference handle={account.handle} />
      </VStack>
    </HStack>
  );
}
