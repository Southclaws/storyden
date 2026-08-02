import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import { handle } from "@/api/client";
import {
  categoryCreate,
  getCategoryListKey,
} from "@/api/openapi-client/categories";
import { Asset } from "@/api/openapi-schema";
import { slugify, slugifyDraft } from "@/utils/slugify";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the category."),
  slug: z
    .string()
    .transform(slugify)
    .pipe(z.string().min(1, "Please enter a URL slug for the category.")),
  description: z.string().min(1, "Please enter a short description."),
  colour: z.string().default("#8577ce"),
  parent: z.string().optional(),
  cover_image_asset_id: z.string().optional(),
});
export type Form = z.infer<typeof FormSchema>;

export interface CategoryCreateProps extends UseDisclosureProps {
  defaultParent?: string;
}

export function useCategoryCreate(props: CategoryCreateProps) {
  const { register, handleSubmit, control, setValue } = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      colour: "#8577ce",
      parent: props.defaultParent,
    },
  });

  const slugInput = register("slug", {
    onChange(event) {
      const value = slugifyDraft(event.target.value);
      event.target.value = value;
      setValue("slug", value, {
        shouldDirty: true,
      });
    },
    onBlur(event) {
      const value = slugify(event.target.value);
      event.target.value = value;
      setValue("slug", value, {
        shouldDirty: true,
        shouldValidate: true,
      });
    },
  });

  const onSubmit = handleSubmit(async (data) => {
    await handle(async () => {
      await categoryCreate(data);
      props.onClose?.();
      mutate(getCategoryListKey());
    });
  });

  function handleImageUpload(asset: Asset) {
    setValue("cover_image_asset_id", asset.id);
  }

  return {
    onSubmit,
    register,
    slugInput,
    control,
    handleImageUpload,
  };
}
