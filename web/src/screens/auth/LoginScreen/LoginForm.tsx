"use client";

import { FingerPrintIcon } from "@heroicons/react/24/outline";

import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

import { Flex, styled } from "@/styled-system/jsx";

import { Props, useLoginForm } from "./useLoginForm";

export function LoginForm(props: Props) {
  const {
    form: {
      register,
      isWebauthnEnabled,
      handlePassword,
      handleWebauthn,
      errors,
    },
  } = useLoginForm();

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
      <styled.p color="red.600" fontSize="sm">
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
        {props.webauthn && isWebauthnEnabled && (
          <>
            <styled.span>or</styled.span>

            <Button
              w="full"
              kind="neutral"
              size="sm"
              type="button"
              onClick={handleWebauthn}
            >
              <styled.span display="flex" gap="1" alignItems="center" px="4">
                device
                <FingerPrintIcon />
              </styled.span>
            </Button>
          </>
        )}
      </Flex>
      <styled.p color="red.600" fontSize="sm">
        {errors.token?.message}
      </styled.p>
      <Button type="submit" w="full" kind="primary" onClick={handlePassword}>
        Login
      </Button>
      <styled.p color="red.600" fontSize="sm">
        {errors.root?.message}
      </styled.p>
    </styled.form>
  );
}
