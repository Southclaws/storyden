"use client";

import { FormControl } from "@/components/ui/FormControl";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { useRegisterEmailForm } from "./useRegisterEmailForm";

export function RegisterEmailForm() {
  const { form, handlers } = useRegisterEmailForm();

  return (
    <styled.form
      className={vstack()}
      w="full"
      gap="2"
      textAlign="center"
      onSubmit={handlers.handleSubmit}
    >
      <FormControl>
        <Input
          type="email"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="email address"
          required
          {...form.register("email")}
        />
        <FormErrorText>{form.formState.errors["email"]?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <Input
          type="text"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="@username"
          required
          {...form.register("handle")}
        />
        <FormErrorText>
          {form.formState.errors["handle"]?.message}
        </FormErrorText>
      </FormControl>

      <FormControl>
        <Input
          type="password"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="password"
          autoComplete="new-password"
          {...form.register("password")}
        />

        <FormErrorText>
          {form.formState.errors["password"]?.message}
        </FormErrorText>
      </FormControl>

      <Button type="submit" w="full">
        Register
      </Button>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
    </styled.form>
  );
}
