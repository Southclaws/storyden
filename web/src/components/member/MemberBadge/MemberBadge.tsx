"use client";

import Link from "next/link";

import { AccountRoleList, ProfileReference } from "@/api/openapi-schema";
import { WEB_ADDRESS } from "@/config";
import { css, cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

import { MemberOptionsMenu } from "../MemberOptions/MemberOptionsMenu";

import { MemberIdent } from "./MemberIdent";

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  showRoles?: "hidden" | "badge" | "all";
  avatar?: "hidden" | "visible";

  // NOTE: If you don't need either of these, just render a <MemberIdent />.
  as?: "menu" | "link";
};

const identContainerStyles = css({
  maxW: "full",
  minW: "0",
  flexShrink: "0",
  flex: "1",
});

export function MemberBadge({
  profile,
  size = "md",
  name = "hidden",
  avatar = "visible",
  showRoles = "hidden",
  as = "menu",
}: Props) {
  const permalink = `${WEB_ADDRESS}/m/${profile.handle}`;

  if (as === "menu") {
    return (
      <HStack
        className={cx("member-badge__menu-container", identContainerStyles)}
      >
        <MemberOptionsMenu profile={profile}>
          <MemberIdent
            profile={profile}
            size={size}
            name={name}
            showRoles={showRoles}
            avatar={avatar}
          />
        </MemberOptionsMenu>
      </HStack>
    );
  }

  return (
    <Link
      className={cx("member-badge__link-container", identContainerStyles)}
      href={permalink}
    >
      <MemberIdent
        profile={profile}
        size={size}
        name={name}
        showRoles={showRoles}
        avatar={avatar}
      />
    </Link>
  );
}
