import { zodResolver } from "@hookform/resolvers/zod";
import { usePathname, useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { useForm } from "react-hook-form";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset, Category } from "@/api/openapi-schema";
import {
  ImageCropperRef,
  ImageCropperSnapshot,
} from "@/components/ui/image-cropper";
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
  const cropperRef = useRef<ImageCropperRef>(null);
  const [pendingCrop, setPendingCrop] = useState<ImageCropperSnapshot | null>(
    null,
  );
  const [imageLoading, setImageLoading] = useState(false);
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
    if (imageLoading) {
      return;
    }

    const originalSlug = props.category.slug;
    const slugChanged = data.slug !== originalSlug;

    await handle(
      async () => {
        let { cover_image, ...rest } = data;
        const crop = pendingCrop;
        if (crop) {
          cover_image = await assetUpload(await crop.getBlob(), {
            filename: "category-cover",
          });
        }

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
    setImageLoading(false);
    setPendingCrop(null);
    form.reset();
    props.onClose?.();
  }

  async function handleImageFileChange() {
    setImageLoading(true);
    setPendingCrop(null);
  }

  function handleImageChange() {
    setImageLoading(false);
    setPendingCrop(cropperRef.current?.getCrop() ?? null);
  }

  function handleImageError() {
    setImageLoading(false);
    setPendingCrop(null);
  }

  async function handleImageRemove() {
    setImageLoading(false);
    setPendingCrop(null);
    form.setValue("cover_image", null, { shouldDirty: true });
  }

  return {
    form,
    cropperRef,
    imageLoading,
    slugInput,
    handlers: {
      handleSubmit,
      handleCancel,
      handleImageFileChange,
      handleImageChange,
      handleImageError,
      handleImageRemove,
    },
  };
}
