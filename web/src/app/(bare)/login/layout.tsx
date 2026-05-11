import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { tServer } from "@/i18n/server";
import { HStack, VStack } from "@/styled-system/jsx";

export default async function Layout({ children }: PropsWithChildren) {
  const forgotPassword = await tServer("Forgot password");
  const register = await tServer("Register");

  return (
    <VStack w="full">
      {children}

      <HStack>
        <LinkButton size="xs" variant="ghost" href="/password-reset">
          {forgotPassword}
        </LinkButton>

        <LinkButton size="xs" variant="subtle" href="/register">
          {register}
        </LinkButton>
      </HStack>
    </VStack>
  );
}
