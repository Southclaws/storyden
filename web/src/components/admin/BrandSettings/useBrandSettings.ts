import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { iconUpload } from "src/api/openapi-client/misc";
import { getColourVariants } from "src/utils/colour";

import { handle } from "@/api/client";
import { Translate } from "@/i18n/format";
import { useI18n } from "@/i18n/provider";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { MotdAlertTypeSchema, Settings } from "@/lib/settings/settings";
import { getIconURL } from "@/utils/icon";

export type Props = {
  settings: Settings;
};

export const FormSchema = z.object({
  title: z.string(),
  description: z.string(),
  content: z.string().optional(),
  accentColour: z.string(),
  motdContent: z.string().optional(),
  motdStartAt: z.string().optional(),
  motdEndAt: z.string().optional(),
  motdType: MotdAlertTypeSchema.optional(),
});

function getFormSchema(t: Translate) {
  return FormSchema.refine(
    (data) => {
      if (!data.motdStartAt || !data.motdEndAt) {
        return true;
      }

      const start = new Date(data.motdStartAt).getTime();
      const end = new Date(data.motdEndAt).getTime();
      if (Number.isNaN(start) || Number.isNaN(end)) {
        return true;
      }

      return start <= end;
    },
    {
      message: t("MOTD end date must be after start date."),
      path: ["motdEndAt"],
    },
  );
}
export type Form = z.infer<typeof FormSchema>;

export function useBrandSettings({ settings }: Props) {
  const { t } = useI18n();
  const { revalidate, updateSettings } = useSettingsMutation();
  const [motdContentInitialValue, setMotdContentInitialValue] = useState<
    string | undefined
  >(settings.motd?.content);
  const [motdContentResetKey, setMotdContentResetKey] = useState<string>();
  const form = useForm<Form>({
    resolver: zodResolver(getFormSchema(t)),
    defaultValues: {
      title: settings.title,
      description: settings.description,
      content: settings.content,
      accentColour: settings.accent_colour,
      motdContent: settings.motd?.content,
      motdStartAt: settings.motd?.start_at,
      motdEndAt: settings.motd?.end_at,
      motdType: settings.motd?.metadata?.type,
    },
  });
  const [currentIcon, setCurrentIcon] = useState<File | undefined>(undefined);
  const [contrast, setContrast] = useState(1);

  useEffect(() => {
    (async () => {
      const icon = await fetch(getIconURL("512x512"));

      const blob = await icon.blob();

      const file = new File([blob], "icon-512x512");

      setCurrentIcon(file);
    })();
  }, []);

  const updateColour = (colour: string) => {
    try {
      const cv = getColourVariants(colour);

      Object.entries(cv).forEach((property) =>
        document.documentElement.style.setProperty(property[0], property[1]),
      );
    } catch (e) {
      // NOTE: do nothing on invalid colours.
      console.warn("failed to update colour variable for previews", e);
    }
  };

  const onSubmit = form.handleSubmit(async (data) => {
    handle(
      async () => {
        const motdType = data["motdType"];
        updateColour(data["accentColour"]);
        const hasMotd =
          Boolean(data["motdContent"]?.trim()) ||
          Boolean(data["motdStartAt"]) ||
          Boolean(data["motdEndAt"]) ||
          Boolean(motdType);

        await updateSettings({
          title: data["title"],
          description: data["description"],
          content: data["content"],
          accent_colour: data["accentColour"],
          motd: hasMotd
            ? {
                content: data["motdContent"],
                start_at: toISODateTime(data["motdStartAt"]),
                end_at: toISODateTime(data["motdEndAt"]),
                metadata: motdType ? { type: motdType } : undefined,
              }
            : {},
        });
      },
      {
        promiseToast: {
          loading: t("Saving settings..."),
          success: t("Settings saved"),
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
  });

  const onSaveIcon = async (file: Blob | null) => {
    if (!file) {
      return;
    }

    await iconUpload(file);
    revalidate();
  };

  const onColourChangePreview = (colour: string) => {
    updateColour(colour);
  };

  const onContrastChange = (v: number) => {
    setContrast(v);
    updateColour(form.getValues("accentColour"));
  };

  const onClearMotdDates = () => {
    form.setValue("motdStartAt", "", {
      shouldDirty: true,
      shouldValidate: true,
    });
    form.setValue("motdEndAt", "", {
      shouldDirty: true,
      shouldValidate: true,
    });
  };

  const onClearMotd = () => {
    form.setValue("motdContent", undefined, {
      shouldDirty: true,
      shouldValidate: true,
    });
    form.setValue("motdStartAt", "", {
      shouldDirty: true,
      shouldValidate: true,
    });
    form.setValue("motdEndAt", "", {
      shouldDirty: true,
      shouldValidate: true,
    });
    form.setValue("motdType", undefined, {
      shouldDirty: true,
      shouldValidate: true,
    });

    setMotdContentInitialValue(undefined);
    setMotdContentResetKey(String(Date.now()));
  };

  return {
    register: form.register,
    control: form.control,
    formState: form.formState,
    onSubmit: onSubmit,
    currentIcon,
    onSaveIcon,
    onColourChangePreview,
    onContrastChange,
    onClearMotdDates,
    onClearMotd,
    motdContentInitialValue,
    motdContentResetKey,
    contrast,
  };
}

function toISODateTime(value?: string): string | undefined {
  if (!value) return undefined;

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return undefined;

  return parsed.toISOString();
}
