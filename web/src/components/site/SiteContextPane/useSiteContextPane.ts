"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
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

export const EditingSchema = z.preprocess(
  (v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  },
  z.enum(["settings", "feed"]),
);
export type Editing = z.infer<typeof EditingSchema>;

export function useSiteContextPane({ session, initialSettings }: Props) {
  session = useSession(session);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: initialSettings,
  });

  const { revalidate, updateSettings } = useSettingsMutation();

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
      settings,
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
