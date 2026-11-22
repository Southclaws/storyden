"use client";

import Link from "next/link";

import { ProfileReference } from "@/api/openapi-schema";
import { WEB_ADDRESS } from "@/config";
import { css, cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

import { MemberOptionsMenu } from "../MemberOptions/MemberOptionsMenu";

import { MemberIdent } from "./MemberIdent";
import type { MemberIdentity } from "./types";

export type Props = {
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  showRoles?: "hidden" | "badge" | "all";
  avatar?: "hidden" | "visible";

  // NOTE: If you don't need either of these, just render a <MemberIdent />.
} & (
  | { profile: ProfileReference; as?: "menu" }
  | { profile: MemberIdentity; as: "link" }
);

const identContainerStyles = css({
  maxW: "full",
  minW: "0",
  flexShrink: "0",
  flex: "1",
});

export function MemberBadge(props: Props) {
  const {
    profile,
    size = "md",
    name = "hidden",
    avatar = "visible",
    showRoles = "hidden",
  } = props;
  const permalink = `${WEB_ADDRESS}/m/${profile.handle}`;

  if (props.as !== "link") {
    return (
      <HStack
        className={cx("member-badge__menu-container", identContainerStyles)}
      >
        <MemberOptionsMenu profile={props.profile}>
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
      prefetch={false}
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
