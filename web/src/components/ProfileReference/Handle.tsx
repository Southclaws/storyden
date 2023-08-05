import { Text } from "@chakra-ui/react";

import { ProfileReference } from "src/api/openapi/schemas";

export type Props = {
  profileReference: ProfileReference;
  size?: "sm" | "lg";
};

export function Handle({ profileReference, size }: Props) {
  return (
    <Text fontSize={size === "lg" ? "md" : "sm"}>
      @{profileReference.handle}
    </Text>
  );
}
