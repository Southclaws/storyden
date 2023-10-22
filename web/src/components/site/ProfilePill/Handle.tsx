import { CheckBadgeIcon } from "@heroicons/react/24/outline";

import { ProfileReference } from "src/api/openapi/schemas";

import { styled } from "@/styled-system/jsx";

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
        <styled.span title="Admin">
          <CheckBadgeIcon height="1rem" />
        </styled.span>
      )}
    </styled.p>
  );
}
