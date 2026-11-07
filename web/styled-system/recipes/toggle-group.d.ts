import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ToggleGroupVariant = {
  variant?: "outline" | "ghost"
  size?: "xs" | "sm" | "md" | "lg"
}

export type ToggleGroupVariantProps = {
  [K in keyof ToggleGroupVariant]?: ToggleGroupVariant[K] | undefined
}

export type ToggleGroupVariantMap = RecipeVariantMap<ToggleGroupVariant>

export type ToggleGroupSlot = "root" | "item"

export type ToggleGroupRecipe = ConfigSlotRecipeRuntimeFn<ToggleGroupSlot, ToggleGroupVariantProps, ToggleGroupVariantMap>

export declare const toggleGroup: ToggleGroupRecipe;