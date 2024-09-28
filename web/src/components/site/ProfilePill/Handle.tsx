import { ShieldCheckIcon } from "@heroicons/react/24/outline";

import { Permission, ProfileReference } from "src/api/openapi-schema";

import { styled } from "@/styled-system/jsx";

export type Props = {
  profileReference: ProfileReference;
  size?: "sm" | "lg";
};

export function Handle({ profileReference, size }: Props) {
  const isAdmin = profileReference.roles.find((role) =>
    role.permissions.includes(Permission.ADMINISTRATOR),
  );

  return (
    <styled.p fontSize={size === "lg" ? "md" : "sm"} display="flex" gap="1">
      <styled.span
        whiteSpace="nowrap"
        textOverflow="ellipsis"
        overflow="hidden"
      >
        @{profileReference.handle}
      </styled.span>

      {isAdmin && (
        <styled.span
          display="flex"
          justifyContent="center"
          alignItems="center"
          title="Admin"
        >
          <ShieldCheckIcon height="1rem" />
        </styled.span>
      )}
    </styled.p>
  );
}
