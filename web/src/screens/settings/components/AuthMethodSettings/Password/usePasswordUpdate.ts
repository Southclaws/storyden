"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi-client/accounts";
import { authPasswordUpdate } from "src/api/openapi-client/auth";
import { APIError } from "src/api/openapi-schema";
import {
  ExistingPasswordSchema,
  PasswordSchema,
} from "src/screens/auth/schemas";
import { deriveError } from "src/utils/error";

const FormSchema = z.object({
  old: ExistingPasswordSchema,
  new: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordUpdate() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate } = useAccountGet();
  const [success, setSuccess] = useState(false);

  async function handlePasswordChange(payload: Form) {
    await authPasswordUpdate(payload)
      .then(() => {
        mutate();
        setSuccess(true);
      })
      .catch((e: APIError) => setError("root", { message: deriveError(e) }));
  }

  function handleCloseNotification() {
    setSuccess(false);
  }

  return {
    form: {
      register,
      handlePasswordChange: handleSubmit(handlePasswordChange),
      errors,
    },
    success,
    handleCloseNotification,
  };
}
