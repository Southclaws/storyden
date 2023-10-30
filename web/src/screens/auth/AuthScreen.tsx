"use client";

import { Box, Flex, HStack, Heading, Text, VStack } from "@chakra-ui/react";

import { Back } from "src/components/site/Action/Action";

import { AuthMethod } from "./components/AuthMethod/AuthMethod";

type Props = {
  method: string;
};
export function AuthScreen({ method }: Props) {
  return (
    <Flex
      id="AuthScreen"
      height="100vh"
      width="full"
      justifyContent="center"
      flexDirection="column"
      alignItems="center"
      backgroundPosition="center"
      backgroundSize="cover"
      gap={4}
      padding={6}
    >
      <VStack
        width="full"
        gap={2}
        p={6}
        borderRadius="lg"
        maxW="xs"
        bg="whiteAlpha.700"
        boxShadow="0 10px 30px rgba(0, 0, 0, 0.05)"
      >
        <HStack w="full" justifyContent="space-between">
          <Back href={method ? "/auth" : "/"} />

          <Box>
            <Heading size="md">
              Sign up
              <br />
            </Heading>
            <Text size="sm" fontWeight="medium" color="blackAlpha.600">
              or sign in
            </Text>
          </Box>

          <Box w="1.4em" />
        </HStack>

        <AuthMethod method={method} />
      </VStack>
    </Flex>
  );
}
