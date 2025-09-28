import { createListCollection } from "@ark-ui/react";

import { useCategoryList } from "src/api/openapi-client/categories";

import { AssetUploadEditor } from "@/components/asset/AssetUploadEditor/AssetUploadEditor";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { SelectField } from "@/components/ui/form/SelectField";
import { Input } from "@/components/ui/input";
import { VStack, WStack, styled } from "@/styled-system/jsx";

import { CategoryCreateProps, useCategoryCreate } from "./useCategoryCreate";

export type { CategoryCreateProps };

export function CategoryCreateScreen(props: CategoryCreateProps) {
  const { register, onSubmit, control, handleImageUpload } =
    useCategoryCreate(props);

  const { data: categoryListResult } = useCategoryList();
  const categories = categoryListResult?.categories || [];

  const categoryCollection = createListCollection({
    items: [
      { label: "No parent (root category)", value: "" },
      ...categories.map((category) => ({
        label: category.name,
        value: category.id,
      })),
    ],
  });

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>
        Use categories to organise posts. A post can only have one category,
        unlike tags. So it&apos;s best to keep categories high-level and
        different enough so that it&apos;s not easy to get confused between
        them.
      </styled.p>
      <styled.form
        display="flex"
        flexDir="column"
        gap="4"
        w="full"
        onSubmit={onSubmit}
      >
        <FormControl>
          <FormLabel>Cover Image</FormLabel>
          <AssetUploadEditor
            aspectRatio="16 / 9"
            onUpload={handleImageUpload}
          />
          <FormHelperText>
            Upload a cover image for the category (16:9 aspect ratio)
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...register("name")} type="text" />
          <FormHelperText>The name for your category</FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>Description</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>Describe your category</FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Parent Category</FormLabel>
          <SelectField
            name="parent"
            control={control}
            collection={categoryCollection}
            placeholder="Select a parent category"
          />
          <FormHelperText>
            Choose a parent category to create a subcategory, or leave as root
            category
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <ColourPickerField control={control} name="colour" />
          <FormHelperText>The colour for the category</FormHelperText>
        </FormControl>

        <WStack>
          <Button flexGrow="1" type="button" onClick={props.onClose}>
            Cancel
          </Button>
          <Button flexGrow="1" type="submit">
            Create
          </Button>
        </WStack>
      </styled.form>
    </VStack>
  );
}
