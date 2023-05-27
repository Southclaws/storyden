import {
  Box,
  Button,
  FormLabel,
  HStack,
  Input,
  VStack,
} from "@chakra-ui/react";
import { usePassword } from "./usePassword";

export function Password() {
  const {
    form: { register, onSubmit, errors },
  } = usePassword();

  return (
    <VStack gap={4}>
      <Box as="form" width="full">
        <VStack>
          <Input
            {...register("identifier")}
            placeholder="username"
            bgColor="whiteAlpha.900"
            required
          />
          <FormLabel>{errors.identifier?.message}</FormLabel>

          <Input
            {...register("token")}
            placeholder="password"
            bgColor="whiteAlpha.900"
            type="password"
            required
          />
          <FormLabel>{errors.token?.message}</FormLabel>

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
