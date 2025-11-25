"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { useSiteEditing } from "@/lib/site-editing/useSiteEditing";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { getIconURL } from "@/utils/icon";
import { hasPermission } from "@/utils/permissions";

export type Props = {
  session: Account | undefined;
  initialSettings: Settings;
};

export const FormSchema = z.object({
  title: z.string().min(1, "Please write a title."),
  description: z.string().min(1, "Please write a short description."),
  content: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export function useSiteContextPane({ session, initialSettings }: Props) {
  session = useSession(session);
  const { editing, isEditingSettings, toggleSettingsEditing, stopEditing } =
    useSiteEditing(session);

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: initialSettings,
  });

  const { revalidate, updateSettings } = useSettingsMutation(initialSettings);

  const { ready, error, settings } = useSettings(initialSettings);
  if (!ready) {
    return {
      ready: false as const,
      error,
    };
  }

  const iconURL = getIconURL("512x512");

  const isEditingEnabled = hasPermission(session, Permission.MANAGE_SETTINGS);
  const isAdmin = hasPermission(session, Permission.ADMINISTRATOR);

  const handleSaveSettings = form.handleSubmit(async (value) => {
    await handle(
      async () => {
        await updateSettings(value);
        stopEditing();
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
      settings,
      iconURL,
      isEditingEnabled,
      isAdmin,
      editing,
    },
    handlers: {
      handleEnableEditing: toggleSettingsEditing,
      handleSaveSettings,
    },
  };
}
