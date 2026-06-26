import { PropsWithChildren } from "react";

import { LinkButton } from "@/components/ui/link-button";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";
import { HStack, VStack } from "@/styled-system/jsx";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Layout({ children }: PropsWithChildren) {
  const { registration_mode } = await getSettings();
  const canRegister = allowsPublicRegistration(registration_mode);

  return (
    <VStack w="full">
      {children}

      <HStack>
        <LinkButton size="xs" variant="ghost" href="/password-reset">
          Forgot password
        </LinkButton>

        {canRegister && (
          <LinkButton size="xs" variant="subtle" href="/register">
            Register
          </LinkButton>
        )}
      </HStack>
    </VStack>
  );
}
