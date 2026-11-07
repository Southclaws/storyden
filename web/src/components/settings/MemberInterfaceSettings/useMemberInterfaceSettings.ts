import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { handle } from "@/api/client";
import { useAccountGet } from "@/api/openapi-client/accounts";
import { useProfileMutations } from "@/lib/profile/mutation";
import { Member } from "@/lib/settings/member-settings";
import * as z from "@/lib/validation";

export const FormSchema = z.object({
  editorMode: z.enum(["richtext", "markdown"]),
  markReadOnClick: z.boolean(),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  session: Member;
};

export function useMemberInterfaceSettings({ session }: Props) {
  const { update, revalidate } = useProfileMutations(session.handle);
  const { mutate: revalidateAccount } = useAccountGet();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      editorMode: session.meta.editor.mode,
      markReadOnClick: session.meta.notifications.markReadOnClick,
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await update({
          meta: {
            ...session.meta,
            editor: {
              mode: data.editorMode,
            },
            notifications: {
              ...session.meta.notifications,
              markReadOnClick: data.markReadOnClick,
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
          await Promise.all([revalidate(), revalidateAccount()]);
        },
      },
    );
  });

  return {
    ready: true as const,
    control: form.control,
    formState: form.formState,
    onSubmit,
    session,
  };
}
