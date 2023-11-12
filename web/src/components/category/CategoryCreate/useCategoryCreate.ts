import { useToast } from "@chakra-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { categoryCreate, getCategoryListKey } from "src/api/openapi/categories";
import { APIError } from "src/api/openapi/schemas";
import { errorToast } from "src/components/site/ErrorBanner";
import { UseDisclosureProps } from "src/theme/components";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the category."),
  description: z.string().min(1, "Please enter a short description."),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryCreate(props: UseDisclosureProps) {
  const toast = useToast();
  const { register, handleSubmit } = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    try {
      const category = await categoryCreate(data);
      props.onClose?.();
      mutate(getCategoryListKey());
      toast({
        title: "Category created",
        description: `${category.name} is now ready to be filled with stuff!`,
      });
    } catch (e: unknown) {
      errorToast(toast)(e as APIError);
    }
  });

  return {
    onSubmit,
    register,
  };
}
