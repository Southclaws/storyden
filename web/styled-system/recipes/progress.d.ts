import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ProgressVariant = {
  size?: "lg" | "md" | "sm"
}

export type ProgressVariantProps = {
  [K in keyof ProgressVariant]?: ConditionalValue<ProgressVariant[K]>
}

export type ProgressVariantMap = RecipeVariantMap<ProgressVariant>

export type ProgressSlot = "root" | "label" | "track" | "range" | "valueText" | "view" | "circle" | "circleTrack" | "circleRange"

export type ProgressRecipe = SlotRecipeRuntimeFn<ProgressSlot, ProgressVariantProps, ProgressVariantMap>

export declare const progress: ProgressRecipe;