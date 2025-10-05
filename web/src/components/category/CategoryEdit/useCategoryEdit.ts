import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Asset, Category } from "src/api/openapi-schema";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import { useCategoryMutations } from "@/lib/category/mutation";

export type Props = {
  category: Category;
} & UseDisclosureProps;

export const FormSchema = z.object({
  name: z.string().min(1),
  description: z.string().min(1),
  colour: z.string().default("#fff"),
  cover_image: z.custom<Asset>().nullable().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryEdit(props: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: props.category.name,
      description: props.category.description,
      colour: props.category.colour,
      cover_image: props.category.cover_image || null,
    },
  });

  const { revalidateList, updateCategory } = useCategoryMutations();

  const handleSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const { cover_image, ...rest } = data;
        const updateData = {
          ...rest,
          cover_image_asset_id: cover_image?.id || null,
        };
        await updateCategory(props.category.slug, updateData);

        props.onClose?.();
      },
      {
        promiseToast: {
          loading: "Updating category...",
          success: "Category updated.",
        },
        cleanup: () => revalidateList(),
      },
    );
  });

  function handleCancel() {
    form.reset();
    props.onClose?.();
  }

  function handleImageUpload(asset: Asset) {
    form.setValue("cover_image", asset);
  }

  return {
    form,
    handlers: {
      handleSubmit,
      handleCancel,
      handleImageUpload,
    },
  };
}
