"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  accessKeyCreate,
  getAccessKeyListKey,
} from "@/api/openapi-client/auth";
import { AccessKeyCreateBody } from "@/api/openapi-schema";
import { useI18n } from "@/i18n/provider";
import { UseDisclosureProps } from "@/utils/useDisclosure";

const getFormSchema = (t: (key: string) => string) =>
  z.object({
    name: z.string().min(1, t("Name is required")).max(50, t("Name is too long")),
    expires_at: z.string().optional(),
  });
export const FormSchema = getFormSchema((key) => key);
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  onClose?: () => void;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useCreateAccessKeyScreen(props: WithDisclosure<Props>) {
  const { mutate } = useSWRConfig();
  const { t } = useI18n();

  const [createdSecret, setCreatedSecret] = useState<string | null>(null);

  const form = useForm<Form>({
    resolver: zodResolver(getFormSchema(t)),
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
          loading: t("Creating access key..."),
          success: t("Access key created successfully"),
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
