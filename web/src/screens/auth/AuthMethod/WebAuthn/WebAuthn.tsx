import { Button, Heading, Input, VStack } from "@chakra-ui/react";
import { useWebAuthn } from "./useWebAuthn";

export function WebAuthn() {
  const { register, onSubmit } = useWebAuthn();

  return (
    <VStack gap={4}>
      <Heading size="sm">Passkey</Heading>
      <form onSubmit={onSubmit}>
        <VStack>
          <Input
            {...register("username")}
            variant="filled"
            placeholder="username"
            w="10em"
          />
          <Button type="submit" width="full">
            Authenticate
          </Button>
        </VStack>
      </form>
    </VStack>
  );
}
