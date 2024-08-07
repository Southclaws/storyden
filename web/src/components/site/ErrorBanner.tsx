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
      <Flex flexDir="column" gap="2" bgColor="bg.error" borderRadius="xl" p="4">
        <Flex id="error__heading" gap="2" alignItems="center">
          <ExclamationTriangleIcon width={32} height={32} />
          <styled.h1 fontSize="md" fontWeight="bold" my="0">
            that&apos;s a yikes from me
          </styled.h1>
        </Flex>

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
      </Flex>
    </Flex>
  );
}

export const handleError = (e: APIError) => {
  console.error(e);
};
