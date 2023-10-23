import { CreateToastFnReturn } from "@chakra-ui/react";
import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";

import { APIError } from "src/api/openapi/schemas";

import { Flex, styled } from "@/styled-system/jsx";

export default function ErrorBanner({
  message,
  error,
  metadata,
}: Partial<APIError>) {
  const showDetails = error !== undefined || metadata !== undefined;

  return (
    <Flex width="full" justifyContent="center">
      <Flex flexDir="column" gap={2} bgColor="red.50" borderRadius="xl" p={4}>
        <Flex id="error__heading" gap={2} alignItems="center">
          <ExclamationTriangleIcon width={32} height={32} />
          <styled.h1 fontSize="md" fontWeight="bold" my="0">
            that&apos;s a yikes from me
          </styled.h1>
        </Flex>

        <styled.p id="error__message">
          <span>{message ?? "something bad happened :("}</span>
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
        </styled.p>
      </Flex>
    </Flex>
  );
}

export const errorToast = (toast: CreateToastFnReturn) => (e: APIError) => {
  console.error(e);
  toast({
    title: "Error",
    status: "error",
    description: e.message,
  });
};
