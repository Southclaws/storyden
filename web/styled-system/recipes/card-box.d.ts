import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CardBoxVariant = {
  kind?: "default" | "edge"
}

export type CardBoxVariantProps = {
  [K in keyof CardBoxVariant]?: ConditionalValue<CardBoxVariant[K]>
}

export type CardBoxVariantMap = RecipeVariantMap<CardBoxVariant>

export type CardBoxRecipe = RecipeRuntimeFn<CardBoxVariantProps, CardBoxVariantMap>

export declare const cardBox: CardBoxRecipe;