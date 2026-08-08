"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useProfileGet } from "@/api/openapi-client/profiles";
import {
  Account,
  AccountMutableProps,
  ProfileGetOKResponse,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";

import { handle } from "@/api/client";
import { useProfileMutations } from "@/lib/profile/mutation";
import { UsernameSchema } from "@/lib/auth/schemas";
import type { SignatureConfig } from "@/lib/settings/settings";
import { hasPermissionOr } from "@/utils/permissions";

export type Props = {
  initialSession?: Account;
  profile: ProfileGetOKResponse;
  initialSignatureConfig: SignatureConfig;
};

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  handle: UsernameSchema,

  bio: z.string(),
  signature: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useProfileScreen({
  initialSession,
  profile,
  initialSignatureConfig,
}: Props) {
  const router = useRouter();
  const session = useSession(initialSession);
  const signaturesEnabled = initialSignatureConfig.enabled;
  const [isEditing, setEditing] = useQueryState("edit", {
    ...parseAsBoolean,
    defaultValue: false,
    clearOnDefault: true,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: profile.name,
      handle: profile.handle,
      bio: profile.bio,
      signature: profile.signature ?? "",
    },
  });

  const { update, revalidate } = useProfileMutations(profile.handle);

  const { data, error } = useProfileGet(profile.handle, {
    swr: { fallbackData: profile },
  });
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  function handleSetEditing() {
    setEditing(true);
  }

  const handleSave = form.handleSubmit(async (data) => {
    const payload: AccountMutableProps = signaturesEnabled
      ? data
      : { ...data, signature: undefined };

    await handle(
      async () => {
        await update(payload);

        if (data.handle !== profile.handle) {
          router.replace(`/m/${data.handle}`);
        } else {
          setEditing(false);
        }
      },
      {
        cleanup: async () => await revalidate(),
        promiseToast: {
          loading: "Updating profile...",
          success: "Profile updated",
        },
      },
    );
  });

  const isSelf = session?.id === data.id;
  const canViewAccount = hasPermissionOr(
    session,
    () => isSelf,
    "VIEW_ACCOUNTS",
  );

  return {
    ready: true as const,
    form,
    state: {
      isEditing,
      isSelf,
      canViewAccount,
      signaturesEnabled,
    },
    data: {
      session,
      profile: data,
    },
    handlers: {
      handleSetEditing,
      handleSave,
    },
  };
}
