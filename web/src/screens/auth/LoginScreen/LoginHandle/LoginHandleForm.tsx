"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Flex, styled } from "@/styled-system/jsx";

import { useLoginHandleForm } from "./useLoginHandleForm";

export function LoginHandleForm() {
  const {
    form: {
      register,
      isWebauthnEnabled,
      handlePassword,
      handleWebauthn,
      errors,
    },
  } = useLoginHandleForm();

  return (
    <styled.form
      w="full"
      display="flex"
      flexDir="column"
      gap="2"
      textAlign="center"
    >
      <Input
        type="text"
        w="full"
        size="sm"
        textAlign="center"
        placeholder="username"
        required
        {...register("identifier")}
      />
      <styled.p color="fg.error" fontSize="sm">
        {errors.identifier?.message}
      </styled.p>
      <Flex alignItems="center" gap="2">
        <Input
          type="password"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="password"
          autoComplete="current-password"
          {...register("token")}
        />
      </Flex>
      <styled.p color="fg.error" fontSize="sm">
        {errors.token?.message}
      </styled.p>
      <Button type="submit" w="full" onClick={handlePassword}>
        Login
      </Button>
      <styled.p color="fg.error" fontSize="sm">
        {errors.root?.message}
      </styled.p>
    </styled.form>
  );
}
