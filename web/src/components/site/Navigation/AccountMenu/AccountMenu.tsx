import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { Account } from "@/api/openapi-schema";
import { MemberAvatar } from "@/components/member/MemberBadge/MemberAvatar";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import * as Menu from "@/components/ui/menu";
import { hasPermission } from "@/utils/permissions";

import { AdminMenuItem } from "../Anchors/Admin";
import { DraftsMenuItem } from "../Anchors/Drafts";
import { LogoutMenuItem } from "../Anchors/Logout";
import { ProfileMenuItem } from "../Anchors/Profile";
import { QueueMenuItem } from "../Anchors/Queue";
import { SettingsMenuItem } from "../Anchors/Settings";

type Props = {
  account: Account;
};

export function AccountMenu({ account }: Props) {
  const isAdmin = hasPermission(account, "ADMINISTRATOR");

  return (
    <Menu.Root>
      <Menu.Trigger cursor="pointer">
        <MemberAvatar profile={account} size="md" />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="account">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <MemberBadge
                  profile={account}
                  as="link"
                  size="md"
                  name="full-vertical"
                />
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <ProfileMenuItem handle={account.handle} />
              <SettingsMenuItem />
              {isAdmin && <AdminMenuItem />}
            </Menu.ItemGroup>

            <Menu.ItemGroup id="content">
              <DraftsMenuItem />
              <QueueMenuItem />
            </Menu.ItemGroup>

            <Menu.Separator />

            <Menu.ItemGroup id="logout">
              <LogoutMenuItem />
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
