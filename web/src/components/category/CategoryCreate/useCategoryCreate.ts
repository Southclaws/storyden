import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import {
  categoryCreate,
  getCategoryListKey,
} from "src/api/openapi-client/categories";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the category."),
  description: z.string().min(1, "Please enter a short description."),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryCreate(props: UseDisclosureProps) {
  const { register, handleSubmit } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    await handle(async () => {
      await categoryCreate(data);
      props.onClose?.();
      mutate(getCategoryListKey());
    });
  });

  return {
    onSubmit,
    register,
  };
}
