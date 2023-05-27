import {
  Box,
  Button,
  FormLabel,
  Heading,
  HStack,
  Input,
  VStack,
} from "@chakra-ui/react";
import { useWebAuthn } from "./useWebAuthn";

export function WebAuthn() {
  const { register, onSubmit, errors } = useWebAuthn();

  return (
    <VStack gap={4}>
      <Heading size="sm">Passkey</Heading>
      <Box as="form" width="full">
        <VStack>
          <Input
            {...register("username")}
            variant="filled"
            placeholder="username"
            width="full"
            bgColor="whiteAlpha.900"
            required
          />
          <FormLabel>{errors.username?.message}</FormLabel>

          <HStack width="full">
            <Button
              type="submit"
              name="auth"
              value="signup"
              width="full"
              bgColor="primary.500"
              color="light"
              variant="primary"
              onClick={onSubmit("signup")}
            >
              Register
            </Button>
            <Button
              type="submit"
              name="auth"
              value="signin"
              width="full"
              variant="secondary"
              onClick={onSubmit("signin")}
            >
              Login
            </Button>
          </HStack>
        </VStack>
      </Box>
    </VStack>
  );
}
