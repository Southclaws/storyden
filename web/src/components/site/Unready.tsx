"use client";

import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";

import { Box, CardBox, Center, HStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import { Spinner } from "../ui/Spinner";

type Props = {
  error?: unknown;
};

export function Unready({ error }: Props) {
  if (!error) {
    return (
      <Center w="full" h="full">
        <Spinner />
      </Center>
    );
  }

  const message = deriveError(error);

  return (
    <HStack maxW="xs" alignItems="center" color="fg.subtle">
      <Box w="5" flexShrink="0">
        <ExclamationTriangleIcon />
      </Box>
      <p id="error__message">{message}</p>
    </HStack>
  );
}

export function UnreadyBanner({ error }: Props) {
  if (!error) {
    return (
      <Center w="full" height="96">
        <Spinner />
      </Center>
    );
  }

  const message = deriveError(error);

  return (
    <Center width="full" justifyContent="center">
      <CardBox maxW="xs">
        <HStack id="error__heading" gap="2" alignItems="center">
          <ExclamationTriangleIcon width={24} height={24} />
          <styled.h1 fontSize="md" fontWeight="bold" my="0">
            Something went wrong
          </styled.h1>
        </HStack>

        <styled.p id="error__message">
          <span>{message}</span>
        </styled.p>
      </CardBox>
    </Center>
  );
}
