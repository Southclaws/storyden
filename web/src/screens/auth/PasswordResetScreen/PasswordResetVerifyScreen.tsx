"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { handle } from "@/api/client";
import { authPasswordReset } from "@/api/openapi-client/auth";
import { FormControl } from "@/components/ui/FormControl";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { Input } from "@/components/ui/input";
import { PasswordSchema } from "@/lib/auth/schemas";
import { styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

type Props = {
  token: string;
};

const FormSchema = z.object({
  password: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordResetVerifyScreen({ token }: Props) {
  const router = useRouter();

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        await authPasswordReset({
          token,
          new: payload.password,
        });
        router.push("/");
      },
      {
        promiseToast: {
          loading: "Resetting password...",
          success: "Password reset successfully.",
        },
      },
    );
  });

  return {
    form,
    handlers: {
      handleSubmit,
    },
  };
}

export function PasswordResetVerifyScreen(props: Props) {
  const { form, handlers } = usePasswordResetVerifyScreen(props);

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
          type="password"
          w="full"
          textAlign="center"
          placeholder="Your new password..."
          required
          {...form.register("password")}
        />
        <FormErrorText>{form.formState.errors["new"]?.message}</FormErrorText>
      </FormControl>

      <Button
        type="submit"
        size="sm"
        w="full"
        loading={form.formState.isSubmitting}
      >
        Reset
      </Button>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
    </styled.form>
  );
}
