"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi/accounts";
import { authPasswordUpdate } from "src/api/openapi/auth";
import { APIError } from "src/api/openapi/schemas";
import { PasswordSchema } from "src/screens/auth/schemas";
import { deriveError } from "src/utils/error";

const FormSchema = z.object({
  old: PasswordSchema,
  new: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function usePassword() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { push } = useRouter();
  const { mutate } = useAccountGet();

  async function handlePasswordChange(payload: Form) {
    await authPasswordUpdate(payload)
      .then(() => {
        push("/");
        mutate();
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  return {
    form: {
      register,
      handlePasswordChange: handleSubmit(handlePasswordChange),
      errors,
    },
  };
}
