"use client";

import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form-error-text";
import { BiometricIcon } from "@/components/ui/icons/Biometric";
import { Input } from "@/components/ui/input";
import { Flex, styled } from "@/styled-system/jsx";

import { Props, useRegisterHandleForm } from "./useRegisterHandleForm";

export function RegisterHandleForm(props: Props) {
  const {
    form: {
      register,
      isWebauthnEnabled,
      handlePassword,
      handleWebauthn,
      errors,
    },
  } = useRegisterHandleForm(props);

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
          autoComplete="new-password"
          {...register("token")}
        />
        {props.webauthn && isWebauthnEnabled && (
          <>
            <span>or</span>

            <Button
              w="full"
              variant="ghost"
              size="sm"
              type="button"
              onClick={handleWebauthn}
            >
              <styled.span display="flex" gap="1" alignItems="center" px="4">
                device
                <BiometricIcon />
              </styled.span>
            </Button>
          </>
        )}
      </Flex>
      <FormErrorText>{errors.token?.message}</FormErrorText>
      <Button type="submit" w="full" onClick={handlePassword}>
        Register
      </Button>
      <FormErrorText>{errors.root?.message}</FormErrorText>
    </styled.form>
  );
}
