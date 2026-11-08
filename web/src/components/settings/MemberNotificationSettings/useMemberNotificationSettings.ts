import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, useWatch } from "react-hook-form";

import { handle } from "@/api/client";
import { accountUpdate, useAccountGet } from "@/api/openapi-client/accounts";
import {
  notificationSettingsUpdate,
  useNotificationSettingsGet,
} from "@/api/openapi-client/notifications";
import type {
  NotificationEvent,
  NotificationSettings,
} from "@/api/openapi-schema";
import type { Member } from "@/lib/settings/member-settings";
import * as z from "@/lib/validation";

const FormSchema = z.object({
  markReadOnClick: z.boolean(),
  pushEnabled: z.boolean(),
  events: z.record(
    z.string(),
    z.object({ inbox: z.boolean(), push: z.boolean() }),
  ),
});

type Form = z.infer<typeof FormSchema>;

export function useMemberNotificationSettings(
  session: Member,
  notificationSettings: NotificationSettings,
) {
  const { mutate } = useAccountGet();
  const { mutate: mutateSettings } = useNotificationSettingsGet();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      markReadOnClick: session.meta.notifications.markReadOnClick,
      pushEnabled: notificationSettings.push.enabled,
      events: notificationSettings.events,
    },
  });
  const events = useWatch({ control: form.control, name: "events" });
  const pushEnabled = useWatch({ control: form.control, name: "pushEnabled" });

  function setGroupEnabled(
    groupEvents: NotificationEvent[],
    channel: "inbox" | "push",
    enabled: boolean,
  ) {
    const overrides = { ...form.getValues("events") };

    for (const event of groupEvents) {
      const settings = {
        ...(overrides[event] ?? { inbox: false, push: false }),
        [channel]: enabled,
      };

      if (channel === "push" && enabled) {
        settings.inbox = true;
      }

      if (channel === "inbox" && !enabled) {
        settings.push = false;
      }

      overrides[event] = settings;
    }

    form.setValue("events", overrides, { shouldDirty: true });
  }

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const saved = await notificationSettingsUpdate({
          events: data.events,
          push: { enabled: data.pushEnabled },
        });
        await mutateSettings(saved, { revalidate: false });

        if (
          data.markReadOnClick !== session.meta.notifications.markReadOnClick
        ) {
          const account = await accountUpdate({
            meta: {
              ...session.meta,
              notifications: {
                ...session.meta.notifications,
                markReadOnClick: data.markReadOnClick,
              },
            },
          });

          await mutate(account, { revalidate: false });
        }

        form.reset(data);
      },
      {
        promiseToast: {
          loading: "Saving settings...",
          success: "Settings saved",
        },
      },
    );
  });

  return {
    control: form.control,
    formState: form.formState,
    events,
    pushEnabled,
    setGroupEnabled,
    onSubmit,
  };
}
