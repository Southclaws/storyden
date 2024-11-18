"use client";

import { Button } from "@/components/ui/button";
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
  } = useRegisterHandleForm();

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
        placeholder="choose your username"
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
          autoComplete="new-password"
          {...register("token")}
        />
        {props.webauthn && isWebauthnEnabled && (
          <>
            <styled.span>or</styled.span>

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
      <styled.p color="fg.error" fontSize="sm">
        {errors.token?.message}
      </styled.p>
      <Button type="submit" w="full" onClick={handlePassword}>
        Register
      </Button>
      <styled.p color="fg.error" fontSize="sm">
        {errors.root?.message}
      </styled.p>
    </styled.form>
  );
}
