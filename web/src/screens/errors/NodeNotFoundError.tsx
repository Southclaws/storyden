import { UnreadyBanner } from "src/components/site/Unready";

import { LinkButton } from "@/components/ui/link-button";
import { VStack } from "@/styled-system/jsx";

export function NodeNotFoundError() {
  return (
    <VStack p="4" h="dvh" justify="center">
      <VStack maxW="sm" minH="60" gap="8">
        <UnreadyBanner error="The link to this page did not lead anywhere." />
        <LinkButton variant="subtle" href="/l">
          Library
        </LinkButton>
      </VStack>
    </VStack>
  );
}
