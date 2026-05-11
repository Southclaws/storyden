"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryState } from "nuqs";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { Account, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Editing, EditingSchema } from "@/components/site/editing";
import type { Translate } from "@/i18n/format";
import { useI18n } from "@/i18n/provider";
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
  title: z.string().min(1),
  description: z.string().min(1),
  content: z.string().optional(),
});

function getFormSchema(t: Translate) {
  return z.object({
    title: z.string().min(1, t("Please write a title.")),
    description: z.string().min(1, t("Please write a short description.")),
    content: z.string().optional(),
  });
}

export type Form = z.infer<typeof FormSchema>;

export function useSiteContextPane({ session, initialSettings }: Props) {
  const { t } = useI18n();
  session = useSession(session);
  const [editing, setEditing] = useQueryState<null | Editing>("editing", {
    defaultValue: null,
    clearOnDefault: true,
    parse: EditingSchema.parse,
  });

  const form = useForm<Form>({
    resolver: zodResolver(getFormSchema(t)),
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
          loading: t("Saving settings..."),
          success: t("Settings saved"),
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
