import { zodResolver } from "@hookform/resolvers/zod";
import { intervalToDuration } from "date-fns";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { useI18n } from "@/i18n/provider";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";

export type Props = {
  settings: AdminSettings;
};

export const DEFAULT_RATE_LIMIT = 5000;
export const DEFAULT_RATE_LIMIT_PERIOD = 3600;
export const DEFAULT_RATE_LIMIT_BUCKET = 60;
export const DEFAULT_RATE_LIMIT_GUEST_COST = 1;
export const DEFAULT_CLIENT_IP_MODE = "remote_addr";
export const DEFAULT_CLIENT_IP_HEADER = "X-Real-IP";

export function formatSeconds(
  seconds: number | undefined,
  t?: (key: string) => string,
): string {
  const safeSeconds =
    typeof seconds === "number" && Number.isFinite(seconds) ? seconds : 0;
  const duration = intervalToDuration({ start: 0, end: safeSeconds * 1000 });
  const translate = t ?? ((key: string) => key);
  const units = [
    { value: duration.days, singular: "day", plural: "days" },
    { value: duration.hours, singular: "hour", plural: "hours" },
    { value: duration.minutes, singular: "minute", plural: "minutes" },
    { value: duration.seconds, singular: "second", plural: "seconds" },
  ];

  const parts = units
    .filter((unit) => unit.value && unit.value > 0)
    .map((unit) => {
      const value = unit.value ?? 0;
      return `${value} ${translate(value === 1 ? unit.singular : unit.plural)}`;
    });

  return parts.length > 0 ? parts.join(" ") : `0 ${translate("seconds")}`;
}

export const FormSchema = z.object({
  rate_limit: z.number().default(DEFAULT_RATE_LIMIT),
  rate_limit_period: z.number().default(DEFAULT_RATE_LIMIT_PERIOD),
  rate_limit_bucket: z.number().default(DEFAULT_RATE_LIMIT_BUCKET),
  rate_limit_guest_cost: z.number().default(DEFAULT_RATE_LIMIT_GUEST_COST),
  cost_overrides: z.record(z.string(), z.number()).default({}),
  client_ip_mode: z
    .enum(["remote_addr", "single_header", "xff_trusted_proxies"])
    .default("remote_addr"),
  client_ip_header: z.string().default(DEFAULT_CLIENT_IP_HEADER),
  trusted_proxy_cidrs: z.string().default(""),
});
export type Form = z.infer<typeof FormSchema>;

export function parseTrustedProxyCidrs(value: string): string[] {
  return value
    .split(/[\n,]/)
    .map((v) => v.trim())
    .filter(Boolean);
}

export function buildClientIPSettingsPayload(
  data: Pick<Form, "client_ip_mode" | "client_ip_header">,
  trustedProxyCidrs: string[],
) {
  const payload: {
    client_ip_mode: Form["client_ip_mode"];
    client_ip_header?: string;
    trusted_proxy_cidrs?: string[];
  } = {
    client_ip_mode: data.client_ip_mode,
  };

  if (data.client_ip_mode === "single_header") {
    payload.client_ip_header = data.client_ip_header;
  }

  if (data.client_ip_mode === "xff_trusted_proxies") {
    payload.trusted_proxy_cidrs = trustedProxyCidrs;
  }

  return payload;
}

export function useSystemSettings({ settings }: Props) {
  const { t } = useI18n();
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
      client_ip_mode:
        settings.services?.client_ip?.client_ip_mode ?? DEFAULT_CLIENT_IP_MODE,
      client_ip_header:
        settings.services?.client_ip?.client_ip_header ??
        DEFAULT_CLIENT_IP_HEADER,
      trusted_proxy_cidrs: (
        settings.services?.client_ip?.trusted_proxy_cidrs ?? []
      ).join(", "),
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    const trustedProxyCidrs = parseTrustedProxyCidrs(data.trusted_proxy_cidrs);
    const clientIP = buildClientIPSettingsPayload(data, trustedProxyCidrs);

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
            client_ip: clientIP,
          },
        });

        form.reset(data);
        await revalidate();
      },
      {
        promiseToast: {
          loading: t("Saving settings..."),
          success: t("Settings saved"),
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
