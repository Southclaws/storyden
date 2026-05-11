import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { moderationActionCreate } from "@/api/openapi-client/admin";
import { ModerationActionPurgeAccountContentType } from "@/api/openapi-schema";
import { useI18n } from "@/i18n/provider";

export type Props = {
  accountId: string;
  handle: string;
  onSave?: () => void;
};

function getSchema(t: (key: string) => string) {
  return z.object({
    contentTypes: z
      .array(z.nativeEnum(ModerationActionPurgeAccountContentType))
      .min(1, t("Select at least one content type to purge")),
  });
}

type FormValues = z.infer<ReturnType<typeof getSchema>>;

export function useAccountPurgeScreen({ accountId, onSave }: Props) {
  const { t } = useI18n();
  const form = useForm<FormValues>({
    resolver: zodResolver(getSchema(t)),
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
          loading: t("Purging account content..."),
          success: t("Account content purged successfully"),
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
