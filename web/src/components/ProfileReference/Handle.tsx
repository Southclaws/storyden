"use client";

import { Icon, Text } from "@chakra-ui/react";
import { CheckBadgeIcon } from "@heroicons/react/24/outline";

import { ProfileReference } from "src/api/openapi/schemas";

export type Props = {
  profileReference: ProfileReference;
  size?: "sm" | "lg";
};

export function Handle({ profileReference, size }: Props) {
  return (
    <Text
      fontSize={size === "lg" ? "md" : "sm"}
      display="flex"
      alignItems="center"
      gap={1}
    >
      @{profileReference.handle}
      {profileReference.admin && (
        <Text as="span" title="Admin">
          <Icon>
            <CheckBadgeIcon />
          </Icon>
        </Text>
      )}
    </Text>
  );
}
