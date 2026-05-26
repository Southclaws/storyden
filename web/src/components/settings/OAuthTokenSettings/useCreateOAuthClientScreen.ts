"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  getOAuthClientListKey,
  oAuthClientCreate,
} from "@/api/openapi-client/auth";
import {
  OAuthClientIssued,
  OAuthClientSelfCreateBody,
} from "@/api/openapi-schema";
import { PermissionSchema } from "@/lib/permission/permission";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export const FormSchema = z.object({
  name: z.string().min(1, "Name is required").max(80, "Name is too long"),
  permissions: z.array(PermissionSchema),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  onClose?: () => void;
};

export type WithDisclosure<T> = UseDisclosureProps & T;

export function useCreateOAuthClientScreen(props: WithDisclosure<Props>) {
  const { mutate } = useSWRConfig();
  const [createdClient, setCreatedClient] = useState<OAuthClientIssued | null>(
    null,
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: "",
      permissions: [],
    },
  });

  const handleSubmit = form.handleSubmit(async (data) => {
    const body: OAuthClientSelfCreateBody = {
      name: data.name,
      allowed_scopes: Array.from(new Set(data.permissions)),
    };

    await handle(
      async () => {
        const result = await oAuthClientCreate(body);

        setCreatedClient(result);
        await mutate(getOAuthClientListKey());
      },
      {
        promiseToast: {
          loading: "Creating OAuth client...",
          success: "OAuth client created",
        },
      },
    );
  });

  return {
    form,
    createdClient,
    handleSubmit,
    onClose: props.onClose,
  };
}
