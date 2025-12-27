import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { EditorModeSchema } from "@/lib/settings/editor";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";
import { SidebarDefaultStateSchema } from "@/lib/settings/sidebar";

export type Props = {
  settings: Settings;
};

export const FormSchema = z.object({
  editorMode: EditorModeSchema,
  sidebarDefaultState: SidebarDefaultStateSchema,
});
export type Form = z.infer<typeof FormSchema>;

export function useInterfaceSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      editorMode: settings.metadata.editor.mode,
      sidebarDefaultState: settings.metadata.sidebar.defaultState,
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
    formState: form.formState,
    onSubmit,
  };
}
