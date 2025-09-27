import { PropsWithChildren } from "react";

import {
  Box,
  CardBox,
  Center,
  HStack,
  LStack,
  styled,
} from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import { Spinner } from "../ui/Spinner";
import { WarningIcon } from "../ui/icons/Warning";
import { LinkButton } from "../ui/link-button";

type Props = {
  error?: unknown;
};

export function Unready({ error }: Props) {
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
          <Spinner />
        </div>
      </Center>
    );
  }

  const message = deriveError(error);

  return (
    <HStack
      maxW="xs"
      alignItems="center"
      color="fg.subtle"
      role="alert"
      aria-atomic="true"
    >
      <Box w="5" flexShrink="0">
        <WarningIcon aria-hidden="true" />
      </Box>
      <p id="error__message">{message}</p>
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

          <styled.p id="error__message">
            <span>{message}</span>
          </styled.p>

          <LStack>{children}</LStack>
        </LStack>
      </CardBox>
    </Center>
  );
}

export function UnauthenticatedBanner() {
  return (
    <UnreadyBanner error="Please log in to see this page.">
      <HStack w="full">
        <LinkButton w="full" size="xs" href="/register">
          Register
        </LinkButton>
        <LinkButton w="full" size="xs" variant="outline" href="/login">
          Login
        </LinkButton>
      </HStack>
    </UnreadyBanner>
  );
}
