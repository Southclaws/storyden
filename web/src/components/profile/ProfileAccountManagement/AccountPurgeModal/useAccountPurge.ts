import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { moderationActionCreate } from "@/api/openapi-client/admin";
import { ModerationActionPurgeAccountContentType } from "@/api/openapi-schema";

export type Props = {
  accountId: string;
  handle: string;
  onSave?: () => void;
};

const schema = z.object({
  contentTypes: z
    .array(z.nativeEnum(ModerationActionPurgeAccountContentType))
    .min(1, "Select at least one content type to purge"),
});

type FormValues = z.infer<typeof schema>;

export function useAccountPurgeScreen({ accountId, onSave }: Props) {
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      contentTypes: [],
    },
  });

  async function handlePurge(data: FormValues) {
    await handle(
      async () => {
        await moderationActionCreate({
          action: "purge_account",
          account_id: accountId,
          include: data.contentTypes,
        });

        onSave?.();
      },
      {
        promiseToast: {
          loading: "Purging account content...",
          success: "Account content purged successfully",
        },
      },
    );
  }

  return {
    form,
    handlers: {
      handlePurge: form.handleSubmit(handlePurge),
    },
  };
}
