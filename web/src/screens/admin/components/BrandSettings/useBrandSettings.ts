import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { adminSettingsUpdate } from "src/api/openapi-client/admin";
import { getGetInfoKey, iconUpload } from "src/api/openapi-client/misc";
import { Info } from "src/api/openapi-schema";
import { getColourVariants } from "src/utils/colour";

import { getIconURL } from "@/utils/icon";

export type Props = Info;

export const FormSchema = z.object({
  title: z.string(),
  description: z.string(),
  accentColour: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useBrandSettings(props: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      title: props.title,
      description: props.description,
      accentColour: props.accent_colour,
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
    updateColour(data.accentColour);
    await adminSettingsUpdate({
      title: data.title,
      description: data.description,
      accent_colour: data.accentColour,
    });
    mutate(getGetInfoKey());
  });

  const onSaveIcon = async (file: Blob | null) => {
    if (!file) {
      return;
    }

    await iconUpload(file);
    mutate(getGetInfoKey());
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
    onSubmit: onSubmit,
    currentIcon,
    onSaveIcon,
    onColourChangePreview,
    onContrastChange,
    contrast,
  };
}
