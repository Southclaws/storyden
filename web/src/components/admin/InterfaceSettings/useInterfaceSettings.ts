import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useI18n } from "@/i18n/provider";
import { EditorModeSchema } from "@/lib/settings/editor";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";
import { SidebarDefaultStateSchema } from "@/lib/settings/sidebar";

export type Props = {
  settings: AdminSettings;
};

export const FormSchema = z.object({
  editorMode: EditorModeSchema,
  sidebarDefaultState: SidebarDefaultStateSchema,
  signaturesEnabled: z.boolean(),
  signatureMaxHeight: z.number().int().min(32).max(2000),
  signatureMaxChars: z.number().int().min(1).max(10000),
});
export type Form = z.infer<typeof FormSchema>;

export function useInterfaceSettings({ settings }: Props) {
  const { t } = useI18n();
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      editorMode: settings.metadata.editor.mode,
      sidebarDefaultState: settings.metadata.sidebar.defaultState,
      signaturesEnabled: settings.metadata.signatures.enabled,
      signatureMaxHeight: settings.metadata.signatures.maxHeight,
      signatureMaxChars:
        settings.services?.moderation?.signature_length_max ?? 500,
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    handle(
      async () => {
        await updateSettings({
          metadata: {
            ...settings.metadata,
            editor: {
              mode: data.editorMode,
            },
            sidebar: {
              defaultState: data.sidebarDefaultState,
            },
            signatures: {
              enabled: data.signaturesEnabled,
              maxHeight: data.signatureMaxHeight,
            },
          },
          services: {
            ...(settings.services ?? {}),
            moderation: {
              ...(settings.services?.moderation ?? {}),
              signature_length_max: data.signatureMaxChars,
            },
          },
        });
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
    register: form.register,
    control: form.control,
    signaturesEnabled: form.watch("signaturesEnabled"),
    formState: form.formState,
    onSubmit,
  };
}
