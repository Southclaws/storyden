import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { handle } from "@/api/client";
import { EditorModeSchema } from "@/lib/settings/editor";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";
import * as z from "@/lib/validation";

export type Props = {
  settings: AdminSettings;
};

export const FormSchema = z.object({
  editorMode: EditorModeSchema,
  signaturesEnabled: z.boolean(),
  signatureMaxHeight: z.number().check(z.int(), z.minimum(32), z.maximum(2000)),
  signatureMaxChars: z.number().check(z.int(), z.minimum(1), z.maximum(10000)),
});
export type Form = z.infer<typeof FormSchema>;

export function useInterfaceSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      editorMode: settings.metadata.editor.mode,
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
    register: form.register,
    control: form.control,
    signaturesEnabled: form.watch("signaturesEnabled"),
    formState: form.formState,
    onSubmit,
  };
}
