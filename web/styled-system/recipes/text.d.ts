import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TextVariant = {
  variant?: "body" | "metadata" | "supporting"
}

export type TextVariantProps = {
  [K in keyof TextVariant]?: ConditionalValue<TextVariant[K]>
}

export type TextVariantMap = RecipeVariantMap<TextVariant>

export type TextRecipe = RecipeRuntimeFn<TextVariantProps, TextVariantMap>

export declare const text: TextRecipe;