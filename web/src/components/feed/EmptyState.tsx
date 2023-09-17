import { Text, VStack } from "@chakra-ui/react";

import { ThreadStaircase } from "src/components/graphics/ThreadStaircase";

export function EmptyState() {
  return (
    <VStack height="full" justify="center" pb={32}>
      <ThreadStaircase width="100%" />
      <Text textAlign="center" fontStyle="italic" color="gray.500">
        *tumbleweed*&nbsp;there&nbsp;are&nbsp;no&nbsp;posts...
        you&nbsp;could&nbsp;be&nbsp;the&nbsp;first!
      </Text>
    </VStack>
  );
}
