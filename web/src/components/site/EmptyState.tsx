import { ThreadStaircase } from "src/components/graphics/ThreadStaircase";

import { VStack } from "@/styled-system/jsx";
import { styled } from "@/styled-system/jsx";

export function EmptyState() {
  return (
    <VStack height="full" justify="center" pb="32">
      <ThreadStaircase width="100%" />
      <styled.p textAlign="center" fontStyle="italic" color="gray.5">
        *tumbleweed*&nbsp;there&apos;s&nbsp;nothing&nbsp;here...
      </styled.p>
    </VStack>
  );
}
