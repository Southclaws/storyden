"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { parseAsBoolean, useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useProfileGet } from "src/api/openapi-client/profiles";
import { ProfileGetOKResponse } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { handle } from "@/api/client";
import { useProfileMutations } from "@/lib/profile/mutation";

export type Props = {
  profile: ProfileGetOKResponse;
};

export const profileHandleRegex =
  /^(?!-)(?!.*--)[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  handle: z
    .string()
    .min(1, "Please enter a handle")
    .max(30, "Handle must be 30 characters or less")
    .regex(
      profileHandleRegex,
      "Handle can only contain lowercase letters, numbers, and hyphens",
    ),

  bio: z.string().max(200, "Bio must be 200 characters or less"),
});
export type Form = z.infer<typeof FormSchema>;

export function useProfileScreen({ profile }: Props) {
  const router = useRouter();
  const session = useSession();
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

  return {
    ready: true as const,
    form,
    state: {
      isEditing,
      isSelf,
    },
    data: {
      profile: data,
    },
    handlers: {
      handleSetEditing,

      handleSave,
    },
  };
}
