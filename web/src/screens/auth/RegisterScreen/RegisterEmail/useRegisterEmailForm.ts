"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { handle } from "@/api/client";
import { useAccountGet } from "@/api/openapi-client/accounts";
import { authEmailPasswordSignup } from "@/api/openapi-client/auth";
import { PasswordSchema, UsernameSchema } from "@/lib/auth/schemas";
import { refreshFeed } from "@/lib/feed/refresh";
import { deriveError } from "@/utils/error";

const FormSchema = z.object({
  handle: UsernameSchema,
  email: z.string().email(),
  password: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useRegisterEmailForm() {
  const router = useRouter();

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate } = useAccountGet();

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        await authEmailPasswordSignup(payload);
        mutate();
        refreshFeed();
        router.push("/");
      },
      {
        errorToast: false,
        onError: async (error) => {
          form.setError("root", { message: deriveError(error) });
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
