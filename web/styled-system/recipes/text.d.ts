import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TextVariant = {
  variant?: "body" | "supporting" | "metadata"
}

export type TextVariantProps = {
  [K in keyof TextVariant]?: ConditionalValue<TextVariant[K]>
}

export type TextVariantMap = RecipeVariantMap<TextVariant>

export type TextRecipe = ConfigRecipeRuntimeFn<TextVariantProps, TextVariantMap>

export declare const text: TextRecipe;