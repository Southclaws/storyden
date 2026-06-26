import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { VStack } from "@/styled-system/jsx";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <VStack w="full">
      {children}

      <LinkButton size="xs" variant="subtle" href="/login">
        Sign in
      </LinkButton>
    </VStack>
  );
}
