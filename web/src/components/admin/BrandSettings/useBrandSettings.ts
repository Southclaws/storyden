import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { iconUpload } from "src/api/openapi-client/misc";
import { getColourVariants } from "src/utils/colour";

import { handle } from "@/api/client";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { Settings } from "@/lib/settings/settings";
import { getIconURL } from "@/utils/icon";

export type Props = {
  settings: Settings;
};

export const FormSchema = z.object({
  title: z.string(),
  description: z.string(),
  content: z.string().optional(),
  accentColour: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useBrandSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      title: settings.title,
      description: settings.description,
      content: settings.content,
      accentColour: settings.accent_colour,
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
        updateColour(data.accentColour);
        await updateSettings({
          title: data.title,
          description: data.description,
          content: data.content,
          accent_colour: data.accentColour,
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
    updateColour(form.getValues().accentColour);
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
    contrast,
  };
}
