"use client";

import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { Account } from "@/api/openapi-schema";
import { MemberAvatar } from "@/components/member/MemberBadge/MemberAvatar";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import * as Menu from "@/components/ui/menu";
import { hasPermission } from "@/utils/permissions";
import { useDisclosure } from "@/utils/useDisclosure";

import { AdminMenuItem } from "../Anchors/Admin";
import { DraftsMenuItem } from "../Anchors/Drafts";
import { LogoutMenuItem } from "../Anchors/Logout";
import { ProfileMenuItem } from "../Anchors/Profile";
import { QueueMenuItem } from "../Anchors/Queue";
import { ReportsMenuItem } from "../Anchors/Reports";
import { SettingsMenuItem } from "../Anchors/Settings";

import {
  InvitationID,
  InvitationMenuItem,
  InvitationModal,
  useInvitation,
} from "./InvitationMenuItem";

type Props = {
  account: Account;
  size?: "sm" | "md";
};

export function AccountMenu({ account, size = "md" }: Props) {
  const isAdmin = hasPermission(account, "ADMINISTRATOR");
  const canCreateInvitations = hasPermission(account, "CREATE_INVITATION");
  const invitationDisclosure = useDisclosure();
  const invitation = useInvitation();

  async function handleInvitationSelect() {
    invitationDisclosure.onOpen();
    await invitation.createInvitation();
  }

  function handleSelect(details: MenuSelectionDetails) {
    if (details.value === InvitationID) {
      void handleInvitationSelect();
    }
  }

  return (
    <>
      <Menu.Root
        size="md"
        onSelect={handleSelect}
        positioning={{
          fitViewport: true,
          slide: true,
          placement: "bottom-end",
          shift: size === "md" ? 24 : 0,
        }}
      >
        <Menu.Trigger cursor="pointer" aria-label="Account menu">
          <MemberAvatar profile={account} size={size} />
        </Menu.Trigger>

        <Portal>
          <Menu.Positioner>
            <Menu.Content minW="72" userSelect="none">
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
                {canCreateInvitations && <InvitationMenuItem />}
              </Menu.ItemGroup>

              <Menu.ItemGroup id="content">
                <DraftsMenuItem />
                <QueueMenuItem />
                <ReportsMenuItem />
              </Menu.ItemGroup>

              <Menu.Separator />

              <Menu.ItemGroup id="logout">
                <LogoutMenuItem />
              </Menu.ItemGroup>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>

      <InvitationModal
        isOpen={invitationDisclosure.isOpen}
        invitation={invitation.invitation}
        onRetry={invitation.createInvitation}
        onClose={invitationDisclosure.onClose}
      />
    </>
  );
}
