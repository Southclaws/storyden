import Image from "next/image";

import { Identifier, ProfileReference } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";
import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";

import { MemberAvatarEditable } from "./MemberAvatarEditable";

const avatarStyles = css({
  borderRadius: "full",
});

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  editable?: boolean;
};

export function MemberAvatar({ profile, size, editable }: Props) {
  const avatarURL = getAvatarURL(profile.handle);

  const { width, height } = avatarSize(size);

  return (
    <Box position="relative" flexShrink="0">
      {editable && <MemberAvatarEditable profile={profile} />}
      <Image
        className={avatarStyles}
        src={avatarURL}
        alt={`${profile.handle}'s avatar`}
        width={width}
        height={height}
      />
    </Box>
  );
}

export function getAvatarURL(id: Identifier): string {
  const [path] = [`/accounts/${id}/avatar`] as const;

  return `${API_ADDRESS}/api${path}`;
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
