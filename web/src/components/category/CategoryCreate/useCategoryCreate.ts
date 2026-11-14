import { zodResolver } from "@hookform/resolvers/zod";
import { BaseSyntheticEvent, useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { mutate } from "swr";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import {
  categoryCreate,
  getCategoryListKey,
} from "@/api/openapi-client/categories";
import {
  ImageCropperRef,
  ImageCropperSnapshot,
} from "@/components/ui/image-cropper";
import * as z from "@/lib/validation";
import { slugify, slugifyDraft } from "@/utils/slugify";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export const FormSchema = z.object({
  name: z
    .string()
    .check(z.minLength(1, "Please enter a name for the category.")),
  slug: z.pipe(
    z.pipe(z.string(), z.transform(slugify)),
    z
      .string()
      .check(z.minLength(1, "Please enter a URL slug for the category.")),
  ),
  description: z
    .string()
    .check(z.minLength(1, "Please enter a short description.")),
  colour: z.string(),
  parent: z.optional(z.string()),
  cover_image_asset_id: z.optional(z.string()),
});
export type Form = z.infer<typeof FormSchema>;

export interface CategoryCreateProps extends UseDisclosureProps {
  defaultParent?: string;
}

export function useCategoryCreate(props: CategoryCreateProps) {
  const submissionRevision = useRef(0);

  useEffect(() => {
    return () => {
      submissionRevision.current += 1;
    };
  }, [props.isOpen]);

  const cropperRef = useRef<ImageCropperRef>(null);
  const [pendingCrop, setPendingCrop] = useState<ImageCropperSnapshot | null>(
    null,
  );
  const [imageLoading, setImageLoading] = useState(false);
  const { register, handleSubmit, control, setValue, formState } =
    useForm<Form>({
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

  function onSubmit(event?: BaseSyntheticEvent) {
    const revision = ++submissionRevision.current;

    return handleSubmit(async (data) => {
      if (imageLoading || revision !== submissionRevision.current) {
        return;
      }

      await handle(async () => {
        const crop = pendingCrop;
        if (crop) {
          const blob = await crop.getBlob();
          if (revision !== submissionRevision.current) {
            return;
          }

          const asset = await assetUpload(blob, {
            filename: "category-cover",
          });
          data.cover_image_asset_id = asset.id;
        }

        if (revision !== submissionRevision.current) {
          return;
        }

        await categoryCreate(data);
        props.onClose?.();
        mutate(getCategoryListKey());
      });
    })(event);
  }

  function handleCancel() {
    submissionRevision.current += 1;
    setImageLoading(false);
    setPendingCrop(null);
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
    setValue("cover_image_asset_id", undefined, { shouldDirty: true });
  }

  return {
    onSubmit,
    handleCancel,
    register,
    slugInput,
    control,
    cropperRef,
    imageLoading,
    isSubmitting: formState.isSubmitting,
    handleImageFileChange,
    handleImageChange,
    handleImageError,
    handleImageRemove,
  };
}
