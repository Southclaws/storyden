import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { categoryUpdate, useCategoryList } from "src/api/openapi/categories";
import { useGetInfo } from "src/api/openapi/misc";
import { APIError, Category } from "src/api/openapi/schemas";
import { errorToast } from "src/components/site/ErrorBanner";
import { UseDisclosureProps, useToast } from "src/theme/components";

export type Props = {
  category: Category;
} & UseDisclosureProps;

export const FormSchema = z.object({
  name: z.string(),
  description: z.string(),
  colour: z.string().default("#fff"), // not implemented yet
  admin: z.boolean().default(false), // not implemented yet
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryEdit(props: Props) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(FormSchema),
    values: {
      name: props.category.name,
      description: props.category.description,
      colour: props.category.colour,
      admin: props.category.admin,
    },
  });
  const { mutate, data: existing } = useCategoryList();
  const { mutate: mutateInfoStatus } = useGetInfo();

  const onSubmit = handleSubmit(async (data) => {
    try {
      const collection = await categoryUpdate(props.category.id, data);
      const updated = [...(existing?.categories ?? []), collection];

      mutateInfoStatus();
      mutate(
        { categories: updated },
        { populateCache: true, rollbackOnError: true },
      );

      props.onClose?.();
      toast({
        title: "Category updated",
      });
    } catch (e: unknown) {
      errorToast(toast)(e as APIError);
    }
  });

  return {
    onSubmit,
    register,
    errors,
  };
}
