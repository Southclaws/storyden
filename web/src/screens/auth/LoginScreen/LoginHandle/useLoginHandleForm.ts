"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { authPasswordSignin } from "@/api/openapi-client/auth";
import { APIError } from "@/api/openapi-schema";
import { ExistingPasswordSchema, UsernameSchema } from "@/lib/auth/schemas";
import { deriveError } from "@/utils/error";

const FormSchema = z.object({
  identifier: UsernameSchema,
  token: z.string().optional(), // Validated properly during password submission
});
const FormPasswordSchema = z.object({
  identifier: UsernameSchema,
  token: ExistingPasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useLoginHandleForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { push } = useRouter();
  const searchParams = useSearchParams();
  const returnURL = searchParams.get("return_url") ?? "/";
  const { mutate } = useAccountGet();

  async function handlePassword(payload: Form) {
    const parsed = FormPasswordSchema.safeParse(payload);
    if (!parsed.success) {
      const fieldErrors = parsed.error.flatten().fieldErrors;

      if (fieldErrors.identifier) {
        setError("identifier", {
          message: fieldErrors.identifier.join(", "),
        });
      }

      if (fieldErrors.token) {
        setError("token", {
          message: fieldErrors.token.join(", "),
        });
      }

      return;
    }

    await authPasswordSignin(parsed.data)
      .then(() => {
        push(returnURL);
        mutate();
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  return {
    form: {
      register,
      handlePassword: handleSubmit(handlePassword),
      errors,
    },
  };
}
