import { PropsWithChildren } from "react";

import { CardBox } from "@/components/ui/card-box";
import { Text } from "@/components/ui/text";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import type { Settings } from "@/lib/settings/settings";
import { Box, Center, HStack, LStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import { WarningIcon } from "../ui/icons/Warning";
import { LinkButton } from "../ui/link-button";
import { Spinner } from "../ui/spinner";

type Props = {
  error?: unknown;
  size?: "sm" | "md";
};

export function Unready({ error, size = "md" }: Props) {
  const isCompact = size === "sm";

  if (!error) {
    return (
      <Center
        w="full"
        h="full"
        role="status"
        aria-busy="true"
        aria-live="polite"
        aria-label="Loading"
      >
        <div aria-hidden="true">
          <Spinner size={size} />
        </div>
      </Center>
    );
  }

  const message = deriveError(error);

  return (
    <HStack
      maxW={isCompact ? "full" : "xs"}
      alignItems={isCompact ? "flex-start" : "center"}
      gap="2"
      color="text.muted"
      role="alert"
      aria-atomic="true"
    >
      <Box w={isCompact ? "4" : "5"} flexShrink="0">
        <WarningIcon aria-hidden="true" boxSize={isCompact ? "4" : undefined} />
      </Box>
      <Text
        id="error__message"
        variant={isCompact ? "metadata" : "body"}
        color="text.muted"
        minW="0"
        overflowWrap="anywhere"
      >
        {message}
      </Text>
    </HStack>
  );
}

export function UnreadyBanner({ error, children }: PropsWithChildren<Props>) {
  if (!error) {
    return (
      <Center
        w="full"
        height="96"
        role="status"
        aria-busy="true"
        aria-live="polite"
        aria-label="Loading"
      >
        <Spinner aria-hidden="true" />
      </Center>
    );
  }

  const message = deriveError(error);

  return (
    <Center
      width="full"
      justifyContent="center"
      role="alert"
      aria-atomic="true"
    >
      <CardBox maxW="xs">
        <LStack>
          <HStack id="error__heading" gap="2" alignItems="center">
            <WarningIcon aria-hidden />
            <styled.h1 fontSize="md" fontWeight="bold" my="0">
              Something went wrong
            </styled.h1>
          </HStack>

          <Text id="error__message">
            <span>{message}</span>
          </Text>

          <LStack>{children}</LStack>
        </LStack>
      </CardBox>
    </Center>
  );
}

export function UnauthenticatedBanner({
  initialSettings,
}: {
  initialSettings?: Settings;
}) {
  const canRegister = initialSettings
    ? allowsPublicRegistration(initialSettings.registration_mode)
    : false;

  return (
    <UnreadyBanner error="Please log in to see this page.">
      <HStack w="full">
        {canRegister && (
          <LinkButton w="full" href="/register">
            Register
          </LinkButton>
        )}
        <LinkButton w="full" variant="outline" href="/login">
          Login
        </LinkButton>
      </HStack>
    </UnreadyBanner>
  );
}
