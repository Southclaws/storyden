import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ToggleGroupVariant = {
  size?: "lg" | "md" | "sm" | "xs"
  variant?: "ghost" | "outline"
}

export type ToggleGroupVariantProps = {
  [K in keyof ToggleGroupVariant]?: ToggleGroupVariant[K]
}

export type ToggleGroupVariantMap = RecipeVariantMap<ToggleGroupVariant>

export type ToggleGroupSlot = "root" | "item"

export type ToggleGroupRecipe = SlotRecipeRuntimeFn<ToggleGroupSlot, ToggleGroupVariantProps, ToggleGroupVariantMap>

export declare const toggleGroup: ToggleGroupRecipe;