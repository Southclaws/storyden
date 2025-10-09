import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { mutate } from "swr";
import { z } from "zod";

import {
  categoryCreate,
  getCategoryListKey,
} from "src/api/openapi-client/categories";
import { Asset } from "src/api/openapi-schema";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { handle } from "@/api/client";

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name for the category."),
  slug: z.string().min(1, "Please enter a URL slug for the category."),
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
    control,
    handleImageUpload,
  };
}
