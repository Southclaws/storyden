import { ProfileReference } from "@/api/openapi-schema";
import { Flex, HStack, styled } from "@/styled-system/jsx";

import { RoleBadgeList } from "../RoleBadge/RoleBadgeList";

import { MemberAvatar } from "./MemberAvatar";

export type Props = {
  profile: ProfileReference;
  size?: "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full";
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
    <HStack
      minW="0"
      w="full"
      overflowY="hidden"
      gap={size === "lg" ? "2" : "1"}
    >
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
    case "full":
      return (
        <Flex
          direction={size === "lg" ? "column" : "row"}
          gap={size === "lg" ? "0" : "1"}
          alignItems={size === "lg" ? "start" : "center"}
        >
          <styled.p
            minW="0"
            fontSize={size === "lg" ? "lg" : "sm"}
            fontWeight={size === "lg" ? "bold" : "medium"}
            overflowX="hidden"
            textWrap="nowrap"
            textOverflow="ellipsis"
            color="fg.default"
            _containerSmall={{
              display: "none",
            }}
          >
            {profile.name}
          </styled.p>
          <styled.p textWrap="nowrap" color="fg.muted">
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} roles={roles} />
        </Flex>
      );

    case "handle":
      return (
        <HStack gap="1">
          <styled.p color="fg.muted">@{profile.handle}</styled.p>
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
