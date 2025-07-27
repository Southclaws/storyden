"use client";

import { usePathname } from "next/navigation";

import { UnreadyBanner } from "src/components/site/Unready";

import { Button } from "@/components/ui/button";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, VStack } from "@/styled-system/jsx";

export function GenericError({
  reset,
  message,
}: {
  reset?: () => void;
  message?: string;
}) {
  const pathName = usePathname();

  const isHome = pathName === "/";

  return (
    <VStack p="4" h="dvh" justify="center">
      <VStack maxW="sm" minH="60" gap="8">
        <UnreadyBanner error={message ?? "An unexpected error occurred."} />
        <HStack>
          {!isHome && (
            <LinkButton variant="subtle" href="/">
              Home
            </LinkButton>
          )}
          <Button variant="outline" onClick={reset}>
            Retry
          </Button>
        </HStack>
      </VStack>
    </VStack>
  );
}
