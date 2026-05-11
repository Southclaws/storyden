"use client";

import { createListCollection } from "@ark-ui/react";

import { useCategoryList } from "src/api/openapi-client/categories";

import { AssetUploadEditor } from "@/components/asset/AssetUploadEditor/AssetUploadEditor";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { SelectField } from "@/components/ui/form/SelectField";
import { Input, InputPrefix } from "@/components/ui/input";
import { WEB_ADDRESS } from "@/config";
import { useI18n } from "@/i18n/provider";
import {
  CATEGORY_COVER_HEIGHT,
  CATEGORY_COVER_WIDTH,
} from "@/lib/category/cover";
import { HStack, VStack, WStack, styled } from "@/styled-system/jsx";

import { CategoryCreateProps, useCategoryCreate } from "./useCategoryCreate";

export type { CategoryCreateProps };

export function CategoryCreateScreen(props: CategoryCreateProps) {
  const { t } = useI18n();
  const { register, onSubmit, control, handleImageUpload } =
    useCategoryCreate(props);

  const { data: categoryListResult } = useCategoryList();
  const categories = categoryListResult?.categories || [];

  const categoryCollection = createListCollection({
    items: [
      { label: t("No parent (root category)"), value: "" },
      ...categories.map((category) => ({
        label: category.name,
        value: category.id,
      })),
    ],
  });

  const hostname = new URL(WEB_ADDRESS).host;

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>
        {t(
          "Use categories to organise posts. A post can only have one category, unlike tags. So it's best to keep categories high-level and different enough so that it's not easy to get confused between them.",
        )}
      </styled.p>
      <styled.form
        display="flex"
        flexDir="column"
        gap="4"
        w="full"
        onSubmit={onSubmit}
      >
        <FormControl>
          <FormLabel>{t("Cover Image")}</FormLabel>
          <AssetUploadEditor
            width={CATEGORY_COVER_WIDTH}
            height={CATEGORY_COVER_HEIGHT}
            onUpload={handleImageUpload}
          />
          <FormHelperText>
            {t("Upload a cover image for the category (4:1 aspect ratio)")}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Name")}</FormLabel>
          <Input {...register("name")} type="text" />
          <FormHelperText>{t("The name for your category")}</FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("URL Slug")}</FormLabel>
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
              {...register("slug")}
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
          <FormHelperText>
            {t(
              'The URL path for your category (e.g., "general", "announcements")',
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Description")}</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>{t("Describe your category")}</FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Parent Category")}</FormLabel>
          <SelectField
            name="parent"
            control={control}
            collection={categoryCollection}
            placeholder={t("Select a parent category")}
          />
          <FormHelperText>
            {t(
              "Choose a parent category to create a subcategory, or leave as root category",
            )}
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>{t("Colour")}</FormLabel>
          <ColourPickerField control={control} name="colour" />
          <FormHelperText>{t("The colour for the category")}</FormHelperText>
        </FormControl>

        <WStack>
          <Button flexGrow="1" type="button" onClick={props.onClose}>
            {t("Cancel")}
          </Button>
          <Button flexGrow="1" type="submit">
            {t("Create")}
          </Button>
        </WStack>
      </styled.form>
    </VStack>
  );
}
