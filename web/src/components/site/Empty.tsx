import { PropsWithChildren } from "react";

import { Text } from "@/components/ui/text";
import { VStack } from "@/styled-system/jsx";

import { EmptyIcon } from "../ui/icons/Empty";

/** @deprecated use EmptyState */
export function Empty({ children }: PropsWithChildren) {
  return (
    <VStack alignItems="center" color="text.subtle">
      <EmptyIcon />
      <Text variant="supporting" fontStyle="italic" textWrap="nowrap">
        {children}
      </Text>
    </VStack>
  );
}
