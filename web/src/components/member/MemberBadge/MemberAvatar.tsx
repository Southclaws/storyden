import Image from "next/image";

import {
  AssetVariant,
  Identifier,
  ProfileReference,
} from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";

import { MemberAvatarEditable } from "./MemberAvatarEditable";
import type { MemberIdentity } from "./types";

const avatarStyles = css({
  borderRadius: "full",
});

export type Props = {
  size?: "xs" | "sm" | "md" | "lg";
  alt?: string;
} & (
  | { profile: MemberIdentity; editable?: false }
  | { profile: ProfileReference; editable: boolean }
);

export function MemberAvatar(props: Props) {
  const { profile, size, alt } = props;
  const avatarURL = getAvatarURL(
    profile.handle,
    size === "xs" || size === "sm"
      ? "thumb_s"
      : size === "md"
        ? "thumb_m"
        : "thumb_l",
  );

  const { width, height } = avatarSize(size);

  return (
    <Box position="relative" flexShrink="0">
      {props.editable && <MemberAvatarEditable profile={props.profile} />}
      <Image
        className={avatarStyles}
        src={avatarURL}
        alt={alt ?? `Avatar of ${profile.name || profile.handle}`}
        width={width}
        height={height}
      />
    </Box>
  );
}

export function getAvatarURL(id: Identifier, variant?: AssetVariant): string {
  const [path] = [`/accounts/${id}/avatar`] as const;

  return `${API_ADDRESS}/api${path}${variant ? `?variant=${variant}` : ""}`;
}

export function avatarSize(size: Props["size"]) {
  switch (size) {
    case "xs":
      return { width: 16, height: 16 };
    case "sm":
      return { width: 24, height: 24 };
    case "md":
      return { width: 36, height: 36 };
    case "lg":
    default:
      return { width: 100, height: 100 };
  }
}
