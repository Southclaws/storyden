import { UnreadyBanner } from "src/components/site/Unready";

import { LinkButton } from "@/components/ui/link-button";
import { VStack } from "@/styled-system/jsx";

export function ThreadNotFoundError() {
  return (
    <VStack p="4" h="dvh" justify="center">
      <VStack maxW="sm" minH="60" gap="8">
        <UnreadyBanner error="The link to this thread did not lead anywhere." />
        <LinkButton variant="subtle" href="/">
          Home
        </LinkButton>
      </VStack>
    </VStack>
  );
}
