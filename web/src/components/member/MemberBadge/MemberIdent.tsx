import { ProfileReference } from "@/api/openapi-schema";
import { Flex, HStack, styled } from "@/styled-system/jsx";

import { RoleBadgeList } from "../RoleBadge/RoleBadgeList";

import { MemberAvatar } from "./MemberAvatar";

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  roles?: "hidden" | "badge" | "all";
  avatar?: "hidden" | "visible";
};

export function MemberIdent({
  profile,
  size,
  name = "hidden",
  avatar = "visible",
  roles = "hidden",
}: Props) {
  return (
    <HStack minW="0" w="full" overflow="hidden" gap={size === "lg" ? "2" : "1"}>
      {avatar === "visible" && <MemberAvatar profile={profile} size={size} />}
      <MemberName profile={profile} size={size} name={name} roles={roles} />
    </HStack>
  );
}

export function MemberName({
  profile,
  size,
  name = "hidden",
  roles = "hidden",
}: Props) {
  switch (name) {
    case "full-horizontal":
      return (
        <Flex direction="row" gap="1" alignItems="center">
          <styled.p
            minW="0"
            fontSize={size}
            fontWeight={size === "lg" ? "bold" : "medium"}
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="fg.default"
            _containerSmall={{
              display: "none",
            }}
          >
            {profile.name}
          </styled.p>
          <styled.p
            fontSize={size}
            fontWeight="normal"
            textWrap="nowrap"
            color="fg.subtle"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} roles={roles} />
        </Flex>
      );

    case "full-vertical":
      return (
        <Flex direction="column" gap="0" alignItems="start">
          <styled.p
            minW="0"
            fontSize={size}
            fontWeight={size === "lg" ? "bold" : "medium"}
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="fg.default"
            _containerSmall={{
              display: "none",
            }}
          >
            {profile.name}
          </styled.p>
          <styled.p
            fontSize={size}
            fontWeight="normal"
            textWrap="nowrap"
            color="fg.subtle"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} roles={roles} />
        </Flex>
      );

    case "handle":
      return (
        <HStack gap="1">
          <styled.p
            fontSize={size}
            fontWeight="normal"
            textWrap="nowrap"
            color="fg.subtle"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} roles={roles} />
        </HStack>
      );

    case "hidden":
      return null;
  }
}

function Roles({ profile, roles }: Pick<Props, "profile" | "roles">) {
  if (roles === "hidden") {
    return null;
  }

  return (
    <RoleBadgeList roles={profile.roles} onlyBadgeRole={roles === "badge"} />
  );
}
