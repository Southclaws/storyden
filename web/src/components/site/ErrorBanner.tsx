import { Box, Flex, Heading, Text } from "@chakra-ui/layout";
import { CreateToastFnReturn } from "@chakra-ui/react";
import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";

import { APIError, APIErrorMetadata } from "src/api/openapi/schemas";

type Props = {
  message?: string | undefined;
  metadata?: APIErrorMetadata | undefined;
};

export default function ErrorBanner({ message, metadata }: Props) {
  return (
    <Flex width="full" justifyContent="center">
      <Flex flexDir="column" gap={2} bgColor="red.50" borderRadius="xl" p={4}>
        <Flex gap={2} alignItems="center">
          <ExclamationTriangleIcon width={32} height={32} />
          <Heading fontSize="md" my="0">
            that&apos;s a yikes from me
          </Heading>
        </Flex>

        <Box>
          <Text>{message ?? "something bad happened :("}</Text>
        </Box>

        {metadata && (
          <Text fontSize="sm">
            <details>
              <summary>
                <em>the code mumbo jumbo for nerds</em>
              </summary>
              <pre>{JSON.stringify(metadata, null, 2)}</pre>
            </details>
          </Text>
        )}
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
