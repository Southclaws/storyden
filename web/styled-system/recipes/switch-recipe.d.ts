import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SwitchRecipeVariant = {
  size?: "lg" | "md" | "sm"
}

export type SwitchRecipeVariantProps = {
  [K in keyof SwitchRecipeVariant]?: ConditionalValue<SwitchRecipeVariant[K]>
}

export type SwitchRecipeVariantMap = RecipeVariantMap<SwitchRecipeVariant>

export type SwitchRecipeSlot = "root" | "label" | "control" | "thumb"

export type SwitchRecipeRecipe = SlotRecipeRuntimeFn<SwitchRecipeSlot, SwitchRecipeVariantProps, SwitchRecipeVariantMap>

export declare const switchRecipe: SwitchRecipeRecipe;