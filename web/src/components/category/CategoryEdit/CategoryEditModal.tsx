import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { AssetUploadEditor } from "@/components/asset/AssetUploadEditor/AssetUploadEditor";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormFeedback } from "@/components/ui/form/FormFeedback";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Input } from "@/components/ui/input";
import {
  CATEGORY_COVER_HEIGHT,
  CATEGORY_COVER_WIDTH,
} from "@/lib/category/cover";
import { HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useCategoryEdit } from "./useCategoryEdit";

export function CategoryEditModal(props: Props) {
  const { form, handlers } = useCategoryEdit(props);

  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      onOpenChange={props.onOpenChange}
      title="Edit category"
    >
      <styled.form
        display="flex"
        flexDir="column"
        justifyContent="space-between"
        alignItems="start"
        height="full"
        onSubmit={handlers.handleSubmit}
        gap="2"
      >
        <VStack w="full">
          <FormControl>
            <FormLabel>Cover Image</FormLabel>
            <AssetUploadEditor
              width={CATEGORY_COVER_WIDTH}
              height={CATEGORY_COVER_HEIGHT}
              value={form.watch("cover_image") || undefined}
              onUpload={handlers.handleImageUpload}
            />
            <FormFeedback error={form.formState.errors["cover_image"]?.message}>
              Upload a cover image for the category (4:1 aspect ratio).
            </FormFeedback>
          </FormControl>

          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input {...form.register("name")} type="text" />
            <FormFeedback error={form.formState.errors["name"]?.message}>
              The name of the category.
            </FormFeedback>
          </FormControl>

          <FormControl>
            <FormLabel>Description</FormLabel>
            <Input {...form.register("description")} type="text" />
            <FormFeedback error={form.formState.errors["description"]?.message}>
              The description for the category.
            </FormFeedback>
          </FormControl>

          <FormControl>
            <FormLabel>Colour</FormLabel>
            <ColourPickerField control={form.control} name="colour" />
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
          <Button type="submit" size="sm">
            Save
          </Button>
        </HStack>
      </styled.form>
    </ModalDrawer>
  );
}
