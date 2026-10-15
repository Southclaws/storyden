"use client";

import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import dynamic from "next/dynamic";

import { Account } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import {
  MemberIdent,
  type Props as MemberIdentProps,
} from "@/components/member/MemberBadge/MemberIdent";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { UnreadyBanner } from "@/components/site/Unready";
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

import { InvitationID, InvitationMenuItem } from "./InvitationMenuItem";

const InvitationScreen = dynamic(
  () => import("./InvitationScreen").then((module) => module.InvitationScreen),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

type Props = {
  account: Account;
  size?: "sm" | "md";
  name?: MemberIdentProps["name"];
};

export function AccountMenu({ account, size = "md", name = "hidden" }: Props) {
  const isAdmin = hasPermission(account, "ADMINISTRATOR");
  const canCreateInvitations = hasPermission(account, "CREATE_INVITATION");
  const invitationDisclosure = useDisclosure();
  function handleSelect(details: MenuSelectionDetails) {
    if (details.value === InvitationID) {
      invitationDisclosure.onOpen();
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
        <Menu.Trigger
          cursor="pointer"
          maxW="full"
          minW="0"
          aria-label="Account menu"
        >
          <MemberIdent profile={account} size={size} name={name} />
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

      <ModalDrawer
        title="Invite someone"
        isOpen={invitationDisclosure.isOpen}
        onClose={invitationDisclosure.onClose}
      >
        <InvitationScreen onClose={invitationDisclosure.onClose} />
      </ModalDrawer>
    </>
  );
}
