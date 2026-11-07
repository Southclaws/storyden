import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SwitchRecipeVariant = {
  size?: "sm" | "md" | "lg"
}

export type SwitchRecipeVariantProps = {
  [K in keyof SwitchRecipeVariant]?: ConditionalValue<SwitchRecipeVariant[K]>
}

export type SwitchRecipeVariantMap = RecipeVariantMap<SwitchRecipeVariant>

export type SwitchRecipeSlot = "root" | "label" | "control" | "thumb"

export type SwitchRecipeRecipe = ConfigSlotRecipeRuntimeFn<SwitchRecipeSlot, SwitchRecipeVariantProps, SwitchRecipeVariantMap>

export declare const switchRecipe: SwitchRecipeRecipe;