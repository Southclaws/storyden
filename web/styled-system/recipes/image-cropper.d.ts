import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ImageCropperVariant = {}

export type ImageCropperVariantProps = {
  [K in keyof ImageCropperVariant]?: ConditionalValue<ImageCropperVariant[K]>
}

export type ImageCropperVariantMap = RecipeVariantMap<ImageCropperVariant>

export type ImageCropperSlot = "root" | "viewport" | "cropper" | "controls" | "toolbar" | "editTrigger" | "reveal" | "actions" | "status" | "checkmark" | "fileUpload" | "hint" | "description"

export type ImageCropperRecipe = SlotRecipeRuntimeFn<ImageCropperSlot, ImageCropperVariantProps, ImageCropperVariantMap>

export declare const imageCropper: ImageCropperRecipe;