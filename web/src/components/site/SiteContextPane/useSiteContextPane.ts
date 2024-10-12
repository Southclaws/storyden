"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useGetInfo } from "@/api/openapi-client/misc";
import { Account, Info, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useInfoMutation } from "@/lib/settings/mutation";
import { getIconURL } from "@/utils/icon";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  session: Account | undefined;
  info: Info;
};

export const FormSchema = z.object({
  title: z.string().min(1, "Please write a title."),
  description: z.string().min(1, "Please write a short description."),
  content: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export const EditingSchema = z.preprocess(
  (v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  },
  z.enum(["settings"]),
);
export type Editing = z.infer<typeof EditingSchema>;

export function useSiteContextPane({ session, info }: Props) {
  session = useSession(session);
  const [editing, setEditing] = useQueryState<null | "settings">("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: info,
  });

  const { revalidate, updateSettings } = useInfoMutation(info);

  const { data, error, isValidating } = useGetInfo({
    swr: { fallbackData: info },
  });
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const iconURL = getIconURL("512x512");

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);
  const isAdmin = hasPermission(session, Permission.ADMINISTRATOR);

  function handleEnableEditing() {
    setEditing("settings");
  }

  const handleSaveSettings = form.handleSubmit(async (value) => {
    await handle(
      async () => {
        await updateSettings(value);
        setEditing(null);
      },
      {
        promiseToast: {
          loading: "Saving settings...",
          success: "Settings saved",
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  });

  return {
    ready: true as const,
    form,
    data: {
      info: data,
      iconURL,
      isEditingEnabled,
      isAdmin,
      editing,
    },
    handlers: {
      handleEnableEditing,
      handleSaveSettings,
    },
  };
}
