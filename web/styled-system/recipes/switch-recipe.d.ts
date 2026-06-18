import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SwitchRecipeVariant = {
  size?: "lg" | "md" | "sm"
}

export type SwitchRecipeVariantProps = {
  [K in keyof SwitchRecipeVariant]?: ConditionalValue<SwitchRecipeVariant[K]>
}

export type SwitchRecipeVariantMap = RecipeVariantMap<SwitchRecipeVariant>

export type SwitchRecipeRecipe = RecipeRuntimeFn<SwitchRecipeVariantProps, SwitchRecipeVariantMap>

export declare const switchRecipe: SwitchRecipeRecipe;