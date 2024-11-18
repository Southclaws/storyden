"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { useAccountGet } from "src/api/openapi-client/accounts";
import { authPasswordUpdate } from "src/api/openapi-client/auth";
import { APIError } from "src/api/openapi-schema";
import { deriveError } from "src/utils/error";

import { handle } from "@/api/client";
import { ExistingPasswordSchema, PasswordSchema } from "@/lib/auth/schemas";

const FormSchema = z.object({
  old: ExistingPasswordSchema,
  new: PasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function usePasswordUpdate() {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate } = useAccountGet();

  const handlePasswordChange = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        await authPasswordUpdate(payload);
        await mutate();
      },
      {
        errorToast: false,
        async cleanup() {
          await mutate();
        },
        async onError(error: unknown) {
          form.setError("old", {
            type: "manual",
            message: deriveError(error),
          });
        },
      },
    );
  });

  return {
    form,
    handlePasswordChange,
  };
}
