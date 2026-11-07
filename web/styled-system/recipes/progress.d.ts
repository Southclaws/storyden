import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ProgressVariant = {
  shape?: "circle" | "horizontal"
  size?: "sm" | "md" | "lg"
}

export type ProgressVariantProps = {
  [K in keyof ProgressVariant]?: ConditionalValue<ProgressVariant[K]>
}

export type ProgressVariantMap = RecipeVariantMap<ProgressVariant>

export type ProgressSlot = "root" | "label" | "track" | "range" | "valueText" | "view" | "circle" | "circleTrack" | "circleRange"

export type ProgressRecipe = ConfigSlotRecipeRuntimeFn<ProgressSlot, ProgressVariantProps, ProgressVariantMap>

export declare const progress: ProgressRecipe;