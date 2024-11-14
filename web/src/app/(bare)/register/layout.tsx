import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { AuthSelection } from "@/screens/auth/components/AuthSelection/AuthSelection";
import { VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <VStack w="full">
      {children}

      {/* TODO: OAuth2 providers */}
      {/* <AuthSelection /> */}

      <LinkButton size="xs" variant="subtle" href="/login">
        Sign in
      </LinkButton>
    </VStack>
  );
}
