"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { handle } from "@/api/client";
import { authPasswordResetRequestEmail } from "@/api/openapi-client/auth";
import { FormControl } from "@/components/ui/FormControl";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { WEB_ADDRESS } from "@/config";
import { styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";
import { deriveError } from "@/utils/error";

const FormSchema = z.object({
  email: z.string().email(),
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordResetEmailScreen() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    setError(null);

    await handle(
      async () => {
        await authPasswordResetRequestEmail({
          email: payload.email,
          token_url: {
            url: `${WEB_ADDRESS}/password-reset/verify`,
            query: "token",
          },
        });
        router.push("/password-reset/verify");
      },
      {
        errorToast: false,
        onError: async (error) => {
          setError(deriveError(error));
        },
      },
    );
  });

  return {
    error,
    setError,
    form,
    handlers: {
      handleSubmit,
    },
  };
}

export function PasswordResetEmailScreen() {
  const { error, setError, form, handlers } = usePasswordResetEmailScreen();

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
          textAlign="center"
          placeholder="Email address..."
          required
          {...form.register("email")}
        />
        <FormErrorText>{form.formState.errors["email"]?.message}</FormErrorText>
      </FormControl>

      <Button type="submit" size="sm" w="full">
        Reset
      </Button>

      <Admonition
        value={Boolean(error)}
        onChange={() => setError(null)}
        kind="failure"
        title="Unable to send reset email"
        textAlign="start"
      >
        {error}
      </Admonition>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
    </styled.form>
  );
}
