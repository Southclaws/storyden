import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type FileUploadVariant = {}

export type FileUploadVariantProps = {
  [K in keyof FileUploadVariant]?: ConditionalValue<FileUploadVariant[K]>
}

export type FileUploadVariantMap = RecipeVariantMap<FileUploadVariant>

export type FileUploadSlot = "root" | "dropzone" | "item" | "itemDeleteTrigger" | "itemGroup" | "itemName" | "itemPreview" | "itemPreviewImage" | "itemSizeText" | "label" | "trigger" | "clearTrigger"

export type FileUploadRecipe = ConfigSlotRecipeRuntimeFn<FileUploadSlot, FileUploadVariantProps, FileUploadVariantMap>

export declare const fileUpload: FileUploadRecipe;