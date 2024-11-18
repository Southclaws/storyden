import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { HStack, VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <VStack w="full">
      {children}

      <HStack>
        <LinkButton size="xs" variant="ghost" href="/login">
          Login
        </LinkButton>

        <LinkButton size="xs" variant="subtle" href="/register">
          Register
        </LinkButton>
      </HStack>
    </VStack>
  );
}
