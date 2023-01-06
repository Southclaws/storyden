import {
  Box,
  Button,
  FormLabel,
  Heading,
  Input,
  VStack,
} from "@chakra-ui/react";
import useSignUp from "./useSignUp";

export function SignUp() {
  const {
    form: { register, onSubmit, errors },
  } = useSignUp();

  return (
    <VStack gap={4}>
      <Heading size="sm">Sign up</Heading>
      <Box as="form" width="full" onSubmit={onSubmit}>
        <VStack>
          <Input
            {...register("identifier")}
            variant="filled"
            placeholder="username"
          />
          <FormLabel>{errors.identifier?.message}</FormLabel>

          <Input
            {...register("token")}
            variant="filled"
            placeholder="password"
          />
          <FormLabel>{errors.token?.message}</FormLabel>

          <Button type="submit" width="full">
            Login
          </Button>
        </VStack>
      </Box>
    </VStack>
  );
}
