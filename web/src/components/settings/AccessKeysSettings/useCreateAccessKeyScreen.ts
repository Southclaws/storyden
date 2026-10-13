"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  accessKeyCreate,
  getAccessKeyListKey,
} from "@/api/openapi-client/auth";
import { AccessKeyCreateBody } from "@/api/openapi-schema";
import * as z from "@/lib/validation";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export const FormSchema = z.object({
  name: z
    .string()
    .check(
      z.minLength(1, "Name is required"),
      z.maxLength(50, "Name is too long"),
    ),
  expires_at: z.optional(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  onClose?: () => void;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useCreateAccessKeyScreen(props: WithDisclosure<Props>) {
  const { mutate } = useSWRConfig();

  const [createdSecret, setCreatedSecret] = useState<string | null>(null);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: "",
      expires_at: "",
    },
  });

  const handleSubmit = form.handleSubmit(async (data) => {
    const body: AccessKeyCreateBody = {
      name: data.name,
      ...(data.expires_at && { expires_at: data.expires_at }),
    };

    await handle(
      async () => {
        const result = await accessKeyCreate(body);

        setCreatedSecret(result.secret);

        // Refresh the access key list
        await mutate(getAccessKeyListKey());
      },
      {
        promiseToast: {
          loading: "Creating access key...",
          success: "Access key created successfully",
        },
      },
    );
  });

  return {
    form,
    handleSubmit,
    createdSecret,
  };
}
