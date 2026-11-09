import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CardBoxVariant = {
  tone?: "default" | "success"
  kind?: "default" | "edge"
}

export type CardBoxVariantProps = {
  [K in keyof CardBoxVariant]?: ConditionalValue<CardBoxVariant[K]>
}

export type CardBoxVariantMap = RecipeVariantMap<CardBoxVariant>

export type CardBoxRecipe = ConfigRecipeRuntimeFn<CardBoxVariantProps, CardBoxVariantMap>

export declare const cardBox: CardBoxRecipe;