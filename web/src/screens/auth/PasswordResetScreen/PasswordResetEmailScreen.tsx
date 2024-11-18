"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { handle } from "@/api/client";
import { authPasswordResetRequestEmail } from "@/api/openapi-client/auth";
import { FormControl } from "@/components/ui/FormControl";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { WEB_ADDRESS } from "@/config";
import { styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

const FormSchema = z.object({
  email: z.string().email(),
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordResetEmailScreen() {
  const router = useRouter();

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(async () => {
      await authPasswordResetRequestEmail({
        email: payload.email,
        token_url: {
          url: `${WEB_ADDRESS}/password-reset/verify`,
          query: "token",
        },
      });
      router.push("/password-reset/verify");
    });
  });

  return {
    form,
    handlers: {
      handleSubmit,
    },
  };
}

export function PasswordResetEmailScreen() {
  const { form, handlers } = usePasswordResetEmailScreen();

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
        <FormErrorText>
          {form.formState.errors["identifier"]?.message}
        </FormErrorText>
      </FormControl>

      <Button type="submit" size="sm" w="full">
        Reset
      </Button>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
    </styled.form>
  );
}
