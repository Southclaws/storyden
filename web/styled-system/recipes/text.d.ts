import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TextVariant = {
  size?: "2xl" | "3xl" | "4xl" | "5xl" | "6xl" | "7xl" | "lg" | "md" | "sm" | "xl" | "xs"
  variant?: "heading"
}

export type TextVariantProps = {
  [K in keyof TextVariant]?: ConditionalValue<TextVariant[K]>
}

export type TextVariantMap = RecipeVariantMap<TextVariant>

export type TextRecipe = RecipeRuntimeFn<TextVariantProps, TextVariantMap>

export declare const text: TextRecipe;