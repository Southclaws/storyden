import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";

export type Props = {
  settings: AdminSettings;
};

export const FormSchema = z.object({
  threadBodyMaxSize: z.number().min(0).max(1_000_000),
  replyBodyMaxSize: z.number().min(0).max(1_000_000),
  wordBlockList: z.array(z.string()),
  wordReportList: z.array(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export function useModerationSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      threadBodyMaxSize:
        settings.services?.moderation?.thread_body_length_max ?? 60_000,
      replyBodyMaxSize:
        settings.services?.moderation?.reply_body_length_max ?? 10_000,
      wordBlockList: settings.services?.moderation?.word_block_list ?? [],
      wordReportList: settings.services?.moderation?.word_report_list ?? [],
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateSettings({
          services: {
            moderation: {
              thread_body_length_max: data.threadBodyMaxSize,
              reply_body_length_max: data.replyBodyMaxSize,
              word_block_list: data.wordBlockList,
              word_report_list: data.wordReportList,
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
