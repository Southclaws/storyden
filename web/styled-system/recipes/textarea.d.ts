import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TextareaVariant = {
  size?: "lg" | "md" | "sm"
  variant?: "ghost" | "inset" | "outline"
}

export type TextareaVariantProps = {
  [K in keyof TextareaVariant]?: TextareaVariant[K] | undefined
}

export type TextareaVariantMap = RecipeVariantMap<TextareaVariant>

export type TextareaRecipe = RecipeRuntimeFn<TextareaVariantProps, TextareaVariantMap>

export declare const textarea: TextareaRecipe;