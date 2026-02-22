import { AccountRoleRefList, ProfileReference } from "@/api/openapi-schema";
import { RoleBadgeList } from "@/components/role/RoleBadge/RoleBadgeList";
import { isDefaultRole } from "@/lib/role/defaults";
import { parseRoleMetadata } from "@/lib/role/metadata";
import { Flex, HStack, styled } from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";

import { MemberAvatar } from "./MemberAvatar";

export type Props = {
  profile: ProfileReference;
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "handle" | "full-horizontal" | "full-vertical";
  showRoles?: "hidden" | "badge" | "all";
  avatar?: "hidden" | "visible";
};

export function MemberIdent({
  profile,
  size,
  name = "hidden",
  avatar = "visible",
  showRoles = "hidden",
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
      />
    </HStack>
  );
}

function topRoleForDecoration(roles: AccountRoleRefList) {
  return roles.find((r) => !isDefaultRole(r));
}

function getRoleDecorationStyle(
  roles: AccountRoleRefList,
  defaultColour: string,
  defaultWeight: string,
  boldWeight: string,
) {
  const role = topRoleForDecoration(roles);
  const metadata = role ? parseRoleMetadata(role.meta) : null;
  const resolvedWeight = metadata?.bold ? boldWeight : defaultWeight;
  const resolvedStyle = metadata?.italic ? "italic" : "normal";
  const resolvedColour =
    metadata?.coloured && role?.colour ? role.colour : defaultColour;

  return {
    "--colors-color-palette": resolvedColour,
    "--decoration-font-style": resolvedStyle,
    "--decoration-font-weight": resolvedWeight,
  } as any;
}

export function MemberName({
  profile,
  size,
  name = "hidden",
  showRoles = "hidden",
}: Props) {
  switch (name) {
    case "full-horizontal": {
      const decorationStyle = getRoleDecorationStyle(
        profile.roles,
        token("colors.fg.default"),
        size === "lg"
          ? token("fontWeights.semibold")
          : token("fontWeights.normal"),
        size === "lg"
          ? token("fontWeights.extrabold")
          : token("fontWeights.bold"),
      );

      return (
        <Flex
          className="member-name__show-horizontal"
          maxW="full"
          direction="row"
          gap="1"
          alignItems="center"
          style={decorationStyle}
        >
          <styled.p
            className="member-name__show-horizontal-display-name"
            minW="0"
            fontSize={size}
            fontWeight="var(--decoration-font-weight)"
            fontStyle="var(--decoration-font-style)"
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="colorPalette"
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
            _containerSmall={{
              color: "colorPalette",
              fontWeight: "var(--decoration-font-weight)",
              fontStyle: "var(--decoration-font-style)",
            }}
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} showRoles={showRoles} />
        </Flex>
      );
    }

    case "full-vertical": {
      const decorationStyle = getRoleDecorationStyle(
        profile.roles,
        token("colors.fg.default"),
        size === "lg"
          ? token("fontWeights.bold")
          : token("fontWeights.semibold"),
        token("fontWeights.extrabold"),
      );

      return (
        <Flex
          className="member-name__show-vertical"
          maxW="full"
          minW="0"
          direction="column"
          gap="0"
          alignItems="start"
          style={decorationStyle}
        >
          <styled.p
            className="member-name__show-vertical-display-name"
            minW="0"
            fontSize={size}
            fontWeight="var(--decoration-font-weight)"
            fontStyle="var(--decoration-font-style)"
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="colorPalette"
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
          <Roles profile={profile} showRoles={showRoles} />
        </Flex>
      );
    }

    case "handle": {
      const decorationStyle = getRoleDecorationStyle(
        profile.roles,
        token("colors.fg.subtle"),
        size === "lg"
          ? token("fontWeights.medium")
          : token("fontWeights.normal"),
        size === "lg"
          ? token("fontWeights.bold")
          : token("fontWeights.semibold"),
      );

      return (
        <HStack
          className="member-name__show-handle"
          maxW="full"
          gap="1"
          minW="0"
          style={decorationStyle}
        >
          <styled.p
            className="member-name__show-handle-handle"
            fontSize={size}
            fontWeight="var(--decoration-font-weight)"
            fontStyle="var(--decoration-font-style)"
            overflowX="hidden"
            overflowY="clip"
            textWrap="nowrap"
            textOverflow="ellipsis"
            lineHeight="tight"
            color="colorPalette"
          >
            @{profile.handle}
          </styled.p>
          <Roles profile={profile} showRoles={showRoles} />
        </HStack>
      );
    }

    case "hidden":
      return null;
  }
}

function Roles({ profile, showRoles }: Pick<Props, "profile" | "showRoles">) {
  if (!showRoles) {
    return null;
  }

  if (showRoles === "hidden") {
    return null;
  }

  return (
    <RoleBadgeList
      roles={profile.roles}
      onlyBadgeRole={showRoles === "badge"}
    />
  );
}
