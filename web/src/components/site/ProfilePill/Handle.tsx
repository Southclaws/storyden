import { CheckBadgeIcon } from "@heroicons/react/24/outline";

import { ProfileReference } from "src/api/openapi/schemas";

import { Center, styled } from "@/styled-system/jsx";

export type Props = {
  profileReference: ProfileReference;
  size?: "sm" | "lg";
};

export function Handle({ profileReference, size }: Props) {
  return (
    <styled.p fontSize={size === "lg" ? "md" : "sm"} display="flex">
      <styled.span
        whiteSpace="nowrap"
        textOverflow="ellipsis"
        overflow="hidden"
      >
        @{profileReference.handle}
      </styled.span>

      {profileReference.admin && (
        <Center title="Admin">
          <CheckBadgeIcon height="1rem" />
        </Center>
      )}
    </styled.p>
  );
}
