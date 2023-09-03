import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { adminSettingsUpdate } from "src/api/openapi/admin";
import { getGetInfoKey } from "src/api/openapi/misc";
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

  const updateColour = (colour: string) => {
    const cv = getColourVariants(colour);

    Object.entries(cv).forEach((property) =>
      document.documentElement.style.setProperty(property[0], property[1])
    );
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

  const onColourChangePreview = (colour: string) => {
    updateColour(colour);
  };

  return {
    register: form.register,
    control: form.control,
    onSubmit: onSubmit,
    onColourChangePreview,
  };
}
