import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { VStack, styled } from "@/styled-system/jsx";

import { usePasswordUpdate } from "./usePasswordUpdate";

export function PasswordUpdateForm() {
  const {
    form: { register, handlePasswordChange, errors },
    success,
    handleCloseNotification,
  } = usePasswordUpdate();

  return (
    <>
      <p>You can change your password here.</p>
      <styled.form
        display="flex"
        flexDir="column"
        w="full"
        gap="2"
        onSubmit={handlePasswordChange}
      >
        <Input
          maxW="xs"
          type="password"
          autoComplete="current-password"
          placeholder="current password"
          {...register("old")}
        />
        <styled.p color="fg.error" fontSize="sm">
          {errors.old?.message}
        </styled.p>
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder="new password"
          {...register("new")}
        />
        <styled.p color="fg.error" fontSize="sm">
          {errors.new?.message}
        </styled.p>
        <styled.p color="fg.error" fontSize="sm">
          {errors.root?.message}
        </styled.p>
        <VStack alignItems="start" w="full">
          <Button type="submit" variant="subtle">
            Change password
          </Button>
          <Admonition
            value={success}
            onChange={handleCloseNotification}
            kind="success"
            title="Success"
          >
            Your password has been updated.
          </Admonition>
        </VStack>
      </styled.form>
    </>
  );
}
