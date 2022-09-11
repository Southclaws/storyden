import { Box, Button, Heading, Input, VStack } from "@chakra-ui/react";

export function Password() {
  return (
    <VStack gap={4}>
      <Heading size="sm">Sign in</Heading>
      <form>
        <VStack>
          <Input variant="filled" placeholder="username" />
          <Input variant="filled" placeholder="password" />
          <Button type="submit" width="full">
            Login
          </Button>
        </VStack>
      </form>
    </VStack>
  );
}
