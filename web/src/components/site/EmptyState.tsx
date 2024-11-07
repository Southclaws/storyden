import { PropsWithChildren } from "react";

import { Center, VStack } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { EmptyIcon } from "../ui/icons/Empty";

export function EmptyState({ children }: PropsWithChildren) {
  return (
    <Center className={vstack()} p="8" gap="2" color="fg.subtle">
      <EmptyIcon />

      <VStack gap="1">
        {children || <p>There&apos;s no content here.</p>}
      </VStack>
    </Center>
  );
}
