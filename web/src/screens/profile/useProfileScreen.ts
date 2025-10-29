"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useProfileGet } from "src/api/openapi-client/profiles";
import { Account, ProfileGetOKResponse } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { handle } from "@/api/client";
import { useProfileMutations } from "@/lib/profile/mutation";
import { hasPermission, hasPermissionOr } from "@/utils/permissions";
import { isSlug } from "@/utils/slugify";

export type Props = {
  initialSession?: Account;
  profile: ProfileGetOKResponse;
};

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  handle: z
    .string()
    .min(1, "Please enter a handle")
    .max(30, "Handle must be 30 characters or less")
    .refine(isSlug, {
      message: "Handle can only contain letters, numbers, hyphens, and underscores",
    }),

  bio: z.string().max(200, "Bio must be 200 characters or less"),
});
export type Form = z.infer<typeof FormSchema>;

export function useProfileScreen({ initialSession, profile }: Props) {
  const router = useRouter();
  const session = useSession(initialSession);
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
    await handle(
      async () => {
        await update(data);

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
