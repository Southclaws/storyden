"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountAuthProviderList } from "src/api/openapi-client/accounts";
import { authPasswordCreate } from "src/api/openapi-client/auth";
import { APIError } from "src/api/openapi-schema";
import { deriveError } from "src/utils/error";

import { PasswordSchema } from "@/lib/auth/schemas";

const FormSchema = z.object({
  newPassword: PasswordSchema,
  confirmPassword: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordCreate() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate } = useAccountAuthProviderList();
  const [success, setSuccess] = useState(false);

  async function handlePasswordCreate(payload: Form) {
    if (payload.newPassword !== payload.confirmPassword) {
      setError("confirmPassword", {
        message: "Passwords do not match",
      });
      return;
    }

    await authPasswordCreate({
      password: payload.newPassword,
    })
      .then(() => {
        setSuccess(true);
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  function handleCloseNotification() {
    setSuccess(false);
    mutate();
  }

  return {
    form: {
      register,
      handlePasswordCreate: handleSubmit(handlePasswordCreate),
      errors,
    },
    success,
    handleCloseNotification,
  };
}
