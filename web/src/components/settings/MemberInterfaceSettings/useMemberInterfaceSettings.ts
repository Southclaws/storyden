import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useSession } from "@/auth";
import { useProfileMutations } from "@/lib/profile/mutation";
import { Member } from "@/lib/settings/member-settings";

export const FormSchema = z.object({
  editorMode: z.enum(["richtext", "markdown"]),
  sidebarDefaultState: z.enum(["open", "closed"]),
});
export type Form = z.infer<typeof FormSchema>;

export type Props = {
  session: Member;
};

export function useMemberInterfaceSettings({ session }: Props) {
  const { update, revalidate } = useProfileMutations(session.handle);
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      editorMode: session.meta.editor.mode,
      sidebarDefaultState: session.meta.sidebar.defaultState,
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
    ready: true as const,
    control: form.control,
    formState: form.formState,
    onSubmit,
    session,
  };
}
