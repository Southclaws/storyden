"use client";

import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form-error-text";
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
        autoCapitalize="none"
        autoCorrect="off"
        autoComplete="username"
        w="full"
        size="sm"
        textAlign="center"
        placeholder="username"
        required
        {...register("identifier")}
      />
      <FormErrorText>{errors.identifier?.message}</FormErrorText>
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
      <FormErrorText>{errors.token?.message}</FormErrorText>
      <Button type="submit" w="full" onClick={handlePassword}>
        Login
      </Button>
      <FormErrorText>{errors.root?.message}</FormErrorText>
    </styled.form>
  );
}
