import { zodResolver } from "@hookform/resolvers/zod";
import { formatDuration, intervalToDuration } from "date-fns";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";

export type Props = {
  settings: AdminSettings;
};

export const DEFAULT_RATE_LIMIT = 5000;
export const DEFAULT_RATE_LIMIT_PERIOD = 3600;
export const DEFAULT_RATE_LIMIT_BUCKET = 60;
export const DEFAULT_RATE_LIMIT_GUEST_COST = 1;

export function formatSeconds(seconds: number): string {
  const duration = intervalToDuration({ start: 0, end: seconds * 1000 });
  return formatDuration(duration, {
    format: ["days", "hours", "minutes", "seconds"],
  });
}

export const FormSchema = z.object({
  rate_limit: z.number().default(DEFAULT_RATE_LIMIT),
  rate_limit_period: z.number().default(DEFAULT_RATE_LIMIT_PERIOD),
  rate_limit_bucket: z.number().default(DEFAULT_RATE_LIMIT_BUCKET),
  rate_limit_guest_cost: z.number().default(DEFAULT_RATE_LIMIT_GUEST_COST),
  cost_overrides: z.record(z.string(), z.number()).default({}),
});
export type Form = z.infer<typeof FormSchema>;

export function useSystemSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      rate_limit:
        settings.services?.rate_limiting?.rate_limit ?? DEFAULT_RATE_LIMIT,
      rate_limit_bucket:
        settings.services?.rate_limiting?.rate_limit_bucket ??
        DEFAULT_RATE_LIMIT_BUCKET,
      rate_limit_period:
        settings.services?.rate_limiting?.rate_limit_period ??
        DEFAULT_RATE_LIMIT_PERIOD,
      rate_limit_guest_cost:
        settings.services?.rate_limiting?.rate_limit_guest_cost ??
        DEFAULT_RATE_LIMIT_GUEST_COST,
      cost_overrides: settings.services?.rate_limiting?.cost_overrides ?? {},
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateSettings({
          services: {
            ...settings.services,
            rate_limiting: {
              rate_limit: data.rate_limit,
              rate_limit_bucket: data.rate_limit_bucket,
              rate_limit_period: data.rate_limit_period,
              rate_limit_guest_cost: data.rate_limit_guest_cost,
              cost_overrides: data.cost_overrides,
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
