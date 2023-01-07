import { Flex, Heading, Box, Text } from "@chakra-ui/layout";
import React, { FC } from "react";
import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";

type Props = {
  error: string;
  message?: string | undefined;
};

const ErrorBanner: FC<Props> = ({ error, message, ...rest }) => {
  return (
    <Flex width="full" justifyContent="center">
      <Flex flexDir="column" gap={2} bgColor="red.50" borderRadius="xl" p={4}>
        <Flex gap={2} alignItems="center">
          <ExclamationTriangleIcon width={4} height={4} />
          <Heading fontSize="md" my="0">
            that&apos;s a yikes from me
          </Heading>
        </Flex>

        <Box>
          <Text>{message ?? "something bad happened :("}</Text>
        </Box>

        <Text fontSize="sm">
          <details>
            <summary>
              <em>the code mumbo jumbo for nerds</em>
            </summary>
            <pre>{JSON.stringify(rest, null, 2)}</pre>
          </details>
        </Text>
      </Flex>
    </Flex>
  );
};

export default ErrorBanner;
