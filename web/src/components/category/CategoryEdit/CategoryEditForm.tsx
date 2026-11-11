import { useWatch } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { ColorPickerField } from "@/components/ui/color-picker";
import { FormControl } from "@/components/ui/form-control";
import { FormFeedback } from "@/components/ui/form-feedback";
import { FormLabel } from "@/components/ui/form-label";
import { ImageCropper } from "@/components/ui/image-cropper";
import { Input, InputPrefix } from "@/components/ui/input";
import { WEB_ADDRESS } from "@/config";
import {
  CATEGORY_COVER_HEIGHT,
  CATEGORY_COVER_SAFE_AREA,
  CATEGORY_COVER_WIDTH,
} from "@/lib/category/cover";
import { HStack, VStack, styled } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

import { type Props, useCategoryEdit } from "./useCategoryEdit";

export function CategoryEditForm(props: Props) {
  const { form, slugInput, handlers, cropperRef, imageLoading } =
    useCategoryEdit(props);
  const coverImage = useWatch({ control: form.control, name: "cover_image" });
  const hostname = new URL(WEB_ADDRESS).host;
  return (
    <styled.form
      display="flex"
      flexDir="column"
      justifyContent="space-between"
      alignItems="start"
      height="full"
      w="full"
      minW="0"
      onSubmit={handlers.handleSubmit}
      gap="2"
    >
      <VStack w="full">
        <FormControl>
          <FormLabel>Cover Image</FormLabel>
          <ImageCropper
            ref={cropperRef}
            aria-label="Category cover image"
            cropMode="fixed"
            safeArea={CATEGORY_COVER_SAFE_AREA}
            saving={form.formState.isSubmitting}
            disabled={form.formState.isSubmitting}
            width={CATEGORY_COVER_WIDTH}
            height={CATEGORY_COVER_HEIGHT}
            src={getAssetURL(coverImage?.path)}
            onFileChange={handlers.handleImageFileChange}
            onLoadError={handlers.handleImageError}
            onChange={handlers.handleImageChange}
            onRemove={handlers.handleImageRemove}
          />
          <FormFeedback error={form.formState.errors["cover_image"]?.message}>
            Safe area is used for social preview cards.
          </FormFeedback>
        </FormControl>

        <HStack w="full" alignItems="start">
          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input {...form.register("name")} type="text" />
            <FormFeedback error={form.formState.errors["name"]?.message}>
              The name of the category.
            </FormFeedback>
          </FormControl>

          <FormControl>
            <FormLabel>URL Slug</FormLabel>
            <HStack gap="0" alignItems="stretch" flex="1">
              <InputPrefix
                display={{
                  base: "none",
                  sm: "flex",
                }}
              >
                {hostname}/d/
              </InputPrefix>
              <Input
                {...slugInput}
                type="text"
                flex="1"
                borderTopLeftRadius={{
                  base: "sm",
                  sm: "none",
                }}
                borderBottomLeftRadius={{
                  base: "sm",
                  sm: "none",
                }}
              />
            </HStack>
            <FormFeedback error={form.formState.errors["slug"]?.message}>
              The URL path for the category.
            </FormFeedback>
          </FormControl>
        </HStack>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input {...form.register("description")} type="text" />
          <FormFeedback error={form.formState.errors["description"]?.message}>
            The description for the category.
          </FormFeedback>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <ColorPickerField
            control={form.control}
            name="colour"
            ariaLabel="Category colour"
          />
          <FormFeedback error={form.formState.errors["colour"]?.message}>
            The colour for the category.
          </FormFeedback>
        </FormControl>
      </VStack>

      <HStack w="full" alignItems="center" justify="end" pb="3" gap="4">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={handlers.handleCancel}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          size="sm"
          loading={form.formState.isSubmitting}
          disabled={imageLoading}
        >
          Save
        </Button>
      </HStack>
    </styled.form>
  );
}
