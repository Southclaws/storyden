import { Box, Button, Heading, Input, VStack } from "@chakra-ui/react";
import { useWebAuthn } from "./useWebAuthn";

export function WebAuthn() {
  const { register, onSubmit } = useWebAuthn();

  return (
    <VStack gap={4}>
      <Heading size="sm">Passkey</Heading>
      <Box as="form" width="full" onSubmit={onSubmit}>
        <VStack>
          <Input
            {...register("username")}
            variant="filled"
            placeholder="username"
            width="full"
          />
          <Button type="submit" width="full">
            Authenticate
          </Button>
        </VStack>
      </Box>
    </VStack>
  );
}
