import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Info } from "src/api/openapi/schemas";

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

  const onSubmit = form.handleSubmit((data) => {
    alert(JSON.stringify(data, null, "  "));
  });

  return {
    register: form.register,
    control: form.control,
    onSubmit: onSubmit,
  };
}
