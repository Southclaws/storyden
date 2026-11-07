import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TextareaVariant = {
  size?: "sm" | "md" | "lg"
  variant?: "outline" | "ghost" | "inset"
}

export type TextareaVariantProps = {
  [K in keyof TextareaVariant]?: TextareaVariant[K] | undefined
}

export type TextareaVariantMap = RecipeVariantMap<TextareaVariant>

export type TextareaRecipe = ConfigRecipeRuntimeFn<TextareaVariantProps, TextareaVariantMap>

export declare const textarea: TextareaRecipe;