"use client";

import { FormControl } from "@/components/ui/FormControl";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { useLoginEmailForm } from "./useLoginEmailForm";

export function LoginEmailForm() {
  const { form, handlers } = useLoginEmailForm();

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
          type="text"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="username or email address"
          required
          {...form.register("identifier")}
        />
        <FormErrorText>
          {form.formState.errors["identifier"]?.message}
        </FormErrorText>
      </FormControl>

      <FormControl>
        <Input
          type="password"
          w="full"
          size="sm"
          textAlign="center"
          placeholder="password"
          autoComplete="current-password"
          {...form.register("password")}
        />

        <FormErrorText>
          {form.formState.errors["password"]?.message}
        </FormErrorText>
      </FormControl>

      <Button type="submit" w="full">
        Login
      </Button>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
    </styled.form>
  );
}
