import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type FileUploadVariant = {}

export type FileUploadVariantProps = {
  [K in keyof FileUploadVariant]?: ConditionalValue<FileUploadVariant[K]>
}

export type FileUploadVariantMap = RecipeVariantMap<FileUploadVariant>

export type FileUploadSlot = "root" | "dropzone" | "item" | "itemDeleteTrigger" | "itemGroup" | "itemName" | "itemPreview" | "itemPreviewImage" | "itemSizeText" | "label" | "trigger" | "clearTrigger"

export type FileUploadRecipe = SlotRecipeRuntimeFn<FileUploadSlot, FileUploadVariantProps, FileUploadVariantMap>

export declare const fileUpload: FileUploadRecipe;