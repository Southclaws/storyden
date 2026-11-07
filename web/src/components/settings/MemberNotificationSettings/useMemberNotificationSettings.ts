import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, useWatch } from "react-hook-form";

import { handle } from "@/api/client";
import { accountUpdate, useAccountGet } from "@/api/openapi-client/accounts";
import type { NotificationEvent } from "@/api/openapi-schema";
import type { Member } from "@/lib/settings/member-settings";
import * as z from "@/lib/validation";

const FormSchema = z.object({
  markReadOnClick: z.boolean(),
  pushEnabled: z.boolean(),
  events: z.record(z.string(), z.boolean()),
});

type Form = z.infer<typeof FormSchema>;

export function useMemberNotificationSettings(session: Member) {
  const { mutate } = useAccountGet();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      markReadOnClick: session.meta.notifications.markReadOnClick,
      pushEnabled: session.notification_settings?.push.enabled ?? true,
      events: session.notification_settings?.push.events ?? {},
    },
  });
  const events = useWatch({ control: form.control, name: "events" });
  const pushEnabled = useWatch({ control: form.control, name: "pushEnabled" });

  function setGroupEnabled(groupEvents: NotificationEvent[], enabled: boolean) {
    const overrides = { ...form.getValues("events") };

    for (const event of groupEvents) {
      overrides[event] = enabled;
    }

    form.setValue("events", overrides, { shouldDirty: true });
  }

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const account = await accountUpdate({
          meta: {
            ...session.meta,
            notifications: {
              ...session.meta.notifications,
              markReadOnClick: data.markReadOnClick,
            },
          },
          notification_settings: {
            ...session.notification_settings,
            push: {
              enabled: data.pushEnabled,
              events: data.events,
            },
          },
        });

        await mutate(account, { revalidate: false });
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
