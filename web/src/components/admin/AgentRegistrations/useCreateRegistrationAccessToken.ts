"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  adminOAuthDynamicRegistrationAccessTokenCreate,
  getAdminOAuthDynamicRegistrationAccessTokenListKey,
} from "@/api/openapi-client/admin";
import { deriveError } from "@/utils/error";

export const FormSchema = z.object({
  label: z
    .string()
    .trim()
    .min(1, "Label is required")
    .max(200, "Label is too long"),
  lifetime: z.enum(["1", "24", "168", "720"]),
  max_registrations: z
    .number()
    .int()
    .min(1, "Allow at least one registration")
    .max(2147483647),
});

export function useCreateRegistrationAccessToken() {
  const { mutate } = useSWRConfig();
  const [createdSecret, setCreatedSecret] = useState<string>();
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: { label: "", lifetime: "24", max_registrations: 1 },
  });
  const handleSubmit = form.handleSubmit(async (data) => {
    form.clearErrors("root");
    await handle(
      async () => {
        const issued = await adminOAuthDynamicRegistrationAccessTokenCreate({
          label: data.label,
          expires_at: new Date(
            Date.now() + Number(data.lifetime) * 3600000,
          ).toISOString(),
          max_registrations: data.max_registrations,
        });
        setCreatedSecret(issued.token);
        await mutate(
          (key) =>
            Array.isArray(key) &&
            key[0] === getAdminOAuthDynamicRegistrationAccessTokenListKey()[0],
        );
      },
      {
        errorToast: false,
        onError: async (error) => {
          form.setError("root", { message: deriveError(error) });
        },
      },
    );
  });

  return { form, handleSubmit, createdSecret };
}
