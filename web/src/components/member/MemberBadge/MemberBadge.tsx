"use client";

import { Portal } from "@ark-ui/react";
import Link from "next/link";

import { handle } from "@/api/client";
import { adminAccountBanCreate } from "@/api/openapi-client/admin";
import { ProfileReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import * as Menu from "@/components/ui/menu";
import { WEB_ADDRESS } from "@/config";
import { HStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { MemberIdent } from "./MemberIdent";

export type Props = {
  profile: ProfileReference;
  size?: "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  roles?: "hidden" | "badge" | "all";
  avatar?: "hidden" | "visible";

  // NOTE: If you don't need either of these, just render a <MemberIdent />.
  as?: "menu" | "link";
};

export function useMemberBadge(profile: ProfileReference) {
  const session = useSession();

  const isSelf = session?.id === profile.id;

  const isSuspendEnabled =
    !isSelf && hasPermission(session, "MANAGE_SUSPENSIONS");

  const isRoleChangeEnabled = hasPermission(session, "MANAGE_ROLES");

  const permalink = `${WEB_ADDRESS}/m/${profile.handle}`;

  // TODO: Add member suspension to API
  // const isMemberSuspended = profile.deletedAt !== null;

  async function handleSuspend() {
    handle(async () => {
      await adminAccountBanCreate(profile.handle);
    });
  }

  return {
    isSuspendEnabled,
    isRoleChangeEnabled,
    permalink,
    handlers: { handleSuspend },
  };
}

export function MemberBadge({
  profile,
  size = "md",
  name = "hidden",
  avatar = "visible",
  roles = "hidden",
  as = "menu",
}: Props) {
  const { isSuspendEnabled, permalink, handlers } = useMemberBadge(profile);

  if (as === "menu") {
    return (
      <HStack w="min" className="feed-item-byline-menu">
        <Menu.Root
          lazyMount
          positioning={{
            strategy: "absolute",
            placement: "bottom-start",
            overflowPadding: 16,
          }}
        >
          <Menu.Trigger
            cursor="pointer"
            title={`${profile.name} @${profile.handle}`}
          >
            <MemberIdent
              profile={profile}
              size={size}
              name={name}
              roles={roles}
              avatar={avatar}
            />
          </Menu.Trigger>

          <Portal>
            <Menu.Positioner>
              <Menu.Content>
                <Menu.ItemGroup>
                  <Menu.ItemGroupLabel>
                    <styled.p fontWeight="bold">{profile.name}</styled.p>
                    <styled.p fontWeight="normal" color="fg.muted">
                      @{profile.handle}
                    </styled.p>
                  </Menu.ItemGroupLabel>
                </Menu.ItemGroup>

                <Menu.Separator />

                <Menu.ItemGroup>
                  <Link href={permalink}>
                    <Menu.Item value="copy">View profile</Menu.Item>
                  </Link>
                </Menu.ItemGroup>

                {isSuspendEnabled && (
                  <Menu.ItemGroup>
                    <Menu.ItemGroupLabel>Admin</Menu.ItemGroupLabel>

                    <Menu.Item
                      value="suspend"
                      onClick={handlers.handleSuspend}
                      colorPalette="red"
                    >
                      Suspend member
                    </Menu.Item>
                  </Menu.ItemGroup>
                )}
              </Menu.Content>
            </Menu.Positioner>
          </Portal>
        </Menu.Root>
      </HStack>
    );
  }

  return (
    <Link className="feed-item-byline-basic" href={permalink}>
      <MemberIdent
        profile={profile}
        size={size}
        name={name}
        roles={roles}
        avatar={avatar}
      />
    </Link>
  );
}
