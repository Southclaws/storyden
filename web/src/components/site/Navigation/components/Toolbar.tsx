import { useSession } from "src/auth";
import { AdminAction } from "src/components/site/Navigation/Anchors/Admin";
import {
  LoginAction,
  RegisterAction,
} from "src/components/site/Navigation/Anchors/Login";
import { LogoutAction } from "src/components/site/Navigation/Anchors/Logout";
import { SettingsAction } from "src/components/site/Navigation/Anchors/Settings";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";

import { HStack } from "@/styled-system/jsx";

export function Toolbar() {
  const account = useSession();
  return (
    <HStack w="full" gap="2" alignItems="center">
      {account ? (
        <HStack w="full" justify="space-between">
          <HStack>
            {/* TODO: Put some of this in a menu */}
            <SettingsAction />
            {account.admin && <AdminAction />}
            <LogoutAction />
          </HStack>

          <ProfilePill
            profileReference={account}
            size="lg"
            showHandle={false}
          />
        </HStack>
      ) : (
        <HStack>
          <RegisterAction w="full" />
          <LoginAction flexShrink={0} />
        </HStack>
      )}
    </HStack>
  );
}
