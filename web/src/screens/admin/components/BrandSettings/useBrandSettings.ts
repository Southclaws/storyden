import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { adminSettingsUpdate } from "src/api/openapi/admin";
import { getGetInfoKey, iconUpload } from "src/api/openapi/misc";
import { Info } from "src/api/openapi/schemas";
import { getColourVariants } from "src/utils/colour";

export type Props = Info;

export const FormSchema = z.object({
  title: z.string(),
  description: z.string(),
  accentColour: z.string(),
});
export type Form = z.infer<typeof FormSchema>;

export function useBrandSettings(props: Props) {
  const toast = useToast();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      title: props.title,
      description: props.description,
      accentColour: props.accent_colour,
    },
  });
  const [currentIcon, setCurrentIcon] = useState<File | undefined>(undefined);

  useEffect(() => {
    (async () => {
      const icon = await fetch(`/api/v1/info/icon/512x512`);

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
    toast({ title: "Settings updated!" });
  });

  const onSaveIcon = async (file: File) => {
    await iconUpload(file);
    mutate(getGetInfoKey());
    toast({ title: "Icon updated!" });
  };

  const onColourChangePreview = (colour: string) => {
    updateColour(colour);
  };

  return {
    register: form.register,
    control: form.control,
    onSubmit: onSubmit,
    currentIcon,
    onSaveIcon,
    onColourChangePreview,
  };
}
