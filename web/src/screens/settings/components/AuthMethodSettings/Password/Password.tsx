import { Heading } from "@chakra-ui/react";

import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

import { VStack, styled } from "@/styled-system/jsx";

import { usePassword } from "./usePassword";

export function Password() {
  const {
    form: { register, handlePasswordChange, errors },
  } = usePassword();
  return (
    <VStack alignItems="start">
      <Heading size="sm">Password</Heading>

      <p>You can change your password here.</p>

      <styled.form
        display="flex"
        flexDir="column"
        gap="2"
        onSubmit={handlePasswordChange}
      >
        <Input
          type="password"
          autoComplete="current-password"
          placeholder="current password"
          {...register("old")}
        />
        <styled.p color="red.600" fontSize="sm">
          {errors.old?.message}
        </styled.p>
        <Input
          type="password"
          autoComplete="new-password"
          placeholder="new password"
          {...register("new")}
        />
        <styled.p color="red.600" fontSize="sm">
          {errors.new?.message}
        </styled.p>
        <styled.p color="red.600" fontSize="sm">
          {errors.root?.message}
        </styled.p>
        <Button type="submit">Change password</Button>
      </styled.form>
    </VStack>
  );
}
