import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form-error-text";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

import { usePasswordCreate } from "./usePasswordCreate";

export function PasswordCreateForm() {
  const {
    form: { register, handlePasswordCreate, errors },
    success,
    handleCloseNotification,
  } = usePasswordCreate();

  return (
    <>
      <p>
        Your account does not currently have a password. You can add a password
        here.
      </p>
      <styled.form
        display="flex"
        flexDir="column"
        w="full"
        gap="2"
        onSubmit={handlePasswordCreate}
      >
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder="new password"
          disabled={success}
          {...register("newPassword")}
        />
        <FormErrorText>{errors.newPassword?.message}</FormErrorText>
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder="confirm new password"
          disabled={success}
          {...register("confirmPassword")}
        />
        <FormErrorText>{errors.confirmPassword?.message}</FormErrorText>
        <FormErrorText>{errors.root?.message}</FormErrorText>
        <VStack alignItems="start" w="full">
          <Button type="submit" disabled={success}>
            Add password
          </Button>
          <Admonition
            value={success}
            onChange={handleCloseNotification}
            kind="success"
            title="Success"
          >
            Your account now has a password! You can now use this to log in, but
            you can continue to use your other authentication methods as well.
          </Admonition>
        </VStack>
      </styled.form>
    </>
  );
}
