import { UIMessagePart } from "ai";

import { Box, HStack, VStack, styled } from "@/styled-system/jsx";

type Props = {
  part: UIMessagePart<any, any>;
};

export function RobotToolCall({ part }: Props) {
  return (
    <VStack
      gap="2"
      p="2"
      bg="bg.muted"
      borderRadius="md"
      fontSize="sm"
      w="full"
      maxW="5/6"
      alignSelf="flex-start"
    >
      <styled.details w="full">
        <styled.summary cursor="pointer" color="fg.subtle">
          {part.type}
        </styled.summary>
        <styled.pre
          fontSize="xs"
          p="2"
          bg="bg.subtle"
          borderRadius="sm"
          overflow="auto"
          maxH="32"
          mt="1"
        >
          {JSON.stringify(part, null, 2)}
        </styled.pre>
      </styled.details>

      <HStack w="full">
        <Box fontSize="xs" color="fg.muted">
          state: {part["state"]}
        </Box>
      </HStack>
    </VStack>
  );
}
