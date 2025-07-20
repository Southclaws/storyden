"use client";

import Link from "next/link";

import { AccountRoleList, ProfileReference } from "@/api/openapi-schema";
import { WEB_ADDRESS } from "@/config";
import { css } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

import { MemberOptionsMenu } from "../MemberOptions/MemberOptionsMenu";

import { MemberIdent } from "./MemberIdent";

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  showRoles?: "hidden" | "badge" | "all";
  roles?: AccountRoleList;
  avatar?: "hidden" | "visible";

  // NOTE: If you don't need either of these, just render a <MemberIdent />.
  as?: "menu" | "link";
};

const identContainerStyles = css({
  maxW: "full",
  flexShrink: "0",
});

export function MemberBadge({
  profile,
  size = "md",
  name = "hidden",
  avatar = "visible",
  showRoles = "hidden",
  roles,
  as = "menu",
}: Props) {
  const permalink = `${WEB_ADDRESS}/m/${profile.handle}`;

  if (as === "menu") {
    return (
      <HStack className={identContainerStyles}>
        <MemberOptionsMenu profile={profile}>
          <MemberIdent
            profile={profile}
            size={size}
            name={name}
            showRoles={showRoles}
            roles={roles}
            avatar={avatar}
          />
        </MemberOptionsMenu>
      </HStack>
    );
  }

  return (
    <Link className={identContainerStyles} href={permalink}>
      <MemberIdent
        profile={profile}
        size={size}
        name={name}
        showRoles={showRoles}
        roles={roles}
        avatar={avatar}
      />
    </Link>
  );
}
