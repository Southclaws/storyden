import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SliderVariant = {
  size?: "lg" | "md" | "sm"
}

export type SliderVariantProps = {
  [K in keyof SliderVariant]?: ConditionalValue<SliderVariant[K]>
}

export type SliderVariantMap = RecipeVariantMap<SliderVariant>

export type SliderSlot = "root" | "label" | "thumb" | "valueText" | "track" | "range" | "control" | "markerGroup" | "marker" | "draggingIndicator"

export type SliderRecipe = SlotRecipeRuntimeFn<SliderSlot, SliderVariantProps, SliderVariantMap>

export declare const slider: SliderRecipe;