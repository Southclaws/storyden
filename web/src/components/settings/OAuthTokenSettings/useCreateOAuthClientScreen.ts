"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useSWRConfig } from "swr";

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
import * as z from "@/lib/validation";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export type OAuthClientPreset = "machine" | "app_integration" | "public_app";

export const FormSchema = z.object({
  name: z
    .string()
    .check(
      z.minLength(1, "Name is required"),
      z.maxLength(80, "Name is too long"),
    ),
  preset: z.enum(["machine", "app_integration", "public_app"]),
  redirectUris: z.optional(
    z.array(
      z.object({ value: z.string().check(z.url("Must be a valid URL")) }),
    ),
  ),
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
      preset: "app_integration",
      redirectUris: [],
      permissions: [],
    },
  });

  const handleSubmit = form.handleSubmit(async (data) => {
    const { type, allowedGrants, pkceRequired } = getPresetConfig(data.preset);

    const body: OAuthClientSelfCreateBody = {
      name: data.name,
      type,
      allowed_scopes: Array.from(new Set(data.permissions)),
      allowed_grants: allowedGrants,
      redirect_uris:
        data.preset === "machine"
          ? undefined
          : data.redirectUris?.map((uri) => uri.value).filter(Boolean),
      pkce_required: pkceRequired,
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

  function getPresetConfig(preset: OAuthClientPreset) {
    switch (preset) {
      case "machine":
        return {
          type: "confidential" as const,
          allowedGrants: ["client_credentials"],
          pkceRequired: false,
        };
      case "app_integration":
        return {
          type: "confidential" as const,
          allowedGrants: ["authorization_code", "refresh_token"],
          pkceRequired: true,
        };
      case "public_app":
        return {
          type: "public" as const,
          allowedGrants: ["authorization_code", "refresh_token"],
          pkceRequired: true,
        };
    }
  }

  return {
    form,
    createdClient,
    handleSubmit,
    onClose: props.onClose,
  };
}
