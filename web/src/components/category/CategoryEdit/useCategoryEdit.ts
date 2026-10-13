import { zodResolver } from "@hookform/resolvers/zod";
import { usePathname, useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";

import { handle } from "@/api/client";
import { Asset, Category } from "@/api/openapi-schema";
import { useCategoryMutations } from "@/lib/category/mutation";
import * as z from "@/lib/validation";
import { slugify, slugifyDraft } from "@/utils/slugify";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export type Props = {
  category: Category;
} & UseDisclosureProps;

export const FormSchema = z.object({
  name: z.string().check(z.minLength(1)),
  slug: z.pipe(
    z.pipe(z.string(), z.transform(slugify)),
    z.string().check(z.minLength(1)),
  ),
  description: z.string().check(z.minLength(1)),
  colour: z.string(),
  cover_image: z.optional(z.nullable(z.custom<Asset>())),
});
export type Form = z.infer<typeof FormSchema>;

export function useCategoryEdit(props: Props) {
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: props.category.name,
      slug: props.category.slug,
      description: props.category.description,
      colour: props.category.colour,
      cover_image: props.category.cover_image || null,
    },
  });
  const pathname = usePathname();
  const router = useRouter();

  const slugInput = form.register("slug", {
    onChange(event) {
      const value = slugifyDraft(event.target.value);
      event.target.value = value;
      form.setValue("slug", value, {
        shouldDirty: true,
      });
    },
    onBlur(event) {
      const value = slugify(event.target.value);
      event.target.value = value;
      form.setValue("slug", value, {
        shouldDirty: true,
        shouldValidate: true,
      });
    },
  });

  const { revalidateList, updateCategory } = useCategoryMutations();

  const handleSubmit = form.handleSubmit(async (data) => {
    const originalSlug = props.category.slug;
    const slugChanged = data.slug !== originalSlug;

    await handle(
      async () => {
        const { cover_image, ...rest } = data;
        const updateData = {
          ...rest,
          cover_image_asset_id: cover_image?.id || null,
        };
        await updateCategory(originalSlug, updateData);

        props.onClose?.();

        if (slugChanged) {
          if (pathname.includes(`/d/${originalSlug}`)) {
            router.replace(`/d/${data.slug}`);
          }
        }
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
    slugInput,
    handlers: {
      handleSubmit,
      handleCancel,
      handleImageUpload,
    },
  };
}
