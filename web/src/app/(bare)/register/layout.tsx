import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <VStack className="sd-screen sd-screen--register" w="full">
      {children}

      <LinkButton size="xs" variant="subtle" href="/login">
        Sign in
      </LinkButton>
    </VStack>
  );
}
