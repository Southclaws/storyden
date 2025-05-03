"use client"; // Error boundaries must be Client Components

import { Box, Card, LinkButton, VStack } from "@/styled-system/jsx";
import { useEffect } from "react";

type Props = {
  error: Error & { digest?: string };
  reset: () => void;
};

export default function Error({ error, reset }: Props) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <VStack w="full" p="16">
      <h2>Something went wrong</h2>
      <Card backgroundColor="Mono.slush" maxW="prose">
        <p>{error.message}</p>
      </Card>
      <LinkButton
        p="2"
        fontSize="sm"
        lineHeight="tight"
        height="auto"
        onClick={() => reset()}
      >
        Try again
      </LinkButton>
    </VStack>
  );
}
