import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";

import { APIError } from "src/api/openapi-schema";

import { Box, CardBox, Center, HStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import { Spinner } from "../ui/Spinner";

type Props = {
  error?: unknown;
};

export default function ErrorBanner({
  message,
  error,
  metadata,
}: Partial<APIError>) {
  const showDetails = error !== undefined || metadata !== undefined;

  return (
    <Center width="full" justifyContent="center">
      <CardBox maxW="xs">
        <HStack id="error__heading" gap="2" alignItems="center">
          <ExclamationTriangleIcon width={32} height={32} />
          <styled.h1 fontSize="md" fontWeight="bold" my="0">
            Something went wrong
          </styled.h1>
        </HStack>

        <styled.p id="error__message">
          <span>{message ?? "something bad happened :("}</span>
        </styled.p>

        {showDetails && (
          <details id="error__details-container">
            <summary id="error__details-summary">
              <em>more information</em>
            </summary>
            <em>{error}</em>
            <pre id="error__details-code">
              {JSON.stringify(metadata, null, 2)}
            </pre>
          </details>
        )}
      </CardBox>
    </Center>
  );
}

export function Unready({ error }: Props) {
  if (error) {
    const message = deriveError(error);

    return (
      <CardBox maxW="xs" alignItems="center">
        <Box w="8" color="fg.subtle">
          <ExclamationTriangleIcon />
        </Box>
        {message}
      </CardBox>
    );
  }

  return (
    <Center h="full">
      <Spinner />
    </Center>
  );
}

export function UnreadyBanner(props: Props) {
  return (
    <Center w="full" height="96">
      <Unready {...props} />
    </Center>
  );
}

export const handleError = (e: APIError) => {
  console.error(e);
};
