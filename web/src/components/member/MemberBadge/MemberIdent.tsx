import { AccountRoleList, ProfileReference } from "@/api/openapi-schema";
import { RoleBadgeList } from "@/components/role/RoleBadge/RoleBadgeList";
import { Flex, HStack, styled } from "@/styled-system/jsx";

import { MemberAvatar } from "./MemberAvatar";

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  showRoles?: "hidden" | "badge" | "all";
  roles?: AccountRoleList;
  avatar?: "hidden" | "visible";
};

export function MemberIdent({
  profile,
  size,
  name = "hidden",
  avatar = "visible",
  showRoles = "hidden",
  roles,
}: Props) {
  return (
    <HStack
      className="member-ident__container"
      minW="0"
      w="full"
      alignItems="center"
      overflow="hidden"
      gap={size === "lg" ? "2" : "1"}
    >
      {avatar === "visible" && <MemberAvatar profile={profile} size={size} />}
      <MemberName
        profile={profile}
        size={size}
        name={name}
        showRoles={showRoles}
        roles={roles}
      />
    </HStack>
  );
}

export function MemberName({
  profile,
  size,
  name = "hidden",
  showRoles = "hidden",
  roles,
}: Props) {
  switch (name) {
    case "full-horizontal":
      return (
        <Flex
          className="member-name__show-horizontal"
          maxW="full"
          direction="row"
          gap="1"
          alignItems="center"
        >
          <styled.p
            className="member-name__show-horizontal-display-name"
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
            className="member-name__show-horizontal-handle"
            color="fg.subtle"
            minW="0"
            textWrap="nowrap"
            textOverflow="ellipsis"
            overflowX="hidden"
            overflowY="clip"
            lineHeight="tight"
            fontSize={size}
            fontWeight="normal"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} showRoles={showRoles} roles={roles} />
        </Flex>
      );

    case "full-vertical":
      return (
        <Flex
          className="member-name__show-vertical"
          maxW="full"
          minW="0"
          direction="column"
          gap="0"
          alignItems="start"
        >
          <styled.p
            className="member-name__show-vertical-display-name"
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
              display: size === "xs" ? "none" : undefined,
            }}
          >
            {profile.name}
          </styled.p>
          <styled.p
            className="member-name__show-vertical-handle"
            fontSize={size}
            fontWeight="normal"
            w="full"
            minW="0"
            textWrap="nowrap"
            color="fg.subtle"
            overflowX="hidden"
            overflowY="clip"
            textOverflow="ellipsis"
            lineHeight="tight"
            // NOTE: Handles are always lowercase so our x-height upper bound is
            // quite low so we can get away with a tighter line height.
            mt={size === "lg" ? "-1" : "0"}
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} showRoles={showRoles} roles={roles} />
        </Flex>
      );

    case "handle":
      return (
        <HStack
          className="member-name__show-handle"
          maxW="full"
          gap="1"
          minW="0"
        >
          <styled.p
            className="member-name__show-handle-handle"
            fontSize={size}
            fontWeight="normal"
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="fg.subtle"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} showRoles={showRoles} roles={roles} />
        </HStack>
      );

    case "hidden":
      return null;
  }
}

function Roles({
  showRoles,
  roles,
}: Pick<Props, "profile" | "showRoles" | "roles">) {
  if (!showRoles) {
    return null;
  }

  if (showRoles === "hidden" || !roles) {
    return null;
  }

  return <RoleBadgeList roles={roles} onlyBadgeRole={showRoles === "badge"} />;
}
