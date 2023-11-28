import { useSession } from "src/auth";
import { AdminAction } from "src/components/site/Action/Admin";
import { LogoutAction } from "src/components/site/Action/Logout";
import { SettingsAction } from "src/components/site/Action/Settings";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";

import { HStack, VStack } from "@/styled-system/jsx";

export function Authbar() {
  const account = useSession();

  if (!account) return null;

  return (
    <HStack alignItems="center">
      <VStack alignItems="start">
        <ProfilePill profileReference={account} size="lg" />
        <HStack>
          <LogoutAction />
          <SettingsAction />
          {account.admin && <AdminAction />}
        </HStack>
      </VStack>
    </HStack>
  );
}
