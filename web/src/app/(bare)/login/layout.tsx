import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <VStack w="full">
      {children}

      {/* TODO: OAuth2 providers */}
      {/* <AuthSelection /> */}

      <LinkButton size="xs" variant="subtle" href="/register">
        Register
      </LinkButton>
    </VStack>
  );
}
